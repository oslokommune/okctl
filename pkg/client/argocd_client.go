package client

import (
	"context"
	"io"

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

// Applier defines necessary functionality for applying manifests to a Kubernetes cluster
type Applier interface {
	// Apply knows how to apply a manifest to a Kubernetes cluster
	Apply(io.Reader) error
}

// GitDeleteRemoteFileFn knows how to remotely delete a file tracked by Git
type GitDeleteRemoteFileFn func(repositoryURL string, path string, commitMessage string) error

// ArgoCDService is a business logic implementation
type ArgoCDService interface {
	// CreateArgoCD defines functionality for installing ArgoCD into a cluster
	CreateArgoCD(context.Context, CreateArgoCDOpts) (*ArgoCD, error)
	// DeleteArgoCD defines functionality for uninstalling ArgoCD from a cluster
	DeleteArgoCD(context.Context, DeleteArgoCDOpts) error
	// SetupApplicationsSync defines functionality for preparing a directory where ArgoCD application manifests will be
	// automatically synced
	SetupApplicationsSync(context.Context, Applier, v1alpha1.Cluster) error
	// SetupNamespacesSync defines functionality for preparing a directory where namespace manifests will be
	// automatically synced
	SetupNamespacesSync(context.Context, Applier, v1alpha1.Cluster) error
}

// ArgoCDState implements the state layer
type ArgoCDState interface {
	SaveArgoCD(cd *ArgoCD) error
	GetArgoCD() (*ArgoCD, error)
	HasArgoCD() (bool, error)
	RemoveArgoCD() error
}
