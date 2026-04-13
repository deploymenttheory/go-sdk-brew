package client

import "time"

const (
	UserAgentBase = "go-sdk-brew"
)

// HTTP client defaults.
const (
	DefaultTimeout   = 30 * time.Second
	MaxRetries       = 3
	RetryWaitTime    = 2 * time.Second
	RetryMaxWaitTime = 30 * time.Second

	// DefaultMaxConcurrentRequests caps parallel in-flight requests.
	// Set to 0 to use WithMaxConcurrentRequests.
	DefaultMaxConcurrentRequests = 10

	// adaptiveDelayMax is the ceiling applied to the adaptive inter-request
	// delay computed from response-time EMA tracking.
	adaptiveDelayMax = 5 * time.Second
)
