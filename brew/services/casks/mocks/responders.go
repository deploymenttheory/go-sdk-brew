package mocks

import (
	"github.com/deploymenttheory/go-sdk-brew/brew/mocks"
)

// CasksServiceMock wraps GenericMock with cask-specific response registration.
type CasksServiceMock struct {
	*mocks.GenericMock
}

// NewCasksServiceMock creates a new mock for the casks service.
func NewCasksServiceMock() *CasksServiceMock {
	return &CasksServiceMock{
		GenericMock: mocks.NewJSONMock("CasksServiceMock"),
	}
}

// RegisterMocks registers all standard success responses.
func (m *CasksServiceMock) RegisterMocks() {
	m.RegisterListCasksMock()
	m.RegisterGetCaskMock()
}

// RegisterListCasksMock registers the List all casks response.
func (m *CasksServiceMock) RegisterListCasksMock() {
	m.Register("GET", "/api/cask.json", 200, "validate_list_casks.json")
}

// RegisterGetCaskMock registers a Get cask by name response.
func (m *CasksServiceMock) RegisterGetCaskMock() {
	m.Register("GET", "/api/cask/iterm2.json", 200, "validate_get_cask.json")
}

// RegisterNotFoundCaskMock registers a 404 for an unknown cask.
func (m *CasksServiceMock) RegisterNotFoundCaskMock() {
	m.RegisterNotFoundError("GET", "/api/cask/nonexistent.json")
}

// RegisterListErrorMock registers an error for the list endpoint.
func (m *CasksServiceMock) RegisterListErrorMock() {
	m.RegisterInternalError("GET", "/api/cask.json")
}
