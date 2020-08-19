package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Vpc represents the state of an aws vpc
type Vpc struct {
	StackName              string
	CloudFormationTemplate []byte

	ID             string
	PublicSubnets  []VpcSubnet
	PrivateSubnets []VpcSubnet
}

// VpcSubnet represents an aws vpc subnet
type VpcSubnet struct {
	ID               string
	Cidr             string
	AvailabilityZone string
}

// CreateVpcOpts defines the inputs to create a vpc
type CreateVpcOpts struct {
	AwsAccountID string
	ClusterName  string // not needed
	Env          string
	RepoName     string
	Cidr         string
	Region       string
}

// DeleteVpcOpts defines the inputs to delete a vpc
type DeleteVpcOpts struct {
	Env      string
	RepoName string
}

// Validate a vpc create request
func (o CreateVpcOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Env, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.Cidr, validation.Required),
		validation.Field(&o.AwsAccountID, validation.Required),
		validation.Field(&o.RepoName, validation.Required),
	)
}

// VpcService defines the service layer of a vpc
type VpcService interface {
	CreateVpc(ctx context.Context, opts CreateVpcOpts) (*Vpc, error)
	DeleteVpc(ctx context.Context, opts DeleteVpcOpts) error
}

// VpcCloudProvider defines the cloud actions that a Vpc service requires
type VpcCloudProvider interface {
	CreateVpc(opts CreateVpcOpts) (*Vpc, error)
	DeleteVpc(opts DeleteVpcOpts) error
}

// VpcStore defines the storage actions that a vpc service requires
type VpcStore interface {
	SaveVpc(vpc *Vpc) error
	DeleteVpc(env, repoName string) error
	GetVpc() (*Vpc, error)
}
