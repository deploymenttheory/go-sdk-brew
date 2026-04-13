package formulae

// ResourceFormula represents a single Homebrew formula.
type ResourceFormula struct {
	Name                         string                 `json:"name"`
	FullName                     string                 `json:"full_name"`
	Tap                          string                 `json:"tap"`
	OldNames                     []string               `json:"oldnames"`
	Aliases                      []string               `json:"aliases"`
	VersionedFormulae            []string               `json:"versioned_formulae"`
	Desc                         string                 `json:"desc"`
	License                      string                 `json:"license"`
	Homepage                     string                 `json:"homepage"`
	Versions                     FormulaVersions        `json:"versions"`
	URLs                         FormulaURLs            `json:"urls"`
	Revision                     int                    `json:"revision"`
	VersionScheme                int                    `json:"version_scheme"`
	CompatibilityVersion         int                    `json:"compatibility_version"`
	Autobump                     bool                   `json:"autobump"`
	NoAutobumpMessage            *string                `json:"no_autobump_message"`
	SkipLivecheck                bool                   `json:"skip_livecheck"`
	Bottle                       FormulaBottle          `json:"bottle"`
	PourBottleOnlyIf             *string                `json:"pour_bottle_only_if"`
	KegOnly                      bool                   `json:"keg_only"`
	KegOnlyReason                *FormulaKegOnlyReason  `json:"keg_only_reason"`
	Options                      []string               `json:"options"`
	BuildDependencies            []string               `json:"build_dependencies"`
	Dependencies                 []string               `json:"dependencies"`
	TestDependencies             []string               `json:"test_dependencies"`
	RecommendedDependencies      []string               `json:"recommended_dependencies"`
	OptionalDependencies         []string               `json:"optional_dependencies"`
	// UsesFromMacOS entries can be either a string (e.g. "zlib") or an object
	// (e.g. {"libarchive": "libarchive"}) depending on whether version bounds are specified.
	UsesFromMacOS                []interface{}          `json:"uses_from_macos"`
	UsesFromMacOSBounds          []FormulaUsesFromMacOSBound `json:"uses_from_macos_bounds"`
	Requirements                 []FormulaRequirement   `json:"requirements"`
	ConflictsWith                []string               `json:"conflicts_with"`
	ConflictsWithReasons         []string               `json:"conflicts_with_reasons"`
	LinkOverwrite                []string               `json:"link_overwrite"`
	Caveats                      *string                `json:"caveats"`
	Installed                    []interface{}          `json:"installed"`
	LinkedKeg                    *string                `json:"linked_keg"`
	Pinned                       bool                   `json:"pinned"`
	Outdated                     bool                   `json:"outdated"`
	Deprecated                   bool                   `json:"deprecated"`
	DeprecationDate              *string                `json:"deprecation_date"`
	DeprecationReason            *string                `json:"deprecation_reason"`
	DeprecationReplacementFormula *string               `json:"deprecation_replacement_formula"`
	DeprecationReplacementCask   *string                `json:"deprecation_replacement_cask"`
	Disabled                     bool                   `json:"disabled"`
	DisableDate                  *string                `json:"disable_date"`
	DisableReason                *string                `json:"disable_reason"`
	DisableReplacementFormula    *string                `json:"disable_replacement_formula"`
	DisableReplacementCask       *string                `json:"disable_replacement_cask"`
	PostInstallDefined           bool                   `json:"post_install_defined"`
	Service                      *FormulaService        `json:"service"`
	TapGitHead                   string                 `json:"tap_git_head"`
	RubySourcePath               string                 `json:"ruby_source_path"`
	RubySourceChecksum           FormulaChecksum        `json:"ruby_source_checksum"`
	HeadDependencies             FormulaHeadDeps        `json:"head_dependencies"`
	Variations                   map[string]interface{} `json:"variations"`
	Analytics                    FormulaAnalytics       `json:"analytics"`
	GeneratedDate                string                 `json:"generated_date"`
}

// FormulaKegOnlyReason describes why a formula is keg-only.
type FormulaKegOnlyReason struct {
	Reason      string `json:"reason"`
	Explanation string `json:"explanation"`
}

// FormulaService describes the launchd/systemd service configuration for a formula.
// Several fields are interface{} due to genuine API polymorphism:
//   - Run: string | []string | map[string][]string (platform-keyed command arrays)
//   - Name: map[string]string (platform-keyed service name, e.g. {"macos": "..."})
//   - Sockets: string | map[string]string (platform-keyed socket address)
type FormulaService struct {
	Run                  interface{}            `json:"run"`
	RunType              string                 `json:"run_type"`
	Name                 interface{}            `json:"name"`
	KeepAlive            map[string]interface{} `json:"keep_alive"`
	Cron                 string                 `json:"cron"`
	Interval             int                    `json:"interval"`
	Sockets              interface{}            `json:"sockets"`
	ProcessType          string                 `json:"process_type"`
	EnvironmentVariables map[string]string      `json:"environment_variables"`
	WorkingDir           string                 `json:"working_dir"`
	InputPath            string                 `json:"input_path"`
	LogPath              string                 `json:"log_path"`
	ErrorLogPath         string                 `json:"error_log_path"`
	MacOSLegacyTimers    bool                   `json:"macos_legacy_timers"`
	RequireRoot          bool                   `json:"require_root"`
}

// FormulaVersions holds version information for a formula.
type FormulaVersions struct {
	Stable string `json:"stable"`
	Head   string `json:"head"`
	Bottle bool   `json:"bottle"`
}

// FormulaURLs holds source URL information.
type FormulaURLs struct {
	Stable FormulaStableURL `json:"stable"`
	Head   *FormulaHeadURL  `json:"head"`
}

// FormulaStableURL describes the stable release source.
type FormulaStableURL struct {
	URL      string  `json:"url"`
	Tag      *string `json:"tag"`
	Revision *string `json:"revision"`
	Using    *string `json:"using"`
	Checksum string  `json:"checksum"`
}

// FormulaHeadURL describes the HEAD (development) source.
type FormulaHeadURL struct {
	URL    string  `json:"url"`
	Branch string  `json:"branch"`
	Using  *string `json:"using"`
}

// FormulaBottle holds bottle (pre-built binary) information.
type FormulaBottle struct {
	Stable FormulaBottleStable `json:"stable"`
}

// FormulaBottleStable holds stable bottle metadata.
type FormulaBottleStable struct {
	Rebuild int                           `json:"rebuild"`
	RootURL string                        `json:"root_url"`
	Files   map[string]FormulaBottleFile  `json:"files"`
}

// FormulaBottleFile describes a single platform bottle file.
type FormulaBottleFile struct {
	Cellar string `json:"cellar"`
	URL    string `json:"url"`
	SHA256 string `json:"sha256"`
}

// FormulaChecksum holds a checksum for the Ruby source file.
type FormulaChecksum struct {
	SHA256 string `json:"sha256"`
}

// FormulaUsesFromMacOSBound describes an optional version constraint for a uses_from_macos entry.
type FormulaUsesFromMacOSBound struct {
	Since string `json:"since,omitempty"`
}

// FormulaRequirement describes a system or architecture requirement for a formula.
type FormulaRequirement struct {
	Name     string   `json:"name"`
	Cask     *string  `json:"cask"`
	Download *string  `json:"download"`
	Version  string   `json:"version"`
	Contexts []string `json:"contexts"`
	Specs    []string `json:"specs"`
}

// FormulaHeadDeps holds HEAD-specific dependency information.
type FormulaHeadDeps struct {
	BuildDependencies       []string `json:"build_dependencies"`
	Dependencies            []string `json:"dependencies"`
	TestDependencies        []string `json:"test_dependencies"`
	RecommendedDependencies []string `json:"recommended_dependencies"`
	OptionalDependencies    []string `json:"optional_dependencies"`
	UsesFromMacOS           []interface{}               `json:"uses_from_macos"`
	UsesFromMacOSBounds     []FormulaUsesFromMacOSBound `json:"uses_from_macos_bounds"`
}

// FormulaAnalytics holds install analytics embedded in a formula's detail response.
type FormulaAnalytics struct {
	Install          FormulaAnalyticsPeriods `json:"install"`
	InstallOnRequest FormulaAnalyticsPeriods `json:"install_on_request"`
	BuildError       FormulaAnalyticsPeriods `json:"build_error"`
}

// FormulaAnalyticsPeriods holds counts for each analytics time period.
type FormulaAnalyticsPeriods struct {
	Period30d  map[string]int `json:"30d"`
	Period90d  map[string]int `json:"90d"`
	Period365d map[string]int `json:"365d"`
}
