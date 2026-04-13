package acceptance

// =============================================================================
// Acceptance Tests: Casks Service
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

// TestAcceptance_Casks_list reads all casks from the live Homebrew API.
func TestAcceptance_Casks_list(t *testing.T) {
	RequireClient(t)

	ctx, cancel := NewContext()
	defer cancel()

	LogTestStage(t, "List", "Fetching all Homebrew casks")

	result, resp, err := Client.Casks.List(ctx)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())

	require.NotEmpty(t, result, "casks list should not be empty")

	first := result[0]
	assert.NotEmpty(t, first.Token, "cask token should be set")
	assert.NotEmpty(t, first.Desc, "cask description should be set")
	assert.NotEmpty(t, first.Homepage, "cask homepage should be set")
	assert.NotEmpty(t, first.Version, "cask version should be set")

	LogTestSuccess(t, "Listed %d casks, first: %s (%s)", len(result), first.Token, first.Version)
}

// TestAcceptance_Casks_get_by_name lists casks then fetches the first result by name.
func TestAcceptance_Casks_get_by_name(t *testing.T) {
	RequireClient(t)

	ctx, cancel := NewContext()
	defer cancel()

	LogTestStage(t, "List", "Listing casks to find a token to look up")

	all, _, err := Client.Casks.List(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, all, "need at least one cask to test GetByName")

	targetToken := all[0].Token
	LogTestStage(t, "GetByName", "Fetching cask %q by name", targetToken)

	ctx2, cancel2 := NewContext()
	defer cancel2()

	result, resp, err := Client.Casks.GetByName(ctx2, targetToken)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, result)

	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, targetToken, result.Token)
	assert.NotEmpty(t, result.Desc)
	assert.NotEmpty(t, result.Version)

	LogTestSuccess(t, "Fetched cask %q version %s", result.Token, result.Version)
}
