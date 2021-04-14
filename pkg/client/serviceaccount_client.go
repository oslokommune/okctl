package client

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/oslokommune/okctl/pkg/api"
)

// ServiceAccount holds state for a service account
type ServiceAccount struct {
	ID        api.ID
	Name      string
	PolicyArn string
	Config    *v1alpha5.ClusterConfig
}

// CreateServiceAccountOpts contains the inputs required
// for creating a new service account
type CreateServiceAccountOpts struct {
	ID        api.ID
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
	ID     api.ID
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

// ServiceAccountService implements the business logic
type ServiceAccountService interface {
	CreateServiceAccount(ctx context.Context, opts CreateServiceAccountOpts) (*ServiceAccount, error)
	DeleteServiceAccount(ctx context.Context, opts DeleteServiceAccountOpts) error
}

// ServiceAccountAPI invokes the remote API
type ServiceAccountAPI interface {
	CreateServiceAccount(opts api.CreateServiceAccountOpts) (*api.ServiceAccount, error)
	DeleteServiceAccount(opts api.DeleteServiceAccountOpts) error
}

// ServiceAccountState defines the state layer
type ServiceAccountState interface {
	SaveServiceAccount(account *ServiceAccount) error
	RemoveServiceAccount(name string) error
	GetServiceAccount(name string) (*ServiceAccount, error)
	UpdateServiceAccount(account *ServiceAccount) error
}
