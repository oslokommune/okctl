package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/helm"
	"helm.sh/helm/v3/pkg/release"
)

// Helm contains the data of a helm release
type Helm struct {
	ClusterName string
	Repository  string
	Environment string
	Release     *release.Release
	Chart       *helm.Chart
}

// CreateExternalSecretsHelmChartOpts contains the data
// required for creating a helm chart
type CreateExternalSecretsHelmChartOpts struct {
	ClusterName string
	Repository  string
	Environment string
}

// Validate the helm create inputs
func (o CreateExternalSecretsHelmChartOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterName, validation.Required),
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

// CreateArgoCDOpts contains the required inputs
type CreateArgoCDOpts struct {
	ClusterName string
	Repository  string
	Environment string

	ArgoDomain         string
	ArgoCertificateARN string

	GithubOrganisation  string
	GithubTeam          string
	GithubRepoURL       string
	GithubRepoName      string
	GithubOauthClientID string

	ExternalSecretName string
	ExternalSecretKey  string
}

// Validate the input options
func (o CreateArgoCDOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterName, validation.Required),
		validation.Field(&o.Repository, validation.Required),
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.ArgoDomain, validation.Required),
		validation.Field(&o.ArgoCertificateARN, validation.Required),
		validation.Field(&o.GithubOrganisation, validation.Required),
		validation.Field(&o.GithubTeam, validation.Required),
		validation.Field(&o.GithubRepoURL, validation.Required),
		validation.Field(&o.GithubRepoName, validation.Required),
		validation.Field(&o.GithubOauthClientID, validation.Required),
	)
}

// HelmService defines the service layer interface
type HelmService interface {
	CreateExternalSecretsHelmChart(ctx context.Context, opts CreateExternalSecretsHelmChartOpts) (*Helm, error)
	CreateAlbIngressControllerHelmChart(ctx context.Context, opts CreateAlbIngressControllerHelmChartOpts) (*Helm, error)
	CreateArgoCD(ctx context.Context, opts CreateArgoCDOpts) (*Helm, error)
}

// HelmRun defines the runner layer
type HelmRun interface {
	CreateExternalSecretsHelmChart(opts CreateExternalSecretsHelmChartOpts) (*Helm, error)
	CreateAlbIngressControllerHelmChart(opts CreateAlbIngressControllerHelmChartOpts) (*Helm, error)
	CreateArgoCD(opts CreateArgoCDOpts) (*Helm, error)
}

// HelmStore defines the storage layer
type HelmStore interface {
	SaveExternalSecretsHelmChart(*Helm) error
	SaveAlbIngressControllerHelmChar(*Helm) error
	SaveArgoCD(*Helm) error
}
