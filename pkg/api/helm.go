package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/helm"
	"helm.sh/helm/v3/pkg/release"
)

type Helm struct {
	Repository  string
	Environment string
	Release     *release.Release
	Chart       *helm.Chart
}

type CreateExternalSecretsHelmChartOpts struct {
	Repository  string
	Environment string
}

func (o CreateExternalSecretsHelmChartOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.Repository, validation.Required),
	)
}

type HelmService interface {
	CreateExternalSecretsHelmChart(ctx context.Context, opts CreateExternalSecretsHelmChartOpts) (*Helm, error)
}

type HelmRun interface {
	CreateExternalSecretsHelmChart(opts CreateExternalSecretsHelmChartOpts) (*Helm, error)
}

type HelmStore interface {
	SaveExternalSecretsHelmChart(*Helm) error
}
