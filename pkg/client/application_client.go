package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ScaffoldApplicationOpts contains information necessary to scaffold application resources
type ScaffoldApplicationOpts struct {
	Cluster     v1alpha1.Cluster
	Application v1alpha1.Application

	HostedZoneID string
}

// Validate ensures presented data is valid
func (o *ScaffoldApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Cluster, validation.Required),
		validation.Field(&o.Application, validation.Required),
		validation.Field(&o.HostedZoneID, validation.Required),
	)
}

// CreateArgoCDApplicationManifestOpts contains data required when creating a ArgoCD Application Manifest
type CreateArgoCDApplicationManifestOpts struct {
	Cluster     v1alpha1.Cluster
	Application v1alpha1.Application
}

// Validate ensures presented data is valid
func (c CreateArgoCDApplicationManifestOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Cluster, validation.Required),
		validation.Field(&c.Application, validation.Required),
	)
}

// DeleteArgoCDApplicationManifestOpts contains data required when deleting an ArgoCD Application Manifest
type DeleteArgoCDApplicationManifestOpts struct {
	Cluster     v1alpha1.Cluster
	Application v1alpha1.Application
}

// Validate ensures presented data is valid
func (c DeleteArgoCDApplicationManifestOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Cluster, validation.Required),
		validation.Field(&c.Application, validation.Required),
	)
}

// ApplicationService applies the scaffolding API and produces the requested resources
type ApplicationService interface {
	// ScaffoldApplication implements functionality for converting an Application.yaml to deployment resources
	ScaffoldApplication(context.Context, *ScaffoldApplicationOpts) error
	// CreateArgoCDApplicationManifest implements functionality for integrating an app with ArgoCD
	CreateArgoCDApplicationManifest(opts CreateArgoCDApplicationManifestOpts) error
	// DeleteArgoCDApplicationManifest implements functionality for removing an ArgoCD integration for an app
	DeleteArgoCDApplicationManifest(opts DeleteArgoCDApplicationManifestOpts) error
}
