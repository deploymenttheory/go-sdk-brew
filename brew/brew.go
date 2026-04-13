package brew

import (
	"fmt"

	"github.com/deploymenttheory/go-sdk-brew/brew/client"
	"github.com/deploymenttheory/go-sdk-brew/brew/constants"
	"github.com/deploymenttheory/go-sdk-brew/brew/services/analytics"
	"github.com/deploymenttheory/go-sdk-brew/brew/services/casks"
	"github.com/deploymenttheory/go-sdk-brew/brew/services/formulae"
)

// Client is the single entry point for the Homebrew Formulae SDK.
// It wraps the HTTP transport and exposes typed service handles.
type Client struct {
	transport *client.Transport

	// Formulae provides access to Homebrew formula metadata.
	Formulae *formulae.Formulae

	// Casks provides access to Homebrew cask metadata.
	Casks *casks.Casks

	// Analytics provides access to Homebrew install analytics.
	Analytics *analytics.Analytics
}

// NewClient creates and configures a new Homebrew API client.
// baseURL defaults to https://formulae.brew.sh if empty.
// Use the WithXxx option functions to customise timeouts, logging, etc.
func NewClient(baseURL string, opts ...ClientOption) (*Client, error) {
	if baseURL == "" {
		baseURL = constants.BaseURL
	}

	transport, err := client.NewTransport(baseURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}

	return &Client{
		transport: transport,
		Formulae:  formulae.NewFormulae(transport),
		Casks:     casks.NewCasks(transport),
		Analytics: analytics.NewAnalytics(transport),
	}, nil
}
