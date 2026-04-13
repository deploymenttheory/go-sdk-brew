package formulae

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-sdk-brew/brew/services/formulae/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMockService creates a Formulae service wired to a fresh mock.
func setupMockService(t *testing.T) (*Formulae, *mocks.FormulaeServiceMock) {
	t.Helper()
	mock := mocks.NewFormulaeServiceMock()
	return NewFormulae(mock), mock
}

// =============================================================================
// List
// =============================================================================

func TestUnit_Formulae_List_Success(t *testing.T) {
	svc, mock := setupMockService(t)
	mock.RegisterListFormulaeMock()

	result, resp, err := svc.List(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, resp)

	assert.Equal(t, 200, resp.StatusCode())
	require.Len(t, result, 2)
	assert.Equal(t, "wget", result[0].Name)
	assert.Equal(t, "Internet file retriever", result[0].Desc)
	assert.Equal(t, "curl", result[1].Name)
}

func TestUnit_Formulae_List_Error(t *testing.T) {
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

func TestUnit_Formulae_GetByName_Success(t *testing.T) {
	svc, mock := setupMockService(t)
	mock.RegisterGetFormulaMock()

	result, resp, err := svc.GetByName(context.Background(), "wget")
	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, resp)

	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "wget", result.Name)
	assert.Equal(t, "Internet file retriever", result.Desc)
	assert.Equal(t, "GPL-3.0-or-later", result.License)
	assert.Equal(t, "1.25.0", result.Versions.Stable)
	assert.Contains(t, result.Dependencies, "openssl@3")
	assert.Equal(t, "homebrew/core", result.Tap)
	assert.Equal(t, "Formula/w/wget.rb", result.RubySourcePath)
}

func TestUnit_Formulae_GetByName_NotFound(t *testing.T) {
	svc, mock := setupMockService(t)
	mock.RegisterNotFoundFormulaMock()

	result, resp, err := svc.GetByName(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, result)
	_ = resp
}

// =============================================================================
// Validation Errors
// =============================================================================

func TestUnit_Formulae_GetByName_EmptyName(t *testing.T) {
	svc, _ := setupMockService(t)

	result, resp, err := svc.GetByName(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "formula name is required")
}
