package formulae

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-sdk-brew/brew/client"
	"github.com/deploymenttheory/go-sdk-brew/brew/constants"
	"resty.dev/v3"
)

// Formulae handles communication with the Homebrew formula endpoints.
//
// Homebrew API docs: https://formulae.brew.sh/docs/api/
type Formulae struct {
	client client.Client
}

// NewFormulae creates a new Formulae service.
func NewFormulae(c client.Client) *Formulae {
	return &Formulae{client: c}
}

// List returns all Homebrew formulae.
// URL: GET /api/formula.json
// https://formulae.brew.sh/docs/api/
func (s *Formulae) List(ctx context.Context) ([]ResourceFormula, *resty.Response, error) {
	var result []ResourceFormula

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointFormulaList)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to list formulae: %w", err)
	}

	return result, resp, nil
}

// GetByName returns a single formula by name.
// URL: GET /api/formula/{name}.json
// https://formulae.brew.sh/docs/api/
func (s *Formulae) GetByName(ctx context.Context, name string) (*ResourceFormula, *resty.Response, error) {
	if name == "" {
		return nil, nil, fmt.Errorf("formula name is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointFormulaByName, name)

	var result ResourceFormula

	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetResult(&result).
		Get(endpoint)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get formula %q: %w", name, err)
	}

	return &result, resp, nil
}
