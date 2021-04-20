package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
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
