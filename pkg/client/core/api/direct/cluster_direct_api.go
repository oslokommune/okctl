package direct

import (
	"context"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)


type clusterDirectClient struct {
	service api.ClusterService
}

func (s *clusterDirectClient) CreateCluster(opts api.ClusterCreateOpts) (*api.Cluster, error) {
	return s.service.CreateCluster(context.Background(), opts)
}

func (s *clusterDirectClient) DeleteCluster(opts api.ClusterDeleteOpts) error {
	return s.service.DeleteCluster(context.Background(), opts)
}

func (s *clusterDirectClient) GetClusterSecurityGroupID(opts api.ClusterSecurityGroupIDGetOpts) (*api.ClusterSecurityGroupID, error) {
	return s.service.GetClusterSecurityGroupID(context.Background(), &opts)
}

// NewClusterAPI returns an initialised cluster API that uses core service directly
func NewClusterAPI(service api.ClusterService) client.ClusterAPI {
	return &clusterDirectClient{
		service: service,
	}
}
