package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type parameterService struct {
	api   client.ParameterAPI
	state client.ParameterState
}

func (s *parameterService) DeleteSecret(_ context.Context, opts client.DeleteSecretOpts) error {
	err := s.api.DeleteSecret(api.DeleteSecretOpts{
		Name: opts.Name,
	})
	if err != nil {
		return err
	}

	return s.state.RemoveSecret(opts.Name)
}

func (s *parameterService) CreateSecret(_ context.Context, opts client.CreateSecretOpts) (*client.SecretParameter, error) {
	secret, err := s.api.CreateSecret(api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   opts.Name,
		Secret: opts.Secret,
	})
	if err != nil {
		return nil, err
	}

	sec := &client.SecretParameter{
		ID:      secret.ID,
		Name:    secret.Name,
		Path:    secret.Path,
		Version: secret.Version,
		Content: secret.Content,
	}

	err = s.state.SaveSecret(sec)
	if err != nil {
		return nil, err
	}

	return sec, nil
}

// NewParameterService returns an initialised service
func NewParameterService(
	api client.ParameterAPI,
	state client.ParameterState,
) client.ParameterService {
	return &parameterService{
		api:   api,
		state: state,
	}
}
