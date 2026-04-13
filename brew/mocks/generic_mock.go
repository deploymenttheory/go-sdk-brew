package mocks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/deploymenttheory/go-sdk-brew/brew/client"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// registeredResponse holds a pre-canned response for a single endpoint.
type registeredResponse struct {
	statusCode int
	rawBody    []byte
	errMsg     string
}

// GenericMock is a reusable test double implementing client.Client.
type GenericMock struct {
	name           string
	responses      map[string]registeredResponse
	logger         *zap.Logger
	LastQueryParams map[string]string
	fixtureDir     string
}

// GenericMockConfig configures a GenericMock instance.
type GenericMockConfig struct {
	Name       string
	FixtureDir string
}

// NewGenericMock creates a new generic mock with the specified configuration.
func NewGenericMock(config GenericMockConfig) *GenericMock {
	if config.Name == "" {
		config.Name = "GenericMock"
	}
	if config.FixtureDir == "" {
		for i := 1; i < 10; i++ {
			_, filename, _, ok := runtime.Caller(i)
			if !ok {
				break
			}
			dir := filepath.Dir(filename)
			if filepath.Base(dir) == "mocks" {
				continue
			}
			config.FixtureDir = filepath.Join(dir, "mocks")
			break
		}
	}

	return &GenericMock{
		name:       config.Name,
		responses:  make(map[string]registeredResponse),
		logger:     zap.NewNop(),
		fixtureDir: config.FixtureDir,
	}
}

// NewJSONMock creates a mock configured for JSON responses.
func NewJSONMock(name string) *GenericMock {
	return NewGenericMock(GenericMockConfig{Name: name})
}

// Register registers a mock response for the given method and path.
func (m *GenericMock) Register(method, path string, statusCode int, fixture string) {
	var body []byte
	if fixture != "" {
		data, err := m.loadFixture(fixture)
		if err != nil {
			panic(fmt.Sprintf("%s: failed to load fixture %q: %v", m.name, fixture, err))
		}
		body = data
	}
	m.responses[method+":"+path] = registeredResponse{statusCode: statusCode, rawBody: body}
}

// RegisterError registers a mock error response.
func (m *GenericMock) RegisterError(method, path string, statusCode int, fixture string, errMsg string) {
	var body []byte
	if fixture != "" {
		var err error
		body, err = m.loadFixture(fixture)
		if err != nil {
			panic(fmt.Sprintf("%s: failed to load error fixture %q: %v", m.name, fixture, err))
		}

		if errMsg == "" {
			var parsed struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}
			if json.Unmarshal(body, &parsed) == nil {
				errMsg = fmt.Sprintf("Homebrew API error (%d) [%s]: %s", statusCode, parsed.Code, parsed.Message)
			}
		}
	}

	if errMsg == "" {
		errMsg = fmt.Sprintf("%s: error response %d", m.name, statusCode)
	}

	m.responses[method+":"+path] = registeredResponse{statusCode: statusCode, rawBody: body, errMsg: errMsg}
}

// RegisterRawBody registers a mock response with raw body bytes.
func (m *GenericMock) RegisterRawBody(method, path string, statusCode int, body []byte) {
	m.responses[method+":"+path] = registeredResponse{statusCode: statusCode, rawBody: body}
}

// loadFixture loads a fixture file from the configured fixture directory.
// Falls back to centralized test_fixtures directory for common errors.
func (m *GenericMock) loadFixture(filename string) ([]byte, error) {
	path := filepath.Join(m.fixtureDir, filename)
	data, err := os.ReadFile(path)
	if err == nil {
		return data, nil
	}

	if isCommonErrorFixture(filename) {
		centralPath := filepath.Join(filepath.Dir(m.fixtureDir), "..", "..", "mocks", "test_fixtures", filename)
		data, err = os.ReadFile(centralPath)
		if err == nil {
			return data, nil
		}
	}

	return nil, fmt.Errorf("read fixture %s: %w", filename, err)
}

// isCommonErrorFixture checks if a filename is a common error fixture.
func isCommonErrorFixture(filename string) bool {
	commonErrors := []string{
		"error_not_found.json",
		"error_bad_request.json",
	}
	return slices.Contains(commonErrors, filename)
}

// dispatch is the core routing logic for all HTTP methods.
func (m *GenericMock) dispatch(method, path string, result any) (*resty.Response, error) {
	r, ok := m.responses[method+":"+path]
	if !ok {
		return nil, fmt.Errorf("%s: no response registered for %s %s", m.name, method, path)
	}

	headers := http.Header{"Content-Type": {"application/json"}}
	resp := NewMockResponse(r.statusCode, headers, r.rawBody)

	if r.errMsg != "" {
		return resp, fmt.Errorf("%s", r.errMsg)
	}

	if result != nil && len(r.rawBody) > 0 {
		if byteSlicePtr, ok := result.(*[]byte); ok {
			*byteSlicePtr = r.rawBody
		} else {
			if err := json.Unmarshal(r.rawBody, result); err != nil {
				return resp, fmt.Errorf("%s: unmarshal into result: %w", m.name, err)
			}
		}
	}

	return resp, nil
}

// -----------------------------------------------------------------------------
// client.Client Interface Implementation
// -----------------------------------------------------------------------------

func (m *GenericMock) NewRequest(ctx context.Context) *client.RequestBuilder {
	return client.NewMockRequestBuilderWithQueryCapture(ctx, func(method, path string, result any) (*resty.Response, error) {
		return m.dispatch(method, path, result)
	}, &m.LastQueryParams)
}

func (m *GenericMock) GetLogger() *zap.Logger { return m.logger }

// Convenience methods for registering common error responses.

// RegisterNotFoundError registers a 404 Not Found error for the given method and path.
func (m *GenericMock) RegisterNotFoundError(method, path string) {
	m.RegisterError(method, path, http.StatusNotFound, "error_not_found.json", "")
}

// RegisterBadRequestError registers a 400 Bad Request error for the given method and path.
func (m *GenericMock) RegisterBadRequestError(method, path string) {
	m.RegisterError(method, path, http.StatusBadRequest, "error_bad_request.json", "")
}

// RegisterInternalError registers a 500 Internal Server Error for the given method and path.
func (m *GenericMock) RegisterInternalError(method, path string) {
	m.RegisterError(method, path, http.StatusInternalServerError, "", fmt.Sprintf("%s: internal server error", m.name))
}
