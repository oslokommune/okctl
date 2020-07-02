// Package mock provides mocks for use with tests
package mock

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

const (
	// DefaultEnv is a default environment used in mocks
	DefaultEnv = "pro"
	// DefaultAWSAccountID is a default aws account id used in mocks
	DefaultAWSAccountID = "123456789012"
	// DefaultCidr is a default cidr used in mocks
	DefaultCidr = "192.168.0.0/20"
	// DefaultRegion is a default aws region used in mocks
	DefaultRegion = "eu-west-1"
	// DefaultAvailabilityZone is a default aws availability zone used in mocks
	DefaultAvailabilityZone = "eu-west-1a"
	// DefaultRepositoryName is a default git repo name used in mocks
	DefaultRepositoryName = "test"
	// DefaultClusterName is a default eks cluster name used in mocks
	DefaultClusterName = "test-cluster-pro"
	// DefaultVpcID is a default aws vpc id used in mocks
	DefaultVpcID = "vpc-0e9801d129EXAMPLE"
	// DefaultPublicSubnetID is a default aws public subnet id used in mocks
	DefaultPublicSubnetID = "subnet-0bb1c79de3EXAMPLE"
	// DefaultPublicSubnetCidr is a default public subnet cidr used in mocks
	DefaultPublicSubnetCidr = "192.168.1.0/24"
	// DefaultPrivateSubnetID is a default private aws subnet id used in mocks
	DefaultPrivateSubnetID = "subnet-8EXAMPLE"
	// DefaultPrivateSubnetCidr is a default private subnet cidr used in mocks
	DefaultPrivateSubnetCidr = "192.168.2.0/24"
)

// DefaultClusterCreateOpts returns options for creating a cluster with defaults set
func DefaultClusterCreateOpts() api.ClusterCreateOpts {
	return api.ClusterCreateOpts{
		Environment:    DefaultEnv,
		AWSAccountID:   DefaultAWSAccountID,
		Cidr:           DefaultCidr,
		RepositoryName: DefaultRepositoryName,
		Region:         DefaultRegion,
		ClusterName:    DefaultClusterName,
	}
}

// DefaultPublicSubnets returns a map of public subnets with defaults set
func DefaultPublicSubnets() map[string]v1alpha1.ClusterNetwork {
	return map[string]v1alpha1.ClusterNetwork{
		DefaultAvailabilityZone: {
			ID:   DefaultPublicSubnetID,
			CIDR: DefaultPublicSubnetCidr,
		},
	}
}

// DefaultPrivateSubnets returns a map of private subnets with defaults set
func DefaultPrivateSubnets() map[string]v1alpha1.ClusterNetwork {
	return map[string]v1alpha1.ClusterNetwork{
		DefaultAvailabilityZone: {
			ID:   DefaultPrivateSubnetID,
			CIDR: DefaultPrivateSubnetCidr,
		},
	}
}

// DefaultClusterConfig returns a cluster config with defaults set
func DefaultClusterConfig() *v1alpha1.ClusterConfig {
	cfg := v1alpha1.NewClusterConfig()

	cfg.Metadata.Name = DefaultClusterName
	cfg.Metadata.Region = DefaultRegion

	cfg.VPC.ID = DefaultVpcID
	cfg.VPC.CIDR = DefaultCidr

	cfg.VPC.Subnets.Public = DefaultPublicSubnets()
	cfg.VPC.Subnets.Private = DefaultPrivateSubnets()

	cfg.IAM.FargatePodExecutionRolePermissionsBoundary = v1alpha1.PermissionsBoundaryARN(DefaultAWSAccountID)
	cfg.IAM.ServiceRolePermissionsBoundary = v1alpha1.PermissionsBoundaryARN(DefaultAWSAccountID)

	return cfg
}

// DefaultCluster returns an api cluster definition with defaults set
func DefaultCluster() *api.Cluster {
	return &api.Cluster{
		Environment:  DefaultEnv,
		AWSAccountID: DefaultAWSAccountID,
		Cidr:         DefaultCidr,
		Config:       DefaultClusterConfig(),
	}
}

// ClusterService provides a mock for the cluster service interface
type ClusterService struct {
	CreateClusterFn func(context.Context, api.ClusterCreateOpts) (*api.Cluster, error)
	DeleteClusterFn func(context.Context, api.ClusterDeleteOpts) error
}

// CreateCluster invokes a mocked function to create a cluster
func (s *ClusterService) CreateCluster(ctx context.Context, opts api.ClusterCreateOpts) (*api.Cluster, error) {
	return s.CreateClusterFn(ctx, opts)
}

// DeleteCluster invokes a mocked function to delete a cluster
func (s *ClusterService) DeleteCluster(ctx context.Context, opts api.ClusterDeleteOpts) error {
	return s.DeleteClusterFn(ctx, opts)
}

// NewGoodClusterService returns a cluster service that will succeed
func NewGoodClusterService() *ClusterService {
	return &ClusterService{
		CreateClusterFn: func(_ context.Context, opts api.ClusterCreateOpts) (*api.Cluster, error) {
			return DefaultCluster(), nil
		},
		DeleteClusterFn: func(context.Context, api.ClusterDeleteOpts) error {
			return nil
		},
	}
}

// NewBadClusterService returns a cluster service that will fail
func NewBadClusterService() *ClusterService {
	return &ClusterService{
		CreateClusterFn: func(context.Context, api.ClusterCreateOpts) (*api.Cluster, error) {
			return nil, fmt.Errorf("something bad")
		},
		DeleteClusterFn: func(context.Context, api.ClusterDeleteOpts) error {
			return fmt.Errorf("something bad")
		},
	}
}

// ClusterCloud provides a mock for the cluster cloud interface
type ClusterCloud struct {
	CreateClusterFn func(awsAccountID, clusterName, env, repoName, cidr, region string) (*v1alpha1.ClusterConfig, error)
	DeleteClusterFn func(env, repoName string) error
}

// CreateCluster invokes the mocked create cluster function
func (c *ClusterCloud) CreateCluster(awsAccountID, clusterName, env, repoName, cidr, region string) (*v1alpha1.ClusterConfig, error) {
	return c.CreateClusterFn(awsAccountID, clusterName, env, repoName, cidr, region)
}

// DeleteCluster invokes the mocked delete cluster function
func (c *ClusterCloud) DeleteCluster(env, repoName string) error {
	return c.DeleteClusterFn(env, repoName)
}

// NewGoodClusterCloud returns a cluster cloud that will succeed
func NewGoodClusterCloud() *ClusterCloud {
	return &ClusterCloud{
		CreateClusterFn: func(awsAccountID, clusterName, env, repoName, cidr, region string) (*v1alpha1.ClusterConfig, error) {
			return DefaultClusterConfig(), nil
		},
		DeleteClusterFn: func(env, repoName string) error {
			return nil
		},
	}
}

// NewBadClusterCloud returns a cluster cloud that will fail
func NewBadClusterCloud() *ClusterCloud {
	return &ClusterCloud{
		CreateClusterFn: func(awsAccountID, clusterName, env, repoName, cidr, region string) (*v1alpha1.ClusterConfig, error) {
			return nil, fmt.Errorf("something bad")
		},
		DeleteClusterFn: func(env, repoName string) error {
			return fmt.Errorf("something bad")
		},
	}
}

// ClusterExe provides a mock for the cluster exe interface
type ClusterExe struct {
	CreateClusterFn func(*v1alpha1.ClusterConfig) error
	DeleteClusterFn func(*v1alpha1.ClusterConfig) error
}

// CreateCluster invokes the mocked create cluster function
func (c *ClusterExe) CreateCluster(config *v1alpha1.ClusterConfig) error {
	return c.CreateClusterFn(config)
}

// DeleteCluster invokes the mocked delete cluster function
func (c *ClusterExe) DeleteCluster(config *v1alpha1.ClusterConfig) error {
	return c.DeleteClusterFn(config)
}

// NewGoodClusterExe returns a cluster exe that will succeed
func NewGoodClusterExe() *ClusterExe {
	return &ClusterExe{
		CreateClusterFn: func(config *v1alpha1.ClusterConfig) error {
			return nil
		},
		DeleteClusterFn: func(config *v1alpha1.ClusterConfig) error {
			return nil
		},
	}
}

// NewBadClusterExe returns a cluster exe that will fail
func NewBadClusterExe() *ClusterExe {
	return &ClusterExe{
		CreateClusterFn: func(config *v1alpha1.ClusterConfig) error {
			return fmt.Errorf("something bad")
		},
		DeleteClusterFn: func(config *v1alpha1.ClusterConfig) error {
			return fmt.Errorf("something bad")
		},
	}
}

// ClusterStore provides a mock for the cluster store interface
type ClusterStore struct {
	SaveClusterFn   func(*api.Cluster) error
	DeleteClusterFn func(env string) error
	GetClusterFn    func(env string) (*api.Cluster, error)
}

// SaveCluster invokes the mocked save function
func (c *ClusterStore) SaveCluster(cluster *api.Cluster) error {
	return c.SaveClusterFn(cluster)
}

// DeleteCluster invokes the mocked delete function
func (c *ClusterStore) DeleteCluster(env string) error {
	return c.DeleteClusterFn(env)
}

// GetCluster invokes the mocked get function
func (c *ClusterStore) GetCluster(env string) (*api.Cluster, error) {
	return c.GetClusterFn(env)
}

// NewGoodClusterStore returns a cluster store that will succeed
func NewGoodClusterStore() *ClusterStore {
	return &ClusterStore{
		SaveClusterFn: func(cluster *api.Cluster) error {
			return nil
		},
		DeleteClusterFn: func(env string) error {
			return nil
		},
		GetClusterFn: func(env string) (*api.Cluster, error) {
			return DefaultCluster(), nil
		},
	}
}

// NewBadClusterStore returns a cluster store that will fail
func NewBadClusterStore() *ClusterStore {
	return &ClusterStore{
		SaveClusterFn: func(cluster *api.Cluster) error {
			return fmt.Errorf("something bad")
		},
		DeleteClusterFn: func(env string) error {
			return fmt.Errorf("something bad")
		},
		GetClusterFn: func(env string) (*api.Cluster, error) {
			return nil, fmt.Errorf("something bad")
		},
	}
}
