package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
)

// ArgoCD contains state about an argo cd deployment
type ArgoCD struct {
	ID             api.ID
	ArgoDomain     string
	ArgoURL        string
	AuthDomain     string
	Certificate    *Certificate
	IdentityClient *IdentityPoolClient
	PrivateKey     *KubernetesManifest
	Secret         *KubernetesManifest
	ClientSecret   *SecretParameter
	SecretKey      *SecretParameter
	Chart          *Helm
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

// ArgoCDState implements the state layer
type ArgoCDState interface {
	SaveArgoCD(cd *ArgoCD) error
	GetArgoCD() (*ArgoCD, error)
	RemoveArgoCD() error
}
