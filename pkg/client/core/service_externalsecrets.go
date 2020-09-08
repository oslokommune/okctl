package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type externalSecretsService struct {
	api   client.ExternalSecretsAPI
	store client.ExternalSecretsStore
}

func (s *externalSecretsService) CreateExternalSecrets(_ context.Context, opts client.CreateExternalSecretsOpts) (*client.ExternalSecrets, error) {
	policy, err := s.api.CreateExternalSecretsPolicy(api.CreateExternalSecretsPolicyOpts{
		ID: opts.ID,
	})
	if err != nil {
		return nil, err
	}

	sa, err := s.api.CreateExternalSecretsServiceAccount(api.CreateExternalSecretsServiceAccountOpts{
		CreateServiceAccountOpts: api.CreateServiceAccountOpts{
			ID:        opts.ID,
			PolicyArn: policy.PolicyARN,
		},
	})
	if err != nil {
		return nil, err
	}

	chart, err := s.api.CreateExternalSecretsHelmChart(api.CreateExternalSecretsHelmChartOpts{
		ID: opts.ID,
	})
	if err != nil {
		return nil, err
	}

	externalSecrets := &client.ExternalSecrets{
		Policy:         policy,
		ServiceAccount: sa,
		Chart:          chart,
	}

	_, err = s.store.SaveExternalSecrets(externalSecrets)
	if err != nil {
		return nil, err
	}

	return externalSecrets, nil
}

// NewExternalSecretsService returns an initialised service
func NewExternalSecretsService(api client.ExternalSecretsAPI, store client.ExternalSecretsStore) client.ExternalSecretsService {
	return &externalSecretsService{
		api:   api,
		store: store,
	}
}
