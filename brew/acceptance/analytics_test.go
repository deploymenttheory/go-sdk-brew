package acceptance

// =============================================================================
// Acceptance Tests: Analytics Service
//
// Strategies used:
//   - Pattern 3 (Read-Only Information): Get analytics→verify response shape
// =============================================================================

import (
	"testing"

	"github.com/deploymenttheory/go-sdk-brew/brew/services/analytics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAcceptance_Analytics_list_by_category fetches install analytics for 30d.
func TestAcceptance_Analytics_list_by_category(t *testing.T) {
	RequireClient(t)

	ctx, cancel := NewContext()
	defer cancel()

	LogTestStage(t, "ListByCategory", "Fetching install analytics for 30d")

	result, resp, err := Client.Analytics.ListByCategory(ctx, analytics.CategoryInstall, analytics.Period30d)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)

	assert.Equal(t, 200, resp.StatusCode())
	assert.NotEmpty(t, result.Items, "analytics items should not be empty")
	assert.Positive(t, result.TotalItems)
	assert.NotEmpty(t, result.StartDate)
	assert.NotEmpty(t, result.EndDate)

	// Verify ranking is sequential.
	for i, item := range result.Items {
		assert.Equal(t, i+1, item.Number, "item %d should have number %d", i, i+1)
		assert.NotEmpty(t, item.Formula)
		assert.NotEmpty(t, item.Count)
		assert.NotEmpty(t, item.Percent)
	}

	LogTestSuccess(t, "Fetched %d analytics items for install/30d, top formula: %s",
		len(result.Items), result.Items[0].Formula)
}

// TestAcceptance_Analytics_list_core_by_category fetches core formula analytics.
// Note: this endpoint returns a different shape (formulae map, not items array).
func TestAcceptance_Analytics_list_core_by_category(t *testing.T) {
	RequireClient(t)

	ctx, cancel := NewContext()
	defer cancel()

	LogTestStage(t, "ListCoreByCategory", "Fetching homebrew-core install analytics for 30d")

	result, resp, err := Client.Analytics.ListCoreByCategory(ctx, analytics.CategoryInstall, analytics.Period30d)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)

	assert.Equal(t, 200, resp.StatusCode())
	assert.Positive(t, result.TotalItems)
	assert.NotEmpty(t, result.Formulae, "formulae map should not be empty")
	assert.NotEmpty(t, result.StartDate)

	LogTestSuccess(t, "Fetched core analytics with %d total formulae entries", len(result.Formulae))
}

// TestAcceptance_Analytics_list_cask_installs fetches cask installation analytics.
func TestAcceptance_Analytics_list_cask_installs(t *testing.T) {
	RequireClient(t)

	ctx, cancel := NewContext()
	defer cancel()

	LogTestStage(t, "ListCaskInstalls", "Fetching cask install analytics for 30d")

	result, resp, err := Client.Analytics.ListCaskInstalls(ctx, analytics.Period30d)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)

	assert.Equal(t, 200, resp.StatusCode())
	assert.NotEmpty(t, result.Formulae)

	LogTestSuccess(t, "Fetched %d cask install analytics entries", len(result.Formulae))
}
