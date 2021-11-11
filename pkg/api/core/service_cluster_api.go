// Package core implements the service layer
package core

import (
	"context"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

type clusterService struct {
	run           api.ClusterRun
	cloudProvider v1alpha1.CloudProvider
}

func (s *clusterService) CreateCluster(_ context.Context, opts api.ClusterCreateOpts) (*api.Cluster, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	cluster, err := s.run.CreateCluster(opts)
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

// GetClusterSecurityGroupID returns the EKS cluster's security group ID.
// See https://docs.aws.amazon.com/eks/latest/userguide/sec-group-reqs.html
func (s *clusterService) GetClusterSecurityGroupID(
	ctx context.Context, opts api.ClusterSecurityGroupIDGetOpts,
) (*api.ClusterSecurityGroupID, error) {
	cluster, err := s.cloudProvider.EKS().DescribeClusterWithContext(ctx, &eks.DescribeClusterInput{
		Name: &opts.ID.ClusterName,
	})
	if err != nil {
		return nil, errors.E(err, "describing cluster", errors.Internal)
	}

	return &api.ClusterSecurityGroupID{
		Value: *cluster.Cluster.ResourcesVpcConfig.ClusterSecurityGroupId,
	}, nil
}

// NewClusterService returns a service operator for the clusterService operations
func NewClusterService(run api.ClusterRun, cloudProvider v1alpha1.CloudProvider) api.ClusterService {
	return &clusterService{
		run:           run,
		cloudProvider: cloudProvider,
	}
}
