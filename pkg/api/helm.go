package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/helm"
	"helm.sh/helm/v3/pkg/release"
)

// Helm contains the data of a helm release
type Helm struct {
	Repository  string
	Environment string
	Release     *release.Release
	Chart       *helm.Chart
}

// CreateExternalSecretsHelmChartOpts contains the data
// required for creating a helm chart
type CreateExternalSecretsHelmChartOpts struct {
	Repository  string
	Environment string
}

// Validate the helm create inputs
func (o CreateExternalSecretsHelmChartOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.Repository, validation.Required),
	)
}

// CreateAlbIngressControllerHelmChartOpts contains the data
// required for creating a helm chart
type CreateAlbIngressControllerHelmChartOpts struct {
	ClusterName string
	Repository  string
	Environment string
	VpcID       string
	Region      string
}

// Validate the input options
func (o CreateAlbIngressControllerHelmChartOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterName, validation.Required),
		validation.Field(&o.Repository, validation.Required),
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.VpcID, validation.Required),
		validation.Field(&o.Region, validation.Required),
	)
}

// HelmService defines the service layer interface
type HelmService interface {
	CreateExternalSecretsHelmChart(ctx context.Context, opts CreateExternalSecretsHelmChartOpts) (*Helm, error)
	CreateAlbIngressControllerHelmChart(ctx context.Context, opts CreateAlbIngressControllerHelmChartOpts) (*Helm, error)
}

// HelmRun defines the runner layer
type HelmRun interface {
	CreateExternalSecretsHelmChart(opts CreateExternalSecretsHelmChartOpts) (*Helm, error)
	CreateAlbIngressControllerHelmChart(opts CreateAlbIngressControllerHelmChartOpts) (*Helm, error)
}

// HelmStore defines the storage layer
type HelmStore interface {
	SaveExternalSecretsHelmChart(*Helm) error
	SaveAlbIngressControllerHelmChar(*Helm) error
}
