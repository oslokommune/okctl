package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type externalSecretsService struct {
	spinner spinner.Spinner
	api     client.ExternalSecretsAPI
	store   client.ExternalSecretsStore
	report  client.ExternalSecretsReport
}

func (s *externalSecretsService) DeleteExternalSecrets(_ context.Context, id api.ID) error {
	err := s.spinner.Start("external-secrets")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = s.api.DeleteExternalSecretsPolicy(id)
	if err != nil {
		return err
	}

	err = s.api.DeleteExternalSecretsServiceAccount(id)
	if err != nil {
		return err
	}

	report, err := s.store.RemoveExternalSecrets(id)
	if err != nil {
		return err
	}

	err = s.report.ReportDeleteExternalSecrets(report)
	if err != nil {
		return err
	}

	return nil
}

func (s *externalSecretsService) CreateExternalSecrets(_ context.Context, opts client.CreateExternalSecretsOpts) (*client.ExternalSecrets, error) {
	err := s.spinner.Start("external-secrets")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

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

	report, err := s.store.SaveExternalSecrets(externalSecrets)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreateExternalSecrets(externalSecrets, report)
	if err != nil {
		return nil, err
	}

	return externalSecrets, nil
}

// NewExternalSecretsService returns an initialised service
func NewExternalSecretsService(
	spinner spinner.Spinner,
	api client.ExternalSecretsAPI,
	store client.ExternalSecretsStore,
	report client.ExternalSecretsReport,
) client.ExternalSecretsService {
	return &externalSecretsService{
		spinner: spinner,
		api:     api,
		store:   store,
		report:  report,
	}
}
