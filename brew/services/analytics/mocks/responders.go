package mocks

import (
	"github.com/deploymenttheory/go-sdk-brew/brew/mocks"
)

// AnalyticsServiceMock wraps GenericMock with analytics-specific response registration.
type AnalyticsServiceMock struct {
	*mocks.GenericMock
}

// NewAnalyticsServiceMock creates a new mock for the analytics service.
func NewAnalyticsServiceMock() *AnalyticsServiceMock {
	return &AnalyticsServiceMock{
		GenericMock: mocks.NewJSONMock("AnalyticsServiceMock"),
	}
}

// RegisterListByCategoryMock registers a response for the category analytics endpoint.
func (m *AnalyticsServiceMock) RegisterListByCategoryMock() {
	m.Register("GET", "/api/analytics/install/30d.json", 200, "validate_analytics.json")
}

// RegisterListCoreByCategoryMock registers a response for the core category analytics endpoint.
func (m *AnalyticsServiceMock) RegisterListCoreByCategoryMock() {
	m.Register("GET", "/api/analytics/install/homebrew-core/30d.json", 200, "validate_core_analytics.json")
}

// RegisterListCaskInstallsMock registers a response for the cask installs analytics endpoint.
func (m *AnalyticsServiceMock) RegisterListCaskInstallsMock() {
	m.Register("GET", "/api/analytics/cask-install/homebrew-cask/30d.json", 200, "validate_cask_analytics.json")
}

// RegisterAnalyticsErrorMock registers an error for the analytics endpoint.
func (m *AnalyticsServiceMock) RegisterAnalyticsErrorMock() {
	m.RegisterInternalError("GET", "/api/analytics/install/30d.json")
}
