package constants

// BaseURL is the default base URL for the Homebrew Formulae API.
const BaseURL = "https://formulae.brew.sh"

// Formulae endpoints.
const (
	// EndpointFormulaList returns all formulae.
	// GET /api/formula.json
	EndpointFormulaList = "/api/formula.json"

	// EndpointFormulaByName returns a single formula by name.
	// GET /api/formula/{name}.json
	// Use fmt.Sprintf(EndpointFormulaByName, name).
	EndpointFormulaByName = "/api/formula/%s.json"
)

// Cask endpoints.
const (
	// EndpointCaskList returns all casks.
	// GET /api/cask.json
	EndpointCaskList = "/api/cask.json"

	// EndpointCaskByName returns a single cask by name.
	// GET /api/cask/{name}.json
	// Use fmt.Sprintf(EndpointCaskByName, name).
	EndpointCaskByName = "/api/cask/%s.json"
)

// Analytics endpoints.
const (
	// EndpointAnalyticsByCategory returns analytics for a category and time period.
	// GET /api/analytics/{category}/{days}.json
	// Use fmt.Sprintf(EndpointAnalyticsByCategory, category, days).
	EndpointAnalyticsByCategory = "/api/analytics/%s/%s.json"

	// EndpointAnalyticsCoreByCategory returns core formulae analytics.
	// GET /api/analytics/{category}/homebrew-core/{days}.json
	// Use fmt.Sprintf(EndpointAnalyticsCoreByCategory, category, days).
	EndpointAnalyticsCoreByCategory = "/api/analytics/%s/homebrew-core/%s.json"

	// EndpointAnalyticsCaskInstalls returns cask installation analytics.
	// GET /api/analytics/cask-install/homebrew-cask/{days}.json
	// Use fmt.Sprintf(EndpointAnalyticsCaskInstalls, days).
	EndpointAnalyticsCaskInstalls = "/api/analytics/cask-install/homebrew-cask/%s.json"
)
