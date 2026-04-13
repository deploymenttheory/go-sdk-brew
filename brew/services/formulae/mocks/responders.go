package mocks

import (
	"github.com/deploymenttheory/go-sdk-brew/brew/mocks"
)

// FormulaeServiceMock wraps GenericMock with formulae-specific response registration.
type FormulaeServiceMock struct {
	*mocks.GenericMock
}

// NewFormulaeServiceMock creates a new mock for the formulae service.
func NewFormulaeServiceMock() *FormulaeServiceMock {
	return &FormulaeServiceMock{
		GenericMock: mocks.NewJSONMock("FormulaeServiceMock"),
	}
}

// RegisterMocks registers all standard success responses.
func (m *FormulaeServiceMock) RegisterMocks() {
	m.RegisterListFormulaeMock()
	m.RegisterGetFormulaMock()
}

// RegisterListFormulaeMock registers the List all formulae response.
func (m *FormulaeServiceMock) RegisterListFormulaeMock() {
	m.Register("GET", "/api/formula.json", 200, "validate_list_formulae.json")
}

// RegisterGetFormulaMock registers a Get formula by name response.
func (m *FormulaeServiceMock) RegisterGetFormulaMock() {
	m.Register("GET", "/api/formula/wget.json", 200, "validate_get_formula.json")
}

// RegisterNotFoundFormulaMock registers a 404 for an unknown formula.
func (m *FormulaeServiceMock) RegisterNotFoundFormulaMock() {
	m.RegisterNotFoundError("GET", "/api/formula/nonexistent.json")
}

// RegisterListErrorMock registers an error for the list endpoint.
func (m *FormulaeServiceMock) RegisterListErrorMock() {
	m.RegisterInternalError("GET", "/api/formula.json")
}
