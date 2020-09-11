package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client/store"
)

// ArgoCD contains state about an argo cd deployment
type ArgoCD struct {
	ID             api.ID
	ArgoDomain     string
	ArgoURL        string
	Certificate    *api.Certificate
	GithubOauthApp *GithubOauthApp
	ExternalSecret *api.ExternalSecretsKube
	SecretKey      *api.SecretParameter
	Chart          *api.Helm
}

// CreateArgoCDOpts contains the required inputs
// for setting up argo cd
type CreateArgoCDOpts struct {
	ID                 api.ID
	Domain             string
	FQDN               string
	HostedZoneID       string
	GithubOrganisation string
	Repository         *GithubRepository
}

// ArgoCDService is a business logic implementation
type ArgoCDService interface {
	CreateArgoCD(ctx context.Context, opts CreateArgoCDOpts) (*ArgoCD, error)
}

// ArgoCDAPI invokes the APIs for creating resources
type ArgoCDAPI interface {
	CreateArgoCD(opts api.CreateArgoCDOpts) (*ArgoCD, error)
}

// ArgoCDStore implements the storage layer
type ArgoCDStore interface {
	SaveArgoCD(cd *ArgoCD) (*store.Report, error)
}

// ArgoCDReport implements the report layer
type ArgoCDReport interface {
	CreateArgoCD(cd *ArgoCD, reports []*store.Report) error
}

// ArgoCDState implements the state layer
type ArgoCDState interface {
	SaveArgoCD(cd *ArgoCD) (*store.Report, error)
}
