package client

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/helm"
	"helm.sh/helm/v3/pkg/release"
)

// Helm contains the data of a helm release
type Helm struct {
	ID      api.ID
	Release *release.Release
	Chart   *helm.Chart
}

// CreateHelmReleaseOpts contains the required inputs for
// installing a Helm release on the Kubernetes cluster
type CreateHelmReleaseOpts struct {
	ID             api.ID
	RepositoryName string
	RepositoryURL  string
	ReleaseName    string
	Version        string
	Chart          string
	Namespace      string
	Values         []byte
}

// Validate the provided inputs
func (o CreateHelmReleaseOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.RepositoryName, validation.Required),
		validation.Field(&o.RepositoryURL, validation.Required),
		validation.Field(&o.ReleaseName, validation.Required),
		validation.Field(&o.Version, validation.Required),
		validation.Field(&o.Chart, validation.Required),
		validation.Field(&o.Namespace, validation.Required),
		validation.Field(&o.Values, validation.Required),
	)
}

// DeleteHelmReleaseOpts contains the required inputs for
// removing a Helm release from the Kubernetes cluster
// Experimenting a little bit with a more generic/general
// interface, starting with Delete operation
type DeleteHelmReleaseOpts struct {
	ID          api.ID
	ReleaseName string
	Namespace   string
}

// Validate the provided inputs
func (o DeleteHelmReleaseOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.ReleaseName, validation.Required),
		validation.Field(&o.Namespace, validation.Required),
	)
}

// GetHelmReleaseOpts contains the required inputs for
// getting a Helm release from the Kubernetes cluster
type GetHelmReleaseOpts struct {
	ClusterID   api.ID
	ReleaseName string
	Namespace   string
}

// Validate the provided inputs
func (o GetHelmReleaseOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterID, validation.Required),
		validation.Field(&o.ReleaseName, validation.Required),
		validation.Field(&o.Namespace, validation.Required),
	)
}

// HelmService provides functionality for interacting with helm chart lifecycle
type HelmService interface {
	CreateHelmRelease(ctx context.Context, opts CreateHelmReleaseOpts) (*Helm, error)
	DeleteHelmRelease(ctx context.Context, opts DeleteHelmReleaseOpts) error
	GetHelmRelease(ctx context.Context, opts GetHelmReleaseOpts) (*Helm, error)
}

// HelmAPI provides the invocation layer
type HelmAPI interface {
	CreateHelmRelease(opts api.CreateHelmReleaseOpts) (*api.Helm, error)
	DeleteHelmRelease(opts api.DeleteHelmReleaseOpts) error
	GetHelmRelease(opts api.GetHelmReleaseOpts) (*api.Helm, error)
}

// HelmState provides the persistence layer
type HelmState interface {
	SaveHelmRelease(helm *Helm) error
	RemoveHelmRelease(releaseName string) error
	GetHelmRelease(releaseName string) (*Helm, error)
}
