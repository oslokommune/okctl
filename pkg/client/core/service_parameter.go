package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type parameterService struct {
	api    client.ParameterAPI
	store  client.ParameterStore
	report client.ParameterReport
}

func (s *parameterService) CreateSecret(_ context.Context, opts api.CreateSecretOpts) (*api.SecretParameter, error) {
	secret, err := s.api.CreateSecret(opts)
	if err != nil {
		return nil, err
	}

	report, err := s.store.SaveSecret(secret)
	if err != nil {
		return nil, err
	}

	err = s.report.SaveSecret(secret, report)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// NewParameterService returns an initialised service
func NewParameterService(api client.ParameterAPI, store client.ParameterStore, report client.ParameterReport) client.ParameterService {
	return &parameterService{
		api:    api,
		store:  store,
		report: report,
	}
}
