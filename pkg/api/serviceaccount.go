package api

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ServiceAccount holds state for a service account
type ServiceAccount struct {
	ID        ID
	PolicyArn string
	Config    *v1alpha1.ClusterConfig
}

// CreateServiceAccountOpts contains opts shared state
type CreateServiceAccountOpts struct {
	ID        ID
	PolicyArn string
}

// ValidateStruct validates the shared state
func (o CreateServiceAccountOpts) ValidateStruct() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.PolicyArn, validation.Required),
	)
}

// CreateExternalSecretsServiceAccountOpts contains the configuration
// for creating a external secrets service account
type CreateExternalSecretsServiceAccountOpts struct {
	CreateServiceAccountOpts
}

// Validate the options
func (o CreateExternalSecretsServiceAccountOpts) Validate() error {
	return o.ValidateStruct()
}

// CreateAlbIngressControllerServiceAccountOpts contains the configuration
// for creating an alb ingress controller service account
type CreateAlbIngressControllerServiceAccountOpts struct {
	CreateServiceAccountOpts
}

// Validate the options
func (o CreateAlbIngressControllerServiceAccountOpts) Validate() error {
	return o.ValidateStruct()
}

// CreateExternalDNSServiceAccountOpts contains the configuration
// for creating an external dns service account
type CreateExternalDNSServiceAccountOpts struct {
	CreateServiceAccountOpts
}

// Validate the input
func (o CreateExternalDNSServiceAccountOpts) Validate() error {
	return o.ValidateStruct()
}

// ServiceAccountService provides the interface for all service account operations
type ServiceAccountService interface {
	CreateExternalSecretsServiceAccount(context.Context, CreateExternalSecretsServiceAccountOpts) (*ServiceAccount, error)
	CreateAlbIngressControllerServiceAccount(context.Context, CreateAlbIngressControllerServiceAccountOpts) (*ServiceAccount, error)
	CreateExternalDNSServiceAccount(context.Context, CreateExternalDNSServiceAccountOpts) (*ServiceAccount, error)
}

// ServiceAccountRun provides the interface for running operations
type ServiceAccountRun interface {
	CreateServiceAccount(*v1alpha1.ClusterConfig) error
}

// ServiceAccountStore provides the storage operations
type ServiceAccountStore interface {
	SaveExternalSecretsServiceAccount(*ServiceAccount) error
	GetExternalSecretsServiceAccount() (*ServiceAccount, error)
	SaveAlbIngressControllerServiceAccount(*ServiceAccount) error
	GetAlbIngressControllerServiceAccount() (*ServiceAccount, error)
	SaveExternalDNSServiceAccount(*ServiceAccount) error
	GetExternalDNSServiceAccount() (*ServiceAccount, error)
}
