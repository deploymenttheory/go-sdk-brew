package client

import (
	"context"

	"go.uber.org/zap"
)

// Client is the interface service implementations depend on.
// The Transport struct in this package satisfies this interface.
type Client interface {
	// NewRequest returns a RequestBuilder that the service layer uses to
	// construct a complete request — headers, query params, result target —
	// before executing it via Get/GetBytes.
	// Retry, throttling, and concurrency limiting are applied by the transport
	// at execution time.
	NewRequest(ctx context.Context) *RequestBuilder

	// GetLogger returns the configured zap logger instance.
	GetLogger() *zap.Logger
}
