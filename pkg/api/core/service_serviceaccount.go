// Package core implements the service layer
package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"
)

type serviceAccount struct {
	run   api.ServiceAccountRun
	store api.ServiceAccountStore
}

var errCreateServiceAccount = func(err error) error {
	return errors.E(err, "failed to create service account", errors.Internal)
}

var errInvalidInputs = func(err error) error {
	return errors.E(err, "failed to validate inputs", errors.Invalid)
}

var errBuildServiceAccount = func(err error) error {
	return errors.E(err, "failed to build service account config", errors.Internal)
}

var errStoreServiceAccount = func(err error) error {
	return errors.E(err, "failed to store service account", errors.Internal)
}

func (c *serviceAccount) CreateExternalDNSServiceAccount(_ context.Context, opts api.CreateExternalDNSServiceAccountOpts) (*api.ServiceAccount, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errInvalidInputs(err)
	}

	config, err := clusterconfig.NewExternalDNSServiceAccount(
		opts.ID.ClusterName,
		opts.ID.Region,
		opts.PolicyArn,
		v1alpha1.PermissionsBoundaryARN(opts.ID.AWSAccountID),
	)
	if err != nil {
		return nil, errBuildServiceAccount(err)
	}

	account, err := c.createServiceAccount(opts.CreateServiceAccountOpts, config)
	if err != nil {
		return nil, errCreateServiceAccount(err)
	}

	err = c.store.SaveExternalDNSServiceAccount(account)
	if err != nil {
		return nil, errStoreServiceAccount(err)
	}

	return account, nil
}

func (c *serviceAccount) CreateAlbIngressControllerServiceAccount(_ context.Context, opts api.CreateAlbIngressControllerServiceAccountOpts) (*api.ServiceAccount, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errInvalidInputs(err)
	}

	config, err := clusterconfig.NewAlbIngressControllerServiceAccount(
		opts.ID.ClusterName,
		opts.ID.Region,
		opts.PolicyArn,
		v1alpha1.PermissionsBoundaryARN(opts.ID.AWSAccountID),
	)
	if err != nil {
		return nil, errBuildServiceAccount(err)
	}

	account, err := c.createServiceAccount(opts.CreateServiceAccountOpts, config)
	if err != nil {
		return nil, errCreateServiceAccount(err)
	}

	err = c.store.SaveAlbIngressControllerServiceAccount(account)
	if err != nil {
		return nil, errStoreServiceAccount(err)
	}

	return account, nil
}

func (c *serviceAccount) CreateExternalSecretsServiceAccount(_ context.Context, opts api.CreateExternalSecretsServiceAccountOpts) (*api.ServiceAccount, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errInvalidInputs(err)
	}

	config, err := clusterconfig.NewExternalSecretsServiceAccount(
		opts.ID.ClusterName,
		opts.ID.Region,
		opts.PolicyArn,
		v1alpha1.PermissionsBoundaryARN(opts.ID.AWSAccountID),
	)
	if err != nil {
		return nil, errBuildServiceAccount(err)
	}

	account, err := c.createServiceAccount(opts.CreateServiceAccountOpts, config)
	if err != nil {
		return nil, errCreateServiceAccount(err)
	}

	err = c.store.SaveExternalSecretsServiceAccount(account)
	if err != nil {
		return nil, errStoreServiceAccount(err)
	}

	return account, nil
}

func (c *serviceAccount) createServiceAccount(opts api.CreateServiceAccountOpts, config *api.ClusterConfig) (*api.ServiceAccount, error) {
	err := c.run.CreateServiceAccount(config)
	if err != nil {
		return nil, errors.E(err, "failed to create service account", errors.Internal)
	}

	account := &api.ServiceAccount{
		ID:        opts.ID,
		PolicyArn: opts.PolicyArn,
		Config:    config,
	}

	return account, nil
}

// NewServiceAccountService returns a service operator for the service account operations
func NewServiceAccountService(store api.ServiceAccountStore, run api.ServiceAccountRun) api.ServiceAccountService {
	return &serviceAccount{
		run:   run,
		store: store,
	}
}
