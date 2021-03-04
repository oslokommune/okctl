package client

import (
	"context"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"io"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// ScaffoldApplicationOpts contains information necessary to scaffold application resources
type ScaffoldApplicationOpts struct {
	In  io.Reader
	Out io.Writer

	ApplicationFilePath string
	RepoDir             string

	ID               *api.ID
	HostedZoneID     string
	HostedZoneDomain string
	IACRepoURL       string
}

// Validate ensures presented data is valid
func (o *ScaffoldApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ApplicationFilePath, validation.Required),
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.HostedZoneID, validation.Required),
		validation.Field(&o.IACRepoURL, validation.Required),
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
		validation.Field(&o.Name, validation.Required, is.DNSName),
		validation.Field(&o.Namespace, validation.Required, is.DNSName),
		validation.Field(&o.Image, validation.Required, is.DNSName),
		validation.Field(&o.Version, validation.Required),
		validation.Field(&o.ImagePullSecret, validation.Required, is.DNSName),
		validation.Field(&o.SubDomain, is.Subdomain),
		validation.Field(&o.Port, is.Port),
		validation.Field(&o.Replicas, validation.Min(0)),
	)
}

// ScaffoldedApplication contains information required by ApplicationStore and ApplicationReport
type ScaffoldedApplication struct {
	ApplicationName string
	Environment     string

	BaseKustomization []byte
	ArgoCDResource    []byte
	Volume            []byte
	Ingress           []byte

	Service         []byte
	Deployment      []byte
	IngressPatch    []byte
	ServicePatch    []byte
	DeploymentPatch []byte
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

// ApplicationReport handles writing output and progress
type ApplicationReport interface {
	// ReportCreateApplication should implement a way of reporting when a ScaffoldedApplication is saved in the ApplicationStore
	ReportCreateApplication(*ScaffoldedApplication, []*store.Report) error
	// ReportDeleteApplication should implement a way of reporting when a ScaffoldedApplication is removed in the ApplicationStore
	ReportDeleteApplication([]*store.Report) error
}
