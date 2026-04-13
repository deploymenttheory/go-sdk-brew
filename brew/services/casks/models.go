package casks

// CaskRename describes a file/app rename operation performed during cask installation.
type CaskRename struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// CaskContainer describes a nested container within a cask artifact.
type CaskContainer struct {
	Type   string `json:"type"`
	Nested string `json:"nested"`
}

// ResourceCask represents a single Homebrew cask.
type ResourceCask struct {
	Token                        string                 `json:"token"`
	FullToken                    string                 `json:"full_token"`
	OldTokens                    []string               `json:"old_tokens"`
	Tap                          string                 `json:"tap"`
	Name                         []string               `json:"name"`
	Desc                         string                 `json:"desc"`
	Homepage                     string                 `json:"homepage"`
	URL                          string                 `json:"url"`
	URLSpecs                     map[string]interface{} `json:"url_specs"`
	Version                      string                 `json:"version"`
	Autobump                     bool                   `json:"autobump"`
	NoAutobumpMessage            *string                `json:"no_autobump_message"`
	SkipLivecheck                bool                   `json:"skip_livecheck"`
	Installed                    *string                `json:"installed"`
	InstalledTime                *string                `json:"installed_time"`
	BundleVersion                *string                `json:"bundle_version"`
	BundleShortVersion           *string                `json:"bundle_short_version"`
	Outdated                     bool                   `json:"outdated"`
	SHA256                       string                 `json:"sha256"`
	Artifacts                    []interface{}          `json:"artifacts"`
	Caveats                      *string                `json:"caveats"`
	// CaveatsRosetta is true when a Rosetta caveat applies, false/null otherwise.
	CaveatsRosetta               *bool                  `json:"caveats_rosetta"`
	DependsOn                    map[string]interface{} `json:"depends_on"`
	ConflictsWith                map[string]interface{} `json:"conflicts_with"`
	Container                    *CaskContainer         `json:"container"`
	Rename                       []CaskRename           `json:"rename"`
	AutoUpdates                  bool                   `json:"auto_updates"`
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
	TapGitHead                   string                 `json:"tap_git_head"`
	Languages                    []string               `json:"languages"`
	RubySourcePath               string                 `json:"ruby_source_path"`
	RubySourceChecksum           map[string]string      `json:"ruby_source_checksum"`
	Variations                   map[string]interface{} `json:"variations"`
	Analytics                    map[string]interface{} `json:"analytics"`
	GeneratedDate                string                 `json:"generated_date"`
}
