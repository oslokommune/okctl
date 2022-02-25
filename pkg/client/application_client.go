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

	CertificateARN string
}

// Validate ensures presented data is valid
func (o *ScaffoldApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Cluster, validation.Required),
		validation.Field(&o.Application, validation.Required),
	)
}

// DeleteApplicationManifestsOpts contains necessary information to delete application manifests
type DeleteApplicationManifestsOpts struct {
	Cluster     v1alpha1.Cluster
	Application v1alpha1.Application
}

// Validate ensures presented data is valid
func (o DeleteApplicationManifestsOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Cluster, validation.Required),
		validation.Field(&o.Application, validation.Required),
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

// HasArgoCDIntegrationOpts defines required information to verify an existing ArgoCD integration
type HasArgoCDIntegrationOpts struct {
	Cluster     v1alpha1.Cluster
	Application v1alpha1.Application
}

// Validate ensures presented data is valid
func (c HasArgoCDIntegrationOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Cluster, validation.Required),
		validation.Field(&c.Application, validation.Required),
	)
}

// ApplicationService applies the scaffolding API and produces the requested resources
type ApplicationService interface {
	// ScaffoldApplication implements functionality for converting an Application.yaml to deployment resources
	ScaffoldApplication(context.Context, *ScaffoldApplicationOpts) error
	// DeleteApplicationManifests implements functionality for deleting application manifests
	DeleteApplicationManifests(context.Context, DeleteApplicationManifestsOpts) error
	// CreateArgoCDApplicationManifest implements functionality for integrating an app with ArgoCD
	CreateArgoCDApplicationManifest(opts CreateArgoCDApplicationManifestOpts) error
	// DeleteArgoCDApplicationManifest implements functionality for removing an ArgoCD integration for an app
	DeleteArgoCDApplicationManifest(opts DeleteArgoCDApplicationManifestOpts) error
	// HasArgoCDIntegration implements functionality for verifying if there is an existing ArgoCD integration for the app
	HasArgoCDIntegration(context.Context, HasArgoCDIntegrationOpts) (bool, error)
}

// ApplicationState defines operations done on application state
type ApplicationState interface {
	// Put knows how to upsert application state
	Put(v1alpha1.Application) error
	// Get knows how to retrieve application state
	Get(name string) (v1alpha1.Application, error)
	// Delete knows how to remove application state
	Delete(name string) error
	// List knows how to retrieve all application state stored
	List() ([]v1alpha1.Application, error)
	// Initialize knows how to do required setup for state to work
	Initialize(clusterManifest v1alpha1.Cluster, absoluteIACRepositoryRootDirectory string) error
}
