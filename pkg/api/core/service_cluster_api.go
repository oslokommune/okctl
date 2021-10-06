// Package core implements the service layer
package core

import (
	"context"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type clusterService struct {
	run api.ClusterRun
}

func (s *clusterService) CreateCluster(ctx context.Context, opts api.ClusterCreateOpts) (*api.Cluster, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	cluster, err := s.run.CreateCluster(ctx, opts)
	if err != nil {
		return nil, errors.E(err, "creating cluster", errors.Internal)
	}

	return cluster, nil
}

func (s *clusterService) DeleteCluster(_ context.Context, opts api.ClusterDeleteOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = s.run.DeleteCluster(opts)
	if err != nil {
		return errors.E(err, "deleting cluster", errors.Internal)
	}

	return nil
}

// NewClusterService returns a service operator for the clusterService operations
func NewClusterService(run api.ClusterRun) api.ClusterService {
	return &clusterService{
		run: run,
	}
}
