package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type vpc struct {
	cloud api.VpcCloudProvider
}

// CreateVpc implements the business logic for creating a vpc
func (v *vpc) CreateVpc(_ context.Context, opts api.CreateVpcOpts) (*api.Vpc, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating the inputs", errors.Invalid)
	}

	got, err := v.cloud.CreateVpc(opts)
	if err != nil {
		return nil, errors.E(err, "creating vpc", errors.Internal)
	}

	return got, nil
}

// DeleteVpc implements the business logic for deleting a vpc
func (v *vpc) DeleteVpc(_ context.Context, opts api.DeleteVpcOpts) error {
	err := v.cloud.DeleteVpc(opts)
	if err != nil {
		return errors.E(err, "deleting vpc", errors.Internal)
	}

	return nil
}

// NewVpcService returns an instantiated vpc service
func NewVpcService(cloud api.VpcCloudProvider) api.VpcService {
	return &vpc{
		cloud: cloud,
	}
}
