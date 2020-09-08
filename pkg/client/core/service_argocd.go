package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client"
)

type argoCDService struct {
	api   client.ArgoCDAPI
	store client.ArgoCDStore
}

func (s *argoCDService) CreateArgoCD(_ context.Context, opts client.CreateArgoCDOpts) (*client.ArgoCD, error) {
	argo, err := s.api.CreateArgoCD(opts)
	if err != nil {
		return nil, err
	}

	_, err = s.store.SaveArgoCD(argo)
	if err != nil {
		return nil, err
	}

	return argo, nil
}

// NewArgoCDService returns an initialised service
func NewArgoCDService(api client.ArgoCDAPI, store client.ArgoCDStore) client.ArgoCDService {
	return &argoCDService{
		api:   api,
		store: store,
	}
}
