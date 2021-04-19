package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/go-ozzo/ozzo-validation/v4/is"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

const (
	lowestPossiblePort     = 1
	highestPossiblePort    = 65535
	lowestPossibleReplicas = 0
)

// ScaffoldApplicationOpts contains information necessary to scaffold application resources
type ScaffoldApplicationOpts struct {
	OutputDir string

	ID               *api.ID
	HostedZoneID     string
	HostedZoneDomain string
	IACRepoURL       string
	Application      v1alpha1.Application
}

// Validate ensures presented data is valid
func (o *ScaffoldApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.HostedZoneID, validation.Required),
		validation.Field(&o.IACRepoURL, validation.Required),
		validation.Field(&o.Application, validation.Required),
	)
}

// OkctlApplication represents the necessary information okctl needs to deploy an application
type OkctlApplication struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`

	Image           string `json:"image"`
	Version         string `json:"version"`
	ImagePullSecret string `json:"imagePullSecret"`

	SubDomain string `json:"subDomain"`
	Port      int32  `json:"port"`

	Replicas int32 `json:"replicas"`

	Environment map[string]string   `json:"environment"`
	Volumes     []map[string]string `json:"volumes"`

	ContainerRepositories []string `json:"containerRepositories"`
}

// HasIngress returns true if the application requires an ingress
func (o OkctlApplication) HasIngress() bool {
	return o.SubDomain != ""
}

// HasService returns true if the application requires a service
func (o OkctlApplication) HasService() bool {
	return o.Port > 0
}

// Validate knows if the application is valid or not
func (o OkctlApplication) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Name, validation.Required, is.Subdomain),
		validation.Field(&o.Namespace, validation.Required, is.Subdomain),
		validation.Field(&o.Image, validation.Required),
		validation.Field(&o.Version, validation.Required),
		validation.Field(&o.ImagePullSecret, is.Subdomain),
		validation.Field(&o.SubDomain, is.Subdomain),
		validation.Field(&o.Port, validation.Min(lowestPossiblePort), validation.Max(highestPossiblePort)),
		validation.Field(&o.Replicas, validation.Min(lowestPossibleReplicas)),
		validation.Field(&o.ContainerRepositories, validation.Each(is.Alphanumeric)),
	)
}

// ScaffoldedApplication contains information required by ApplicationStore and ApplicationReport
type ScaffoldedApplication struct {
	ApplicationName string
	ClusterName     string

	BaseKustomization []byte
	ArgoCDResource    []byte
	Volume            []byte
	Ingress           []byte

	OverlayKustomization []byte
	Service              []byte
	Deployment           []byte
	IngressPatch         []byte
	ServicePatch         []byte
	DeploymentPatch      []byte
}

// ApplicationService applies the scaffolding API and produces the requested resources
type ApplicationService interface {
	// ScaffoldApplication implements functionality for converting an Application.yaml to deployment resources
	ScaffoldApplication(context.Context, *ScaffoldApplicationOpts) error
}

// ApplicationStore handles writing deployment resources to persistent storage
type ApplicationStore interface {
	// SaveApplication should implement functionality for storing the scaffolded application in som form of persistent storage
	SaveApplication(*ScaffoldedApplication) (*store.Report, error)
	// RemoveApplication should implement functionality for removing the scaffolded application from the persistent storage
	RemoveApplication(string) (*store.Report, error)
}
