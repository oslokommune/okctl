// Package store implements the store layer
package store

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/oslokommune/okctl/pkg/storage/state"
)

type cluster struct {
	provider state.PersisterProvider
}

// SaveClusterConfig knows how to save cluster state
func (c *cluster) SaveCluster(clu *api.Cluster) error {
	s := c.provider.Repository().State()

	s.Clusters = append(s.Clusters, repository.Cluster{
		Environment: clu.Environment,
		AWS: repository.AWS{
			AccountID: clu.AWSAccountID,
			Cidr:      clu.Cidr,
		},
	})

	return c.provider.Repository().SaveState()
}

// DeleteClusterConfig knows how to delete cluster state
func (c *cluster) DeleteCluster(env string) error {
	s := c.provider.Repository().State()

	for i, curr := range s.Clusters {
		if curr.Environment == env {
			s.Clusters[i] = s.Clusters[len(s.Clusters)-1]
			s.Clusters[len(s.Clusters)-1] = repository.Cluster{}
			s.Clusters = s.Clusters[:len(s.Clusters)-1]
		}
	}

	return c.provider.Repository().SaveState()
}

// GetClusterConfig knows how to get cluster state
func (c *cluster) GetCluster(env string) (*api.Cluster, error) {
	clusters := c.provider.Repository().State().Clusters

	for _, cluster := range clusters {
		if cluster.Environment == env {
			return &api.Cluster{
				Environment:  cluster.Environment,
				AWSAccountID: cluster.AWS.AccountID,
				Cidr:         cluster.AWS.Cidr,
			}, nil
		}
	}

	return nil, fmt.Errorf("failed to find cluster configuration for env: %s", env)
}

// NewClusterStore returns a store for cluster
func NewClusterStore(provider state.PersisterProvider) api.ClusterStore {
	return &cluster{
		provider: provider,
	}
}
