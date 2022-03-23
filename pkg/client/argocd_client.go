package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/clients/kubectl"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

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
	ClusterManifest v1alpha1.Cluster
	Domain          string
	FQDN            string
	HostedZoneID    string
	UserPoolID      string
	AuthDomain      string
	Repository      *GithubRepository
}

// DeleteArgoCDOpts contains the required inputs
// for deleting an argocd installation
type DeleteArgoCDOpts struct {
	ID api.ID
}

// ArgoCDService is a business logic implementation
type ArgoCDService interface {
	// CreateArgoCD defines functionality for installing ArgoCD into a cluster
	CreateArgoCD(ctx context.Context, opts CreateArgoCDOpts) (*ArgoCD, error)
	// DeleteArgoCD defines functionality for uninstalling ArgoCD from a cluster
	DeleteArgoCD(ctx context.Context, opts DeleteArgoCDOpts) error
	// SetupApplicationsSync defines functionality for preparing a directory where ArgoCD application manifests will be
	// automatically synced
	SetupApplicationsSync(ctx context.Context, cluster v1alpha1.Cluster) error
	// SetupNamespacesSync defines functionality for preparing a directory where namespace manifests will be
	// automatically synced
	SetupNamespacesSync(ctx context.Context, kubectlClient kubectl.Client, cluster v1alpha1.Cluster) error
}

// ArgoCDState implements the state layer
type ArgoCDState interface {
	SaveArgoCD(cd *ArgoCD) error
	GetArgoCD() (*ArgoCD, error)
	HasArgoCD() (bool, error)
	RemoveArgoCD() error
}
