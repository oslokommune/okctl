package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/spinner"
)

type serviceAccountService struct {
	spinner spinner.Spinner
	api     client.ServiceAccountAPI
	store   client.ServiceAccountStore
	report  client.ServiceAccountReport
}

func (m *serviceAccountService) CreateServiceAccount(_ context.Context, opts api.CreateServiceAccountOpts) (*api.ServiceAccount, error) {
	err := m.spinner.Start("service-account")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = m.spinner.Stop()
	}()

	p, err := m.api.CreateServiceAccount(opts)
	if err != nil {
		return nil, err
	}

	report, err := m.store.SaveCreateServiceAccount(p)
	if err != nil {
		return nil, err
	}

	err = m.report.ReportCreateServiceAccount(p, report)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (m *serviceAccountService) DeleteServiceAccount(_ context.Context, opts api.DeleteServiceAccountOpts) error {
	err := m.spinner.Start("service-account")
	if err != nil {
		return err
	}

	defer func() {
		_ = m.spinner.Stop()
	}()

	err = m.api.DeleteServiceAccount(opts)
	if err != nil {
		return err
	}

	report, err := m.store.RemoveDeleteServiceAccount(opts.Name)
	if err != nil {
		return err
	}

	err = m.report.ReportDeleteServiceAccount(opts.Name, report)
	if err != nil {
		return err
	}

	return nil
}

// NewServiceAccountService returns an initialised service
func NewServiceAccountService(
	spinner spinner.Spinner,
	api client.ServiceAccountAPI,
	store client.ServiceAccountStore,
	report client.ServiceAccountReport,
) client.ServiceAccountService {
	return &serviceAccountService{
		spinner: spinner,
		api:     api,
		store:   store,
		report:  report,
	}
}
