package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type parameterService struct {
	service api.ParameterService
	state   client.ParameterState
}

func (s *parameterService) DeleteSecret(context context.Context, opts client.DeleteSecretOpts) error {
	err := s.service.DeleteSecret(context, api.DeleteSecretOpts{
		ID:   opts.ID,
		Name: opts.Name,
	})
	if err != nil {
		return err
	}

	return s.state.RemoveSecret(opts.Name)
}

func (s *parameterService) CreateSecret(context context.Context, opts client.CreateSecretOpts) (*client.SecretParameter, error) {
	secret, err := s.service.CreateSecret(context, api.CreateSecretOpts{
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
	service api.ParameterService,
	state client.ParameterState,
) client.ParameterService {
	return &parameterService{
		service: service,
		state:   state,
	}
}
