package acceptance

import (
	"fmt"
	"os"
	"strconv"
	"time"

	brew "github.com/deploymenttheory/go-sdk-brew/brew"
	"github.com/deploymenttheory/go-sdk-brew/brew/constants"
)

// TestConfig holds configuration for acceptance tests driven by environment variables.
type TestConfig struct {
	// BaseURL is the Homebrew API base URL (default: https://formulae.brew.sh).
	BaseURL string

	// RequestTimeout is the per-test context timeout (default: 30s).
	RequestTimeout time.Duration

	// Verbose enables detailed test logging when true.
	Verbose bool
}

var (
	// Config is the global acceptance test configuration.
	Config *TestConfig
	// Client is the shared Homebrew SDK client for acceptance tests.
	Client *brew.Client
)

func init() {
	Config = &TestConfig{
		BaseURL:        getEnv("BREW_BASE_URL", constants.BaseURL),
		RequestTimeout: getDurationEnv("BREW_REQUEST_TIMEOUT", 30*time.Second),
		Verbose:        getBoolEnv("BREW_VERBOSE", false),
	}
}

// InitClient creates the shared Homebrew client from environment variables.
func InitClient() error {
	var err error
	Client, err = brew.NewClient(Config.BaseURL,
		brew.WithTimeout(getDurationEnv("BREW_TRANSPORT_TIMEOUT", 2*time.Minute)),
	)
	if err != nil {
		return fmt.Errorf("failed to create Homebrew client: %w", err)
	}

	if Config.Verbose {
		fmt.Printf("Acceptance test client initialised: %s\n", Config.BaseURL)
	}
	return nil
}

// IsConfigured always returns true since the Homebrew API requires no credentials.
func IsConfigured() bool {
	return true
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
