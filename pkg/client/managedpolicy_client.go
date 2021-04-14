package client

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/oslokommune/okctl/pkg/api"
)

// ManagedPolicy holds state for a managed policy
type ManagedPolicy struct {
	ID                     api.ID
	StackName              string
	PolicyARN              string
	CloudFormationTemplate []byte
}

// CreatePolicyOpts contains all the required inputs for
// creating a managed policy
type CreatePolicyOpts struct {
	ID                     api.ID
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
	ID        api.ID
	StackName string
}

// Validate the required inputs
func (o DeletePolicyOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.StackName, validation.Required),
	)
}

// ManagedPolicyService implements the business logic
type ManagedPolicyService interface {
	CreatePolicy(ctx context.Context, opts CreatePolicyOpts) (*ManagedPolicy, error)
	DeletePolicy(ctx context.Context, opts DeletePolicyOpts) error
}

// ManagedPolicyAPI invokes the remote API
type ManagedPolicyAPI interface {
	CreatePolicy(opts api.CreatePolicyOpts) (*api.ManagedPolicy, error)
	DeletePolicy(opts api.DeletePolicyOpts) error
}

// ManagedPolicyState provides a persistence layer
type ManagedPolicyState interface {
	SavePolicy(policy *ManagedPolicy) error
	GetPolicy(stackName string) (*ManagedPolicy, error)
	RemovePolicy(stackName string) error
}
