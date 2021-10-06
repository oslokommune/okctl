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
	CreatePolicy(ctx context.Context, opts CreatePolicyOpts) (*ManagedPolicy, error)
	DeletePolicy(ctx context.Context, opts DeletePolicyOpts) error
}

// ManagedPolicyCloudProvider defines the cloud provider layer for managed policies
type ManagedPolicyCloudProvider interface {
	CreatePolicy(ctx context.Context, opts CreatePolicyOpts) (*ManagedPolicy, error)
	DeletePolicy(opts DeletePolicyOpts) error
}
