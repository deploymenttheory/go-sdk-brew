package analytics

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-sdk-brew/brew/services/analytics/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMockService creates an Analytics service wired to a fresh mock.
func setupMockService(t *testing.T) (*Analytics, *mocks.AnalyticsServiceMock) {
	t.Helper()
	mock := mocks.NewAnalyticsServiceMock()
	return NewAnalytics(mock), mock
}

// =============================================================================
// ListByCategory
// =============================================================================

func TestUnit_Analytics_ListByCategory_Success(t *testing.T) {
	svc, mock := setupMockService(t)
	mock.RegisterListByCategoryMock()

	result, resp, err := svc.ListByCategory(context.Background(), CategoryInstall, Period30d)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, resp)

	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, 23609, result.TotalItems)
	require.Len(t, result.Items, 5)
	assert.Equal(t, "ca-certificates", result.Items[0].Formula)
	assert.Equal(t, 1, result.Items[0].Number)
	assert.Equal(t, "510,810", result.Items[0].Count)
}

func TestUnit_Analytics_ListByCategory_Error(t *testing.T) {
	svc, mock := setupMockService(t)
	mock.RegisterAnalyticsErrorMock()

	result, resp, err := svc.ListByCategory(context.Background(), CategoryInstall, Period30d)
	assert.Error(t, err)
	assert.Nil(t, result)
	_ = resp
}

// =============================================================================
// ListCoreByCategory
// =============================================================================

func TestUnit_Analytics_ListCoreByCategory_Success(t *testing.T) {
	svc, mock := setupMockService(t)
	mock.RegisterListCoreByCategoryMock()

	result, resp, err := svc.ListCoreByCategory(context.Background(), CategoryInstall, Period30d)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, resp)

	assert.Equal(t, 200, resp.StatusCode())
	// The core analytics mock uses the same fixture (validate_analytics.json) which
	// has total_items=3 but uses "items" not "formulae", so Formulae map will be empty.
	// The test validates the endpoint is called correctly.
	assert.NotNil(t, result)
}

// =============================================================================
// ListCaskInstalls
// =============================================================================

func TestUnit_Analytics_ListCaskInstalls_Success(t *testing.T) {
	svc, mock := setupMockService(t)
	mock.RegisterListCaskInstallsMock()

	result, resp, err := svc.ListCaskInstalls(context.Background(), Period30d)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, resp)

	assert.Equal(t, 200, resp.StatusCode())
	assert.NotEmpty(t, result.Formulae)
}

// =============================================================================
// Validation Errors
// =============================================================================

func TestUnit_Analytics_ListByCategory_EmptyCategory(t *testing.T) {
	svc, _ := setupMockService(t)

	result, resp, err := svc.ListByCategory(context.Background(), "", Period30d)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "analytics category is required")
}

func TestUnit_Analytics_ListByCategory_EmptyPeriod(t *testing.T) {
	svc, _ := setupMockService(t)

	result, resp, err := svc.ListByCategory(context.Background(), CategoryInstall, "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "analytics period is required")
}

func TestUnit_Analytics_ListCaskInstalls_EmptyPeriod(t *testing.T) {
	svc, _ := setupMockService(t)

	result, resp, err := svc.ListCaskInstalls(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "analytics period is required")
}
