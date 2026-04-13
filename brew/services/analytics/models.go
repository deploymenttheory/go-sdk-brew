package analytics

// AnalyticsCategory represents a valid analytics category for the Homebrew API.
type AnalyticsCategory string

const (
	// CategoryInstall counts all installs.
	CategoryInstall AnalyticsCategory = "install"

	// CategoryInstallOnRequest counts installs explicitly requested by users
	// (excludes installs as dependencies).
	CategoryInstallOnRequest AnalyticsCategory = "install-on-request"

	// CategoryBuildError counts build errors.
	CategoryBuildError AnalyticsCategory = "build-error"
)

// AnalyticsPeriod represents a valid time period for analytics queries.
type AnalyticsPeriod string

const (
	// Period30d is the last 30 days.
	Period30d AnalyticsPeriod = "30d"

	// Period90d is the last 90 days.
	Period90d AnalyticsPeriod = "90d"

	// Period365d is the last 365 days.
	Period365d AnalyticsPeriod = "365d"
)

// ResourceAnalytics represents the response from a Homebrew analytics endpoint
// for the general (non-core) analytics endpoints. Items are ranked with a number,
// count, and percent.
type ResourceAnalytics struct {
	Category   string          `json:"category"`
	TotalItems int             `json:"total_items"`
	StartDate  string          `json:"start_date"`
	EndDate    string          `json:"end_date"`
	TotalCount int             `json:"total_count"`
	Items      []AnalyticsItem `json:"items"`
}

// AnalyticsItem represents a single formula or cask entry in a ranked analytics response.
type AnalyticsItem struct {
	Number  int    `json:"number"`
	Formula string `json:"formula"`
	Count   string `json:"count"`
	Percent string `json:"percent"`
}

// ResourceCoreAnalytics represents the response from the homebrew-core analytics
// endpoint. The formulae field is a map of formula name to a list of version entries.
// This is a different shape from ResourceAnalytics.
//
// Used by: ListCoreByCategory
type ResourceCoreAnalytics struct {
	Category   string                       `json:"category"`
	TotalItems int                          `json:"total_items"`
	StartDate  string                       `json:"start_date"`
	EndDate    string                       `json:"end_date"`
	TotalCount int                          `json:"total_count"`
	Formulae   map[string][]CoreFormulaItem `json:"formulae"`
}

// CoreFormulaItem represents a single install count entry within the core analytics response.
type CoreFormulaItem struct {
	Formula string `json:"formula"`
	Count   string `json:"count"`
}

// ResourceCaskAnalytics represents the response from the cask install analytics
// endpoint. Similar to ResourceCoreAnalytics, the formulae field is a map of
// cask token to a list of install count entries.
//
// Used by: ListCaskInstalls
type ResourceCaskAnalytics struct {
	Category   string                     `json:"category"`
	TotalItems int                        `json:"total_items"`
	StartDate  string                     `json:"start_date"`
	EndDate    string                     `json:"end_date"`
	TotalCount int                        `json:"total_count"`
	Formulae   map[string][]CaskInstallItem `json:"formulae"`
}

// CaskInstallItem represents a single install count entry in the cask analytics response.
type CaskInstallItem struct {
	Cask  string `json:"cask"`
	Count string `json:"count"`
}
