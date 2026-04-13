package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/deploymenttheory/go-sdk-brew/brew/constants"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// Transport is the HTTP transport layer for the Homebrew API.
// It wraps a resty.Client with behaviour: idempotent-only retries with
// exponential backoff, adaptive response-time throttling, optional
// concurrency limiting, and structured logging.
type Transport struct {
	client          *resty.Client
	logger          *zap.Logger
	BaseURL         string
	globalHeaders   map[string]string
	userAgent       string
	sem             *semaphore
	requestDelay    time.Duration
	totalRetryDuration time.Duration
	responseTracker *responseTimeTracker
}

// GetHTTPClient returns the underlying resty client for advanced use.
func (t *Transport) GetHTTPClient() *resty.Client {
	return t.client
}

// GetLogger returns the configured logger.
func (t *Transport) GetLogger() *zap.Logger {
	return t.logger
}

// NewTransport creates and fully configures a Homebrew API transport.
//
// Behaviour applied at construction time (resty native where possible):
//   - Idempotent-only retry (all GET) with exponential backoff
//   - Adaptive inter-request delay derived from response-time EMA tracking
func NewTransport(baseURL string, opts ...ClientOption) (*Transport, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	settings := &TransportSettings{
		GlobalHeaders: make(map[string]string),
	}
	for _, opt := range opts {
		if err := opt(settings); err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	logger := settings.Logger
	if logger == nil {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger: %w", err)
		}
	}

	// Option overrides the caller-supplied baseURL.
	if settings.BaseURL != "" {
		baseURL = settings.BaseURL
	}
	baseURL = trimTrailingSlash(baseURL)

	userAgent := settings.UserAgent
	if userAgent == "" {
		userAgent = fmt.Sprintf("%s/%s", UserAgentBase, constants.Version)
	}

	timeout := settings.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	retryCount := settings.RetryCount
	if retryCount == 0 {
		retryCount = MaxRetries
	}
	retryWait := settings.RetryWaitTime
	if retryWait == 0 {
		retryWait = RetryWaitTime
	}
	retryMaxWait := settings.RetryMaxWaitTime
	if retryMaxWait == 0 {
		retryMaxWait = RetryMaxWaitTime
	}

	restyClient := resty.New()
	restyClient.SetBaseURL(baseURL)
	restyClient.SetTimeout(timeout)
	restyClient.SetRetryCount(retryCount)
	restyClient.SetRetryWaitTime(retryWait)
	restyClient.SetRetryMaxWaitTime(retryMaxWait)
	restyClient.SetHeader("User-Agent", userAgent)

	// Only retry on transient server errors.
	restyClient.AddRetryConditions(retryCondition)

	if settings.Debug {
		restyClient.SetDebug(true)
	}

	if settings.InsecureSkipVerify {
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec
	} else if settings.TLSClientConfig != nil {
		restyClient.SetTLSClientConfig(settings.TLSClientConfig)
	}

	if settings.ProxyURL != "" {
		restyClient.SetProxy(settings.ProxyURL)
	}
	if settings.HTTPTransport != nil {
		restyClient.SetTransport(settings.HTTPTransport)
	}
	for k, v := range settings.GlobalHeaders {
		restyClient.SetHeader(k, v)
	}

	var sem *semaphore
	if settings.MaxConcurrentRequests > 0 {
		sem = newSemaphore(settings.MaxConcurrentRequests)
	}

	transport := &Transport{
		client:             restyClient,
		logger:             logger,
		BaseURL:            baseURL,
		globalHeaders:      settings.GlobalHeaders,
		userAgent:          userAgent,
		responseTracker:    newResponseTimeTracker(),
		sem:                sem,
		requestDelay:       settings.MandatoryRequestDelay,
		totalRetryDuration: settings.TotalRetryDuration,
	}

	// Apply OpenTelemetry instrumentation (always enabled, uses global providers).
	// If no global providers are configured, this is a no-op.
	transport.applyOpenTelemetry()

	logger.Info("Homebrew API transport created",
		zap.String("base_url", transport.BaseURL),
	)
	return transport, nil
}

// trimTrailingSlash removes a trailing slash from s.
func trimTrailingSlash(s string) string {
	if len(s) > 0 && s[len(s)-1] == '/' {
		return s[:len(s)-1]
	}
	return s
}

// NewRequest returns a RequestBuilder for this transport.
func (t *Transport) NewRequest(ctx context.Context) *RequestBuilder {
	return &RequestBuilder{
		req:      t.client.R().SetContext(ctx).SetResponseBodyUnlimitedReads(true),
		executor: t,
	}
}

// execute implements requestExecutor for Transport.
func (t *Transport) execute(req *resty.Request, method, path string, _ any) (*resty.Response, error) {
	return t.executeRequest(req, method, path)
}

// executeGetBytes implements requestExecutor for Transport.
func (t *Transport) executeGetBytes(req *resty.Request, path string) (*resty.Response, []byte, error) {
	resp, err := t.executeRequest(req, "GET", path)
	if err != nil {
		return resp, nil, err
	}
	return resp, resp.Bytes(), nil
}

// executeRequest is the central request executor used by all HTTP verb methods.
func (t *Transport) executeRequest(req *resty.Request, method, path string) (*resty.Response, error) {
	ctx := req.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// Wrap in a deadline for the total allowed retry window if configured.
	if t.totalRetryDuration > 0 {
		if _, hasDeadline := ctx.Deadline(); !hasDeadline {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, t.totalRetryDuration)
			defer cancel()
			req.SetContext(ctx)
		}
	}

	// Acquire concurrency slot — blocks until available or context cancelled.
	if t.sem != nil {
		if err := t.sem.acquire(ctx); err != nil {
			return nil, fmt.Errorf("concurrency limit: %w", err)
		}
		defer t.sem.release()
	}

	t.logger.Debug("Executing API request", zap.String("method", method), zap.String("path", path))

	resp, execErr := req.Execute(method, path)

	if execErr != nil {
		t.logger.Error("Request failed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Error(execErr),
		)
		return resp, fmt.Errorf("request failed: %w", execErr)
	}

	if err := t.validateResponse(resp, method, path); err != nil {
		return resp, err
	}

	if resp.IsError() {
		return resp, ParseErrorResponse(
			[]byte(resp.String()),
			resp.StatusCode(),
			resp.Status(),
			method,
			path,
			t.logger,
		)
	}

	duration := resp.Duration()

	t.logger.Info("Request completed",
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", resp.StatusCode()),
		zap.Duration("duration", duration),
	)

	// Mandatory fixed delay (user-configured).
	if t.requestDelay > 0 {
		time.Sleep(t.requestDelay)
	}

	// Adaptive delay: when the server is responding more slowly than its EMA baseline.
	if adaptive := t.responseTracker.record(duration); adaptive > 0 {
		t.logger.Debug("Adaptive delay applied due to elevated response time",
			zap.Duration("response_time", duration),
			zap.Duration("adaptive_delay", adaptive),
		)
		time.Sleep(adaptive)
	}

	return resp, nil
}
