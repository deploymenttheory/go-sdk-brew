package acceptance

// =============================================================================
// Acceptance Tests: Formulae Service
//
// Strategies used:
//   - Pattern 3 (Read-Only Information): List→verify fields
//   - Pattern 4 (Read-Only with Existing Data): List→GetByName on first result
// =============================================================================

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAcceptance_Formulae_list reads all formulae from the live Homebrew API
// and verifies the response shape.
func TestAcceptance_Formulae_list(t *testing.T) {
	RequireClient(t)

	ctx, cancel := NewContext()
	defer cancel()

	LogTestStage(t, "List", "Fetching all Homebrew formulae")

	result, resp, err := Client.Formulae.List(ctx)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())

	require.NotEmpty(t, result, "formulae list should not be empty")

	// Spot-check the first result has expected fields.
	first := result[0]
	assert.NotEmpty(t, first.Name, "formula name should be set")
	assert.NotEmpty(t, first.Desc, "formula description should be set")
	assert.NotEmpty(t, first.Homepage, "formula homepage should be set")
	assert.NotEmpty(t, first.Versions.Stable, "formula stable version should be set")

	LogTestSuccess(t, "Listed %d formulae, first: %s (%s)", len(result), first.Name, first.Versions.Stable)
}

// TestAcceptance_Formulae_get_by_name lists formulae then fetches the first
// result by name to verify the detail endpoint works.
func TestAcceptance_Formulae_get_by_name(t *testing.T) {
	RequireClient(t)

	ctx, cancel := NewContext()
	defer cancel()

	LogTestStage(t, "List", "Listing formulae to find a name to look up")

	all, _, err := Client.Formulae.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, all, "need at least one formula to test GetByName")

	targetName := all[0].Name
	LogTestStage(t, "GetByName", "Fetching formula %q by name", targetName)

	ctx2, cancel2 := NewContext()
	defer cancel2()

	result, resp, err := Client.Formulae.GetByName(ctx2, targetName)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)

	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, targetName, result.Name)
	assert.NotEmpty(t, result.Desc)
	assert.NotEmpty(t, result.Versions.Stable)

	LogTestSuccess(t, "Fetched formula %q version %s", result.Name, result.Versions.Stable)
}

// TestAcceptance_Formulae_get_wget fetches the well-known "wget" formula
// to verify a known formula returns expected data.
func TestAcceptance_Formulae_get_wget(t *testing.T) {
	RequireClient(t)

	ctx, cancel := NewContext()
	defer cancel()

	LogTestStage(t, "GetByName", "Fetching formula 'wget'")

	result, resp, err := Client.Formulae.GetByName(ctx, "wget")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)

	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "wget", result.Name)
	assert.Equal(t, "homebrew/core", result.Tap)
	assert.NotEmpty(t, result.Versions.Stable)
	assert.Contains(t, result.Homepage, "gnu.org")

	LogTestSuccess(t, "wget version %s verified", result.Versions.Stable)
}
