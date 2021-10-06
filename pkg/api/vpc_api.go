package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Vpc represents the state of an aws vpc
type Vpc struct {
	ID                     ID
	StackName              string
	CloudFormationTemplate []byte

	VpcID                    string
	Cidr                     string
	PublicSubnets            []VpcSubnet
	PrivateSubnets           []VpcSubnet
	DatabaseSubnets          []VpcSubnet
	DatabaseSubnetsGroupName string
}

// VpcSubnet represents an aws vpc subnet
type VpcSubnet struct {
	ID               string
	Cidr             string
	AvailabilityZone string
}

// CreateVpcOpts defines the inputs to create a vpc
type CreateVpcOpts struct {
	ID      ID
	Cidr    string
	Minimal bool
}

// DeleteVpcOpts defines the inputs to delete a vpc
type DeleteVpcOpts struct {
	ID ID
}

// Validate a vpc create request
func (o CreateVpcOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Cidr, validation.Required),
	)
}

// VpcService defines the service layer of a vpc
type VpcService interface {
	CreateVpc(ctx context.Context, opts CreateVpcOpts) (*Vpc, error)
	DeleteVpc(ctx context.Context, opts DeleteVpcOpts) error
}

// VpcCloudProvider defines the cloud actions that a Vpc service requires
type VpcCloudProvider interface {
	CreateVpc(ctx context.Context, opts CreateVpcOpts) (*Vpc, error)
	DeleteVpc(opts DeleteVpcOpts) error
}
