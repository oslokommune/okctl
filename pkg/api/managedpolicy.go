package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ManagedPolicy contains all state for a policy
type ManagedPolicy struct {
	StackName              string
	Repository             string
	Environment            string
	CloudFormationTemplate []byte
	PolicyARN              string
}

// CreateExternalSecretsPolicyOpts contains the options
// that are required for creating an external secrets policy
type CreateExternalSecretsPolicyOpts struct {
	Repository  string
	Environment string
}

// Validate determines if the options are valid
func (o CreateExternalSecretsPolicyOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.Repository, validation.Required),
	)
}

// CreateAlbIngressControllerPolicyOpts contains the input
type CreateAlbIngressControllerPolicyOpts struct {
	Repository  string
	Environment string
}

// Validate the input
func (o CreateAlbIngressControllerPolicyOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Repository, validation.Required),
		validation.Field(&o.Environment, validation.Required),
	)
}

// CreateExternalDNSPolicyOpts contains the input
type CreateExternalDNSPolicyOpts struct {
	Repository  string
	Environment string
}

// Validate the input
func (o CreateExternalDNSPolicyOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Repository, validation.Required),
		validation.Field(&o.Environment, validation.Required),
	)
}

// ManagedPolicyService defines the service layer for managed policies
type ManagedPolicyService interface {
	CreateExternalSecretsPolicy(ctx context.Context, opts CreateExternalSecretsPolicyOpts) (*ManagedPolicy, error)
	CreateAlbIngressControllerPolicy(ctx context.Context, opts CreateAlbIngressControllerPolicyOpts) (*ManagedPolicy, error)
	CreateExternalDNSPolicy(ctx context.Context, opts CreateExternalDNSPolicyOpts) (*ManagedPolicy, error)
}

// ManagedPolicyCloudProvider defines the cloud provider layer for managed policies
type ManagedPolicyCloudProvider interface {
	CreateExternalSecretsPolicy(opts CreateExternalSecretsPolicyOpts) (*ManagedPolicy, error)
	CreateAlbIngressControllerPolicy(opts CreateAlbIngressControllerPolicyOpts) (*ManagedPolicy, error)
	CreateExternalDNSPolicy(opts CreateExternalDNSPolicyOpts) (*ManagedPolicy, error)
}

// ManagedPolicyStore defines the storage layer for managed policies
type ManagedPolicyStore interface {
	SaveExternalSecretsPolicy(policy *ManagedPolicy) error
	GetExternalSecretsPolicy() (*ManagedPolicy, error)
	SaveAlbIngressControllerPolicy(policy *ManagedPolicy) error
	GetAlbIngressControllerPolicy() (*ManagedPolicy, error)
	SaveExternalDNSPolicy(policy *ManagedPolicy) error
	GetExternalDNSPolicy() (*ManagedPolicy, error)
}
