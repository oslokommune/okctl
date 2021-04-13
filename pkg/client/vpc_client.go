package client

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/oslokommune/okctl/pkg/api"
)

// Vpc represents the state of an aws vpc
type Vpc struct {
	ID                     api.ID
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
	ID      api.ID
	Cidr    string
	Minimal bool
}

// DeleteVpcOpts defines the inputs to delete a vpc
type DeleteVpcOpts struct {
	ID api.ID
}

// Validate a vpc create request
func (o CreateVpcOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Cidr, validation.Required),
	)
}

// VPCService orchestrates the creation of a vpc
type VPCService interface {
	CreateVpc(ctx context.Context, opts CreateVpcOpts) (*Vpc, error)
	DeleteVpc(ctx context.Context, opts DeleteVpcOpts) error
	GetVPC(ctx context.Context, id api.ID) (*Vpc, error)
}

// VPCAPI invokes the API calls for creating a vpc
type VPCAPI interface {
	CreateVpc(opts api.CreateVpcOpts) (*api.Vpc, error)
	DeleteVpc(opts api.DeleteVpcOpts) error
}

// VPCState implement the state layer
type VPCState interface {
	SaveVpc(vpc *Vpc) error
	GetVpc(stackName string) (*Vpc, error)
	RemoveVpc(stackName string) error
}
