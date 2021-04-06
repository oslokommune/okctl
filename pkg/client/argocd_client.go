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
	AuthDomain     string
	Certificate    *Certificate
	IdentityClient *api.IdentityPoolClient
	PrivateKey     *KubernetesManifest
	Secret         *KubernetesManifest
	ClientSecret   *api.SecretParameter
	SecretKey      *api.SecretParameter
	Chart          *api.Helm
}

// ArgoCDStateInfo represents a subset of the available
// argocd state information
type ArgoCDStateInfo struct {
	ID         api.ID
	ArgoDomain string
}

// CreateArgoCDOpts contains the required inputs
// for setting up argo cd
type CreateArgoCDOpts struct {
	ID                 api.ID
	Domain             string
	FQDN               string
	HostedZoneID       string
	GithubOrganisation string
	UserPoolID         string
	AuthDomain         string
	Repository         *GithubRepository
}

// DeleteArgoCDOpts contains the required inputs
// for deleting an argocd installation
type DeleteArgoCDOpts struct {
	ID api.ID
}

// ArgoCDService is a business logic implementation
type ArgoCDService interface {
	CreateArgoCD(ctx context.Context, opts CreateArgoCDOpts) (*ArgoCD, error)
	DeleteArgoCD(ctx context.Context, opts DeleteArgoCDOpts) error
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
	GetArgoCD(id api.ID) ArgoCDStateInfo
}
