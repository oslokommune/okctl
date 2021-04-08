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

// CreateLokiHelmChartOpts contains the required inputs
type CreateLokiHelmChartOpts struct {
	ID ID
}

// Validate the inputs
func (o CreateLokiHelmChartOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
	)
}

// CreatePromtailHelmChartOpts contains the required inputs
type CreatePromtailHelmChartOpts struct {
	ID ID
}

// Validate the inputs
func (o CreatePromtailHelmChartOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
	)
}

// CreateKubePrometheusStackOpts defines the required inputs
type CreateKubePrometheusStackOpts struct {
	ID                                  ID
	GrafanaCloudWatchServiceAccountName string
	CertificateARN                      string
	Hostname                            string
	AuthHostname                        string
	ClientID                            string
	SecretsConfigName                   string
	SecretsCookieSecretKey              string
	SecretsClientSecretKey              string
	SecretsAdminUserKey                 string
	SecretsAdminPassKey                 string
}

// Validate the inputs
func (o CreateKubePrometheusStackOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.GrafanaCloudWatchServiceAccountName, validation.Required),
		validation.Field(&o.CertificateARN, validation.Required),
		validation.Field(&o.Hostname, validation.Required),
		validation.Field(&o.AuthHostname, validation.Required),
		validation.Field(&o.ClientID, validation.Required),
		validation.Field(&o.SecretsConfigName, validation.Required),
		validation.Field(&o.SecretsCookieSecretKey, validation.Required),
		validation.Field(&o.SecretsAdminUserKey, validation.Required),
		validation.Field(&o.SecretsAdminPassKey, validation.Required),
	)
}

// CreateHelmReleaseOpts contains the required inputs for
// installing a Helm release on the Kubernetes cluster
type CreateHelmReleaseOpts struct {
	ID             ID
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
	ID          ID
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

// HelmService defines the service layer interface
type HelmService interface {
	CreateKubePrometheusStack(ctx context.Context, opts CreateKubePrometheusStackOpts) (*Helm, error)
	CreateLokiHelmChart(ctx context.Context, opts CreateLokiHelmChartOpts) (*Helm, error)
	CreatePromtailHelmChart(ctx context.Context, opts CreatePromtailHelmChartOpts) (*Helm, error)
	CreateHelmRelease(ctx context.Context, opts CreateHelmReleaseOpts) (*Helm, error)
	DeleteHelmRelease(ctx context.Context, opts DeleteHelmReleaseOpts) error
}

// HelmRun defines the runner layer
type HelmRun interface {
	CreateKubePromStack(opts CreateKubePrometheusStackOpts) (*Helm, error)
	CreateLokiHelmChart(opts CreateLokiHelmChartOpts) (*Helm, error)
	CreatePromtailHelmChart(opts CreatePromtailHelmChartOpts) (*Helm, error)
	CreateHelmRelease(opts CreateHelmReleaseOpts) (*Helm, error)
	DeleteHelmRelease(opts DeleteHelmReleaseOpts) error
}
