package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ServiceAccount holds state for a service account
type ServiceAccount struct {
	ClusterName  string
	Environment  string
	Region       string
	AWSAccountID string
	PolicyArn    string
	Config       *ClusterConfig
}

// CreateExternalSecretsServiceAccountOpts contains the configuration
// for creating a external secrets service account
type CreateExternalSecretsServiceAccountOpts struct {
	ClusterName  string
	Environment  string
	Region       string
	AWSAccountID string
	PolicyArn    string
}

// Validate the options
func (o CreateExternalSecretsServiceAccountOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.PolicyArn, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
		validation.Field(&o.Region, validation.Required),
	)
}

// ServiceAccountService provides the interface for all service account operations
type ServiceAccountService interface {
	CreateExternalSecretsServiceAccount(context.Context, CreateExternalSecretsServiceAccountOpts) (*ServiceAccount, error)
}

// ServiceAccountRun provides the interface for running operations
type ServiceAccountRun interface {
	CreateExternalSecretsServiceAccount(*ClusterConfig) error
}

// ServiceAccountStore provides the storage operations
type ServiceAccountStore interface {
	SaveExternalSecretsServiceAccount(*ServiceAccount) error
	GetExternalSecretsServiceAccount() (*ServiceAccount, error)
}
