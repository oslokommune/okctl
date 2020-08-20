// Package core implements the service layer
package core

import (
	"context"
	"fmt"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clusterconfig"
)

type serviceAccount struct {
	run   api.ServiceAccountRun
	store api.ServiceAccountStore
}

func (c *serviceAccount) CreateExternalSecretsServiceAccount(_ context.Context, opts api.CreateExternalSecretsServiceAccountOpts) (*api.ServiceAccount, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "failed to validate service account input", errors.Invalid)
	}

	config, err := clusterconfig.NewExternalSecretsServiceAccount(
		opts.ClusterName,
		opts.Region,
		opts.PolicyArn,
		v1alpha1.PermissionsBoundaryARN(opts.AWSAccountID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create configuration for service account: %w", err)
	}

	err = c.run.CreateExternalSecretsServiceAccount(config)
	if err != nil {
		return nil, errors.E(err, "failed to create service account")
	}

	account := &api.ServiceAccount{
		ClusterName:  opts.ClusterName,
		Environment:  opts.Environment,
		Region:       opts.Region,
		AWSAccountID: opts.AWSAccountID,
		PolicyArn:    opts.PolicyArn,
		Config:       config,
	}

	err = c.store.SaveExternalSecretsServiceAccount(account)
	if err != nil {
		return nil, errors.E(err, "failed to store service account")
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
