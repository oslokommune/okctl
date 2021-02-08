package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/helm"
	"helm.sh/helm/v3/pkg/release"
)

// Helm contains the data of a helm release
type Helm struct {
	ID      ID
	Release *release.Release
	Chart   *helm.Chart
}

// CreateExternalSecretsHelmChartOpts contains the data
// required for creating a helm chart
type CreateExternalSecretsHelmChartOpts struct {
	ID ID
}

// Validate the helm create inputs
func (o CreateExternalSecretsHelmChartOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
	)
}

// CreateAlbIngressControllerHelmChartOpts contains the data
// required for creating a helm chart
type CreateAlbIngressControllerHelmChartOpts struct {
	ID    ID
	VpcID string
}

// Validate the input options
func (o CreateAlbIngressControllerHelmChartOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.VpcID, validation.Required),
	)
}

// CreateAWSLoadBalancerControllerHelmChartOpts contains the data
// required for creating a helm chart
type CreateAWSLoadBalancerControllerHelmChartOpts struct {
	ID    ID
	VpcID string
}

// Validate the input options
func (o CreateAWSLoadBalancerControllerHelmChartOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.VpcID, validation.Required),
	)
}

// CreateArgoCDOpts contains the required inputs
type CreateArgoCDOpts struct {
	ID ID

	ArgoDomain         string
	ArgoCertificateARN string

	GithubOrganisation string
	GithubRepoURL      string
	GithubRepoName     string

	ClientID   string
	AuthDomain string
	UserPoolID string

	PrivateKeyName string
	PrivateKeyKey  string
}

// Validate the input options
func (o CreateArgoCDOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.ArgoDomain, validation.Required),
		validation.Field(&o.ArgoCertificateARN, validation.Required),
		validation.Field(&o.GithubOrganisation, validation.Required),
		validation.Field(&o.GithubRepoURL, validation.Required),
		validation.Field(&o.GithubRepoName, validation.Required),
		validation.Field(&o.ClientID, validation.Required),
		validation.Field(&o.PrivateKeyName, validation.Required),
		validation.Field(&o.PrivateKeyKey, validation.Required),
	)
}

// HelmService defines the service layer interface
type HelmService interface {
	CreateExternalSecretsHelmChart(ctx context.Context, opts CreateExternalSecretsHelmChartOpts) (*Helm, error)
	CreateAlbIngressControllerHelmChart(ctx context.Context, opts CreateAlbIngressControllerHelmChartOpts) (*Helm, error)
	CreateAWSLoadBalancerControllerHelmChart(ctx context.Context, opts CreateAWSLoadBalancerControllerHelmChartOpts) (*Helm, error)
	CreateArgoCD(ctx context.Context, opts CreateArgoCDOpts) (*Helm, error)
}

// HelmRun defines the runner layer
type HelmRun interface {
	CreateExternalSecretsHelmChart(opts CreateExternalSecretsHelmChartOpts) (*Helm, error)
	CreateAlbIngressControllerHelmChart(opts CreateAlbIngressControllerHelmChartOpts) (*Helm, error)
	CreateAWSLoadBalancerControllerHelmChart(opts CreateAWSLoadBalancerControllerHelmChartOpts) (*Helm, error)
	CreateArgoCD(opts CreateArgoCDOpts) (*Helm, error)
}

// HelmStore defines the storage layer
type HelmStore interface {
	SaveExternalSecretsHelmChart(*Helm) error
	SaveAlbIngressControllerHelmChart(*Helm) error
	SaveAWSLoadBalancerControllerHelmChart(*Helm) error
	SaveArgoCD(*Helm) error
}
