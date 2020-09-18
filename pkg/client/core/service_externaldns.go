package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type externalDNSService struct {
	spinner spinner.Spinner
	api     client.ExternalDNSAPI
	store   client.ExternalDNSStore
	report  client.ExternalDNSReport
}

func (s *externalDNSService) DeleteExternalDNS(_ context.Context, id api.ID) error {
	err := s.spinner.Start("external-dns")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = s.api.DeleteExternalDNSPolicy(id)
	if err != nil {
		return err
	}

	err = s.api.DeleteExternalDNSServiceAccount(id)
	if err != nil {
		return err
	}

	report, err := s.store.RemoveExternalDNS(id)
	if err != nil {
		return err
	}

	err = s.report.ReportDeleteExternalDNS(report)
	if err != nil {
		return err
	}

	return nil
}

func (s *externalDNSService) CreateExternalDNS(_ context.Context, opts client.CreateExternalDNSOpts) (*client.ExternalDNS, error) {
	err := s.spinner.Start("external-dns")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	policy, err := s.api.CreateExternalDNSPolicy(api.CreateExternalDNSPolicyOpts{
		ID: opts.ID,
	})
	if err != nil {
		return nil, err
	}

	account, err := s.api.CreateExternalDNSServiceAccount(api.CreateExternalDNSServiceAccountOpts{
		CreateServiceAccountOpts: api.CreateServiceAccountOpts{
			ID:        opts.ID,
			PolicyArn: policy.PolicyARN,
		},
	})
	if err != nil {
		return nil, err
	}

	kube, err := s.api.CreateExternalDNSKubeDeployment(api.CreateExternalDNSKubeDeploymentOpts{
		ID:           opts.ID,
		HostedZoneID: opts.HostedZoneID,
		DomainFilter: opts.Domain,
	})
	if err != nil {
		return nil, err
	}

	externalDNS := &client.ExternalDNS{
		Policy:         policy,
		ServiceAccount: account,
		Kube:           kube,
	}

	report, err := s.store.SaveExternalDNS(externalDNS)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreateExternalDNS(externalDNS, report)
	if err != nil {
		return nil, err
	}

	return externalDNS, nil
}

// NewExternalDNSService returns an initialised service
func NewExternalDNSService(spinner spinner.Spinner, api client.ExternalDNSAPI, store client.ExternalDNSStore, report client.ExternalDNSReport) client.ExternalDNSService {
	return &externalDNSService{
		spinner: spinner,
		api:     api,
		store:   store,
		report:  report,
	}
}
