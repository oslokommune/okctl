package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/spinner"
)

type managedPolicyService struct {
	spinner spinner.Spinner
	api     client.ManagedPolicyAPI
	store   client.ManagedPolicyStore
	report  client.ManagedPolicyReport
}

func (m *managedPolicyService) CreatePolicy(_ context.Context, opts api.CreatePolicyOpts) (*api.ManagedPolicy, error) {
	err := m.spinner.Start("managed-policy")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = m.spinner.Stop()
	}()

	p, err := m.api.CreatePolicy(opts)
	if err != nil {
		return nil, err
	}

	report, err := m.store.SaveCreatePolicy(p)
	if err != nil {
		return nil, err
	}

	err = m.report.ReportCreatePolicy(p, report)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (m *managedPolicyService) DeletePolicy(_ context.Context, opts api.DeletePolicyOpts) error {
	err := m.spinner.Start("managed-policy")
	if err != nil {
		return err
	}

	defer func() {
		_ = m.spinner.Stop()
	}()

	err = m.api.DeletePolicy(opts)
	if err != nil {
		return err
	}

	report, err := m.store.RemoveDeletePolicy(opts.StackName)
	if err != nil {
		return err
	}

	err = m.report.ReportDeletePolicy(opts.StackName, report)
	if err != nil {
		return err
	}

	return nil
}

// NewManagedPolicyService returns an initialised service
func NewManagedPolicyService(
	spinner spinner.Spinner,
	api client.ManagedPolicyAPI,
	store client.ManagedPolicyStore,
	report client.ManagedPolicyReport,
) client.ManagedPolicyService {
	return &managedPolicyService{
		spinner: spinner,
		api:     api,
		store:   store,
		report:  report,
	}
}
