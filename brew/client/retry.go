package client

import (
	"resty.dev/v3"
)

// retryCondition is the resty AddRetryConditions callback.
// It returns true when the request should be retried.
//
// Resty handles the actual retry scheduling: it uses capped exponential backoff
// with full jitter (min * 2^attempt, capped at max, then halved and jittered).
// The min and max bounds are set by RetryWaitTime and RetryMaxWaitTime on the
// resty client. This function only decides whether a given response warrants
// another attempt.
//
// Since the Homebrew API is entirely GET-based (all requests are idempotent),
// retry is safe for all methods. Retries are applied on transient server errors.
func retryCondition(resp *resty.Response, err error) bool {
	// Network / transport error with no response — always safe to retry GET.
	if err != nil {
		return true
	}

	if resp == nil {
		return false
	}

	code := resp.StatusCode()

	// Never retry definitive client-side failures.
	if isNonRetryableStatusCode(code) {
		return false
	}

	return isTransientStatusCode(code)
}

// isTransientStatusCode returns true for errors that are likely temporary and
// worth retrying with exponential backoff.
func isTransientStatusCode(code int) bool {
	switch code {
	case 408, 500, 502, 503, 504:
		return true
	default:
		return false
	}
}

// isNonRetryableStatusCode returns true for definitive client-side errors
// that will not succeed on retry regardless of timing.
func isNonRetryableStatusCode(code int) bool {
	switch code {
	case 400, 401, 402, 403, 404, 405, 406, 407, 409, 410,
		411, 412, 413, 414, 415, 416, 417, 422, 423, 424,
		426, 428, 429, 431, 451:
		return true
	default:
		return false
	}
}
