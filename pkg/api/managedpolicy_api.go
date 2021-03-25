package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ManagedPolicy contains all state for a policy
type ManagedPolicy struct {
	ID                     ID
	StackName              string
	PolicyARN              string
	CloudFormationTemplate []byte
}

// CreateExternalSecretsPolicyOpts contains the options
// that are required for creating an external secrets policy
type CreateExternalSecretsPolicyOpts struct {
	ID ID
}

// Validate determines if the options are valid
func (o CreateExternalSecretsPolicyOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
	)
}

// CreateAWSLoadBalancerControllerPolicyOpts contains the input
type CreateAWSLoadBalancerControllerPolicyOpts struct {
	ID ID
}

// Validate the input
func (o CreateAWSLoadBalancerControllerPolicyOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
	)
}

// CreateExternalDNSPolicyOpts contains the input
type CreateExternalDNSPolicyOpts struct {
	ID ID
}

// Validate the input
func (o CreateExternalDNSPolicyOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
	)
}

// CreateAutoscalerPolicy contains all required inputs
type CreateAutoscalerPolicy struct {
	ID ID
}

// Validate the inputs
func (o CreateAutoscalerPolicy) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
	)
}

// CreateBlockstoragePolicy contains all required inputs
type CreateBlockstoragePolicy struct {
	ID ID
}

// Validate the inputs
func (o CreateBlockstoragePolicy) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
	)
}

// CreatePolicyOpts contains all the required inputs for
// creating a managed policy
type CreatePolicyOpts struct {
	ID                     ID
	StackName              string
	PolicyOutputName       string
	CloudFormationTemplate []byte
}

// Validate the required inputs
func (o CreatePolicyOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.StackName, validation.Required),
		validation.Field(&o.PolicyOutputName, validation.Required),
		validation.Field(&o.CloudFormationTemplate, validation.Required),
	)
}

// DeletePolicyOpts contains all the required inputs
// for deleting a managed policy
type DeletePolicyOpts struct {
	ID        ID
	StackName string
}

// Validate the required inputs
func (o DeletePolicyOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.StackName, validation.Required),
	)
}

// ManagedPolicyService defines the service layer for managed policies
type ManagedPolicyService interface {
	CreateExternalSecretsPolicy(ctx context.Context, opts CreateExternalSecretsPolicyOpts) (*ManagedPolicy, error)
	DeleteExternalSecretsPolicy(ctx context.Context, id ID) error
	CreateAWSLoadBalancerControllerPolicy(ctx context.Context, opts CreateAWSLoadBalancerControllerPolicyOpts) (*ManagedPolicy, error)
	DeleteAWSLoadBalancerControllerPolicy(ctx context.Context, id ID) error
	CreateExternalDNSPolicy(ctx context.Context, opts CreateExternalDNSPolicyOpts) (*ManagedPolicy, error)
	DeleteExternalDNSPolicy(ctx context.Context, id ID) error
	CreateAutoscalerPolicy(ctx context.Context, opts CreateAutoscalerPolicy) (*ManagedPolicy, error)
	DeleteAutoscalerPolicy(ctx context.Context, id ID) error
	CreateBlockstoragePolicy(ctx context.Context, opts CreateBlockstoragePolicy) (*ManagedPolicy, error)
	DeleteBlockstoragePolicy(ctx context.Context, id ID) error
	CreatePolicy(ctx context.Context, opts CreatePolicyOpts) (*ManagedPolicy, error)
	DeletePolicy(ctx context.Context, opts DeletePolicyOpts) error
}

// ManagedPolicyCloudProvider defines the cloud provider layer for managed policies
type ManagedPolicyCloudProvider interface {
	CreateExternalSecretsPolicy(opts CreateExternalSecretsPolicyOpts) (*ManagedPolicy, error)
	DeleteExternalSecretsPolicy(id ID) error
	CreateAWSLoadBalancerControllerPolicy(opts CreateAWSLoadBalancerControllerPolicyOpts) (*ManagedPolicy, error)
	DeleteAWSLoadBalancerControllerPolicy(id ID) error
	CreateExternalDNSPolicy(opts CreateExternalDNSPolicyOpts) (*ManagedPolicy, error)
	DeleteExternalDNSPolicy(id ID) error
	CreateAutoscalerPolicy(opts CreateAutoscalerPolicy) (*ManagedPolicy, error)
	DeleteAutoscalerPolicy(id ID) error
	CreateBlockstoragePolicy(opts CreateBlockstoragePolicy) (*ManagedPolicy, error)
	DeleteBlockstoragePolicy(id ID) error
	CreatePolicy(opts CreatePolicyOpts) (*ManagedPolicy, error)
	DeletePolicy(opts DeletePolicyOpts) error
}
