package casks

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-sdk-brew/brew/client"
	"github.com/deploymenttheory/go-sdk-brew/brew/constants"
	"resty.dev/v3"
)

// Casks handles communication with the Homebrew cask endpoints.
//
// Homebrew API docs: https://formulae.brew.sh/docs/api/
type Casks struct {
	client client.Client
}

// NewCasks creates a new Casks service.
func NewCasks(c client.Client) *Casks {
	return &Casks{client: c}
}

// List returns all Homebrew casks.
// URL: GET /api/cask.json
// https://formulae.brew.sh/docs/api/
func (s *Casks) List(ctx context.Context) ([]ResourceCask, *resty.Response, error) {
	var result []ResourceCask

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointCaskList)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list casks: %w", err)
	}

	return result, resp, nil
}

// GetByName returns a single cask by token/name.
// URL: GET /api/cask/{name}.json
// https://formulae.brew.sh/docs/api/
func (s *Casks) GetByName(ctx context.Context, name string) (*ResourceCask, *resty.Response, error) {
	if name == "" {
		return nil, nil, fmt.Errorf("cask name is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointCaskByName, name)

	var result ResourceCask

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetResult(&result).
		Get(endpoint)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get cask %q: %w", name, err)
	}

	return &result, resp, nil
}
