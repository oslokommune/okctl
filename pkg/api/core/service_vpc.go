package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type vpc struct {
	cloud api.VpcCloudProvider
	store api.VpcStore
}

// CreateVpc implements the business logic for creating a vpc
func (v *vpc) CreateVpc(_ context.Context, opts api.CreateVpcOpts) (*api.Vpc, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate create vpc options", errors.Invalid)
	}

	got, err := v.cloud.CreateVpc(opts)
	if err != nil {
		return nil, err
	}

	err = v.store.SaveVpc(got)
	if err != nil {
		return nil, err
	}

	return got, nil
}

// DeleteVpc implements the business logic for deleting a vpc
func (v *vpc) DeleteVpc(_ context.Context, opts api.DeleteVpcOpts) error {
	err := v.cloud.DeleteVpc(opts)
	if err != nil {
		return err
	}

	return v.store.DeleteVpc(opts.ID.Environment, opts.ID.Repository)
}

// NewVpcService returns an instantiated vpc service
func NewVpcService(cloud api.VpcCloudProvider, store api.VpcStore) api.VpcService {
	return &vpc{
		cloud: cloud,
		store: store,
	}
}
