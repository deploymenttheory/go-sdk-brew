package casks

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-sdk-brew/brew/services/casks/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMockService creates a Casks service wired to a fresh mock.
func setupMockService(t *testing.T) (*Casks, *mocks.CasksServiceMock) {
	t.Helper()
	mock := mocks.NewCasksServiceMock()
	return NewCasks(mock), mock
}

// =============================================================================
// List
// =============================================================================

func TestUnit_Casks_List_Success(t *testing.T) {
	svc, mock := setupMockService(t)
	mock.RegisterListCasksMock()

	result, resp, err := svc.List(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, resp)

	assert.Equal(t, 200, resp.StatusCode())
	require.Len(t, result, 2)
	assert.Equal(t, "iterm2", result[0].Token)
	assert.Equal(t, "visual-studio-code", result[1].Token)
}

func TestUnit_Casks_List_Error(t *testing.T) {
	svc, mock := setupMockService(t)
	mock.RegisterListErrorMock()

	result, resp, err := svc.List(context.Background())
	assert.Error(t, err)
	assert.Nil(t, result)
	_ = resp
}

// =============================================================================
// GetByName
// =============================================================================

func TestUnit_Casks_GetByName_Success(t *testing.T) {
	svc, mock := setupMockService(t)
	mock.RegisterGetCaskMock()

	result, resp, err := svc.GetByName(context.Background(), "iterm2")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, resp)

	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "iterm2", result.Token)
	assert.Equal(t, "Terminal emulator as alternative to Apple's Terminal app", result.Desc)
	assert.Equal(t, "3.6.9", result.Version)
	assert.True(t, result.AutoUpdates)
	assert.Equal(t, "27c00f476978c0a243144a0e03f01345facd0812ce4112fabaa54168c050b19e", result.SHA256)
}

func TestUnit_Casks_GetByName_NotFound(t *testing.T) {
	svc, mock := setupMockService(t)
	mock.RegisterNotFoundCaskMock()

	result, resp, err := svc.GetByName(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	_ = resp
}

// =============================================================================
// Validation Errors
// =============================================================================

func TestUnit_Casks_GetByName_EmptyName(t *testing.T) {
	svc, _ := setupMockService(t)

	result, resp, err := svc.GetByName(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "cask name is required")
}
