package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client"
)

type argoCDService struct {
	api    client.ArgoCDAPI
	store  client.ArgoCDStore
	report client.ArgoCDReport
}

func (s *argoCDService) CreateArgoCD(_ context.Context, opts client.CreateArgoCDOpts) (*client.ArgoCD, error) {
	argo, err := s.api.CreateArgoCD(opts)
	if err != nil {
		return nil, err
	}

	reports, err := s.store.SaveArgoCD(argo)
	if err != nil {
		return nil, err
	}

	err = s.report.CreateArgoCD(argo, reports)
	if err != nil {
		return nil, err
	}

	return argo, nil
}

// NewArgoCDService returns an initialised service
func NewArgoCDService(api client.ArgoCDAPI, store client.ArgoCDStore, report client.ArgoCDReport) client.ArgoCDService {
	return &argoCDService{
		api:    api,
		store:  store,
		report: report,
	}
}
