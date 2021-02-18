// Package core implements the service layer
package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"
)

type serviceAccount struct {
	run api.ServiceAccountRun
}

var errDeleteServiceAccount = func(err error) error {
	return errors.E(err, "failed to delete service account", errors.Internal)
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

func (c *serviceAccount) CreateBlockstorageServiceAccount(_ context.Context, opts api.CreateBlockstorageServiceAccountOpts) (*api.ServiceAccount, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errInvalidInputs(err)
	}

	config, err := clusterconfig.NewBlockstorageServiceAccount(
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

	return account, nil
}

func (c *serviceAccount) DeleteBlockstorageServiceAccount(_ context.Context, id api.ID) error {
	err := id.Validate()
	if err != nil {
		return errInvalidInputs(err)
	}

	config, err := clusterconfig.NewBlockstorageServiceAccount(
		id.ClusterName,
		id.Region,
		"n/a",
		v1alpha1.PermissionsBoundaryARN(id.AWSAccountID),
	)
	if err != nil {
		return errBuildServiceAccount(err)
	}

	err = c.run.DeleteServiceAccount(config)
	if err != nil {
		return errDeleteServiceAccount(err)
	}

	return nil
}

func (c *serviceAccount) CreateAutoscalerServiceAccount(_ context.Context, opts api.CreateAutoscalerServiceAccountOpts) (*api.ServiceAccount, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errInvalidInputs(err)
	}

	config, err := clusterconfig.NewAutoscalerServiceAccount(
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

	return account, nil
}

func (c *serviceAccount) DeleteAutoscalerServiceAccount(_ context.Context, id api.ID) error {
	err := id.Validate()
	if err != nil {
		return errInvalidInputs(err)
	}

	config, err := clusterconfig.NewAutoscalerServiceAccount(
		id.ClusterName,
		id.Region,
		"n/a",
		v1alpha1.PermissionsBoundaryARN(id.AWSAccountID),
	)
	if err != nil {
		return errBuildServiceAccount(err)
	}

	err = c.run.DeleteServiceAccount(config)
	if err != nil {
		return errDeleteServiceAccount(err)
	}

	return nil
}

func (c *serviceAccount) DeleteExternalSecretsServiceAccount(_ context.Context, id api.ID) error {
	err := id.Validate()
	if err != nil {
		return errInvalidInputs(err)
	}

	config, err := clusterconfig.NewExternalSecretsServiceAccount(
		id.ClusterName,
		id.Region,
		"n/a",
		v1alpha1.PermissionsBoundaryARN(id.AWSAccountID),
	)
	if err != nil {
		return errBuildServiceAccount(err)
	}

	err = c.run.DeleteServiceAccount(config)
	if err != nil {
		return errDeleteServiceAccount(err)
	}

	return nil
}

func (c *serviceAccount) DeleteAlbIngressControllerServiceAccount(_ context.Context, id api.ID) error {
	err := id.Validate()
	if err != nil {
		return errInvalidInputs(err)
	}

	config, err := clusterconfig.NewAlbIngressControllerServiceAccount(
		id.ClusterName,
		id.Region,
		"n/a",
		v1alpha1.PermissionsBoundaryARN(id.AWSAccountID),
	)
	if err != nil {
		return errBuildServiceAccount(err)
	}

	err = c.run.DeleteServiceAccount(config)
	if err != nil {
		return errDeleteServiceAccount(err)
	}

	return nil
}

func (c *serviceAccount) DeleteExternalDNSServiceAccount(_ context.Context, id api.ID) error {
	err := id.Validate()
	if err != nil {
		return errInvalidInputs(err)
	}

	config, err := clusterconfig.NewExternalDNSServiceAccount(
		id.ClusterName,
		id.Region,
		"n/a",
		v1alpha1.PermissionsBoundaryARN(id.AWSAccountID),
	)
	if err != nil {
		return errBuildServiceAccount(err)
	}

	err = c.run.DeleteServiceAccount(config)
	if err != nil {
		return errDeleteServiceAccount(err)
	}

	return nil
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

	return account, nil
}

// nolint: lll
func (c *serviceAccount) CreateAWSLoadBalancerControllerServiceAccount(_ context.Context, opts api.CreateAWSLoadBalancerControllerServiceAccountOpts) (*api.ServiceAccount, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errInvalidInputs(err)
	}

	config, err := clusterconfig.NewAWSLoadBalancerControllerServiceAccount(
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

	return account, nil
}

func (c *serviceAccount) DeleteAWSLoadBalancerControllerServiceAccount(_ context.Context, id api.ID) error {
	err := id.Validate()
	if err != nil {
		return errInvalidInputs(err)
	}

	config, err := clusterconfig.NewAWSLoadBalancerControllerServiceAccount(
		id.ClusterName,
		id.Region,
		"n/a",
		v1alpha1.PermissionsBoundaryARN(id.AWSAccountID),
	)
	if err != nil {
		return errBuildServiceAccount(err)
	}

	err = c.run.DeleteServiceAccount(config)
	if err != nil {
		return errDeleteServiceAccount(err)
	}

	return nil
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

	return account, nil
}

func (c *serviceAccount) createServiceAccount(opts api.CreateServiceAccountOpts, config *v1alpha5.ClusterConfig) (*api.ServiceAccount, error) {
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
func NewServiceAccountService(run api.ServiceAccountRun) api.ServiceAccountService {
	return &serviceAccount{
		run: run,
	}
}
