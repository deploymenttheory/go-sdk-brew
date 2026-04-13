package client

import (
	"crypto/tls"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ClientOption is a function that mutates TransportSettings.
// Options are applied before the Transport is constructed, so all
// configuration is collected into a single TransportSettings value
// before any network state is initialised.
type ClientOption func(*TransportSettings) error

// TransportSettings collects all optional transport configuration. Each field
// has a zero value that signals "use the built-in default".
type TransportSettings struct {
	// BaseURL overrides the base URL passed to NewTransport when non-empty.
	BaseURL string

	// Timeout overrides the default HTTP request timeout (30 s) when non-zero.
	Timeout time.Duration

	// RetryCount overrides the default retry count (3) when non-zero.
	RetryCount int

	// RetryWaitTime overrides the default retry wait (2 s) when non-zero.
	RetryWaitTime time.Duration

	// RetryMaxWaitTime overrides the default max retry wait (30 s) when non-zero.
	RetryMaxWaitTime time.Duration

	// Logger replaces the default production zap logger when non-nil.
	Logger *zap.Logger

	// Debug enables resty's request/response debug logging when true.
	Debug bool

	// UserAgent replaces the default SDK user-agent string when non-empty.
	UserAgent string

	// GlobalHeaders are added to every outgoing request.
	GlobalHeaders map[string]string

	// ProxyURL sets an HTTP proxy for all requests when non-empty.
	ProxyURL string

	// TLSClientConfig sets custom TLS configuration. Ignored when
	// InsecureSkipVerify is true (InsecureSkipVerify takes precedence).
	TLSClientConfig *tls.Config

	// HTTPTransport replaces the default net/http transport when non-nil.
	HTTPTransport http.RoundTripper

	// InsecureSkipVerify disables TLS certificate verification.
	// Use only for testing.
	InsecureSkipVerify bool

	// MaxConcurrentRequests caps parallel in-flight API requests.
	// A value of 0 means no limit.
	MaxConcurrentRequests int

	// MandatoryRequestDelay inserts a fixed pause after every successful request.
	MandatoryRequestDelay time.Duration

	// TotalRetryDuration sets a maximum wall-clock budget for a request
	// including all retry attempts. Zero disables the budget.
	TotalRetryDuration time.Duration
}
