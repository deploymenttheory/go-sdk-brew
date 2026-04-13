package analytics

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-sdk-brew/brew/client"
	"github.com/deploymenttheory/go-sdk-brew/brew/constants"
	"resty.dev/v3"
)

// Analytics handles communication with the Homebrew analytics endpoints.
//
// Homebrew API docs: https://formulae.brew.sh/docs/api/
type Analytics struct {
	client client.Client
}

// NewAnalytics creates a new Analytics service.
func NewAnalytics(c client.Client) *Analytics {
	return &Analytics{client: c}
}

// ListByCategory returns formula analytics for a given category and time period.
// URL: GET /api/analytics/{category}/{days}.json
// Valid categories: "install", "install-on-request", "build-error".
// Valid periods: "30d", "90d", "365d".
// https://formulae.brew.sh/docs/api/
func (s *Analytics) ListByCategory(ctx context.Context, category AnalyticsCategory, period AnalyticsPeriod) (*ResourceAnalytics, *resty.Response, error) {
	if category == "" {
		return nil, nil, fmt.Errorf("analytics category is required")
	}
	if period == "" {
		return nil, nil, fmt.Errorf("analytics period is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointAnalyticsByCategory, category, period)

	var result ResourceAnalytics

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetResult(&result).
		Get(endpoint)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get analytics for category %q period %q: %w", category, period, err)
	}

	return &result, resp, nil
}

// ListCoreByCategory returns analytics for homebrew-core formulae for a given category and period.
// URL: GET /api/analytics/{category}/homebrew-core/{days}.json
//
// Note: This endpoint returns a different response shape from ListByCategory.
// The formulae field is a map of formula name → list of version entries.
// https://formulae.brew.sh/docs/api/
func (s *Analytics) ListCoreByCategory(ctx context.Context, category AnalyticsCategory, period AnalyticsPeriod) (*ResourceCoreAnalytics, *resty.Response, error) {
	if category == "" {
		return nil, nil, fmt.Errorf("analytics category is required")
	}
	if period == "" {
		return nil, nil, fmt.Errorf("analytics period is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointAnalyticsCoreByCategory, category, period)

	var result ResourceCoreAnalytics

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetResult(&result).
		Get(endpoint)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get core analytics for category %q period %q: %w", category, period, err)
	}

	return &result, resp, nil
}

// ListCaskInstalls returns cask installation analytics for a given time period.
// URL: GET /api/analytics/cask-install/homebrew-cask/{days}.json
//
// Note: This endpoint returns a different response shape from ListByCategory.
// The formulae field is a map of cask token → list of install count entries.
// https://formulae.brew.sh/docs/api/
func (s *Analytics) ListCaskInstalls(ctx context.Context, period AnalyticsPeriod) (*ResourceCaskAnalytics, *resty.Response, error) {
	if period == "" {
		return nil, nil, fmt.Errorf("analytics period is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointAnalyticsCaskInstalls, period)

	var result ResourceCaskAnalytics

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetResult(&result).
		Get(endpoint)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get cask install analytics for period %q: %w", period, err)
	}

	return &result, resp, nil
}
