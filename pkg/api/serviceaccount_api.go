package api

import (
	"context"

	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ServiceAccount holds state for a service account
type ServiceAccount struct {
	ID        ID
	Name      string
	PolicyArn string
	Config    *v1alpha5.ClusterConfig
}

// CreateServiceAccountBaseOpts contains opts shared state
type CreateServiceAccountBaseOpts struct {
	ID        ID
	PolicyArn string
}

// ValidateStruct validates the shared state
func (o CreateServiceAccountBaseOpts) ValidateStruct() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.PolicyArn, validation.Required),
	)
}

// CreateExternalSecretsServiceAccountOpts contains the configuration
// for creating a external secrets service account
type CreateExternalSecretsServiceAccountOpts struct {
	CreateServiceAccountBaseOpts
}

// Validate the options
func (o CreateExternalSecretsServiceAccountOpts) Validate() error {
	return o.ValidateStruct()
}

// CreateAWSLoadBalancerControllerServiceAccountOpts contains the configuration
// for creating an alb ingress controller service account
type CreateAWSLoadBalancerControllerServiceAccountOpts struct {
	CreateServiceAccountBaseOpts
}

// Validate the options
func (o CreateAWSLoadBalancerControllerServiceAccountOpts) Validate() error {
	return o.ValidateStruct()
}

// CreateExternalDNSServiceAccountOpts contains the configuration
// for creating an external dns service account
type CreateExternalDNSServiceAccountOpts struct {
	CreateServiceAccountBaseOpts
}

// Validate the input
func (o CreateExternalDNSServiceAccountOpts) Validate() error {
	return o.ValidateStruct()
}

// CreateAutoscalerServiceAccountOpts contains the required inputs
type CreateAutoscalerServiceAccountOpts struct {
	CreateServiceAccountBaseOpts
}

// Validate the inputs
func (o CreateAutoscalerServiceAccountOpts) Validate() error {
	return o.ValidateStruct()
}

// CreateBlockstorageServiceAccountOpts contains the required inputs
type CreateBlockstorageServiceAccountOpts struct {
	CreateServiceAccountBaseOpts
}

// Validate the inputs
func (o CreateBlockstorageServiceAccountOpts) Validate() error {
	return o.ValidateStruct()
}

// CreateServiceAccountOpts contains the inputs required
// for creating a new service account
type CreateServiceAccountOpts struct {
	ID        ID
	Name      string
	PolicyArn string
	Config    *v1alpha5.ClusterConfig
}

// Validate the provided inputs
func (o CreateServiceAccountOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Name, validation.Required),
		validation.Field(&o.PolicyArn, validation.Required),
		validation.Field(&o.Config, validation.Required),
	)
}

// DeleteServiceAccountOpts contains the inputs required
// for deleting a service account
type DeleteServiceAccountOpts struct {
	ID     ID
	Name   string
	Config *v1alpha5.ClusterConfig
}

// Validate the provided inputs
func (o DeleteServiceAccountOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Name, validation.Required),
		validation.Field(&o.Config, validation.Required),
	)
}

// ServiceAccountService provides the interface for all service account operations
type ServiceAccountService interface {
	CreateExternalSecretsServiceAccount(context.Context, CreateExternalSecretsServiceAccountOpts) (*ServiceAccount, error)
	DeleteExternalSecretsServiceAccount(context.Context, ID) error
	CreateAWSLoadBalancerControllerServiceAccount(context.Context, CreateAWSLoadBalancerControllerServiceAccountOpts) (*ServiceAccount, error)
	DeleteAWSLoadBalancerControllerServiceAccount(context.Context, ID) error
	CreateExternalDNSServiceAccount(context.Context, CreateExternalDNSServiceAccountOpts) (*ServiceAccount, error)
	DeleteExternalDNSServiceAccount(context.Context, ID) error
	CreateAutoscalerServiceAccount(context.Context, CreateAutoscalerServiceAccountOpts) (*ServiceAccount, error)
	DeleteAutoscalerServiceAccount(context.Context, ID) error
	CreateBlockstorageServiceAccount(context.Context, CreateBlockstorageServiceAccountOpts) (*ServiceAccount, error)
	DeleteBlockstorageServiceAccount(context.Context, ID) error
	CreateServiceAccount(ctx context.Context, opts CreateServiceAccountOpts) (*ServiceAccount, error)
	DeleteServiceAccount(ctx context.Context, opts DeleteServiceAccountOpts) error
}

// ServiceAccountRun provides the interface for running operations
type ServiceAccountRun interface {
	CreateServiceAccount(*v1alpha5.ClusterConfig) error
	DeleteServiceAccount(*v1alpha5.ClusterConfig) error
}
