// Package store implements the store layer
package store

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/oslokommune/okctl/pkg/storage/state"
	"sigs.k8s.io/yaml"
)

type cluster struct {
	provider state.PersisterProvider
}

// SaveCluster knows how to save cluster state
func (c *cluster) SaveCluster(clu *api.Cluster) error {
	s := c.provider.Repository().State()

	s.Clusters = append(s.Clusters, repository.Cluster{
		Environment: clu.Environment,
		AWS: repository.AWS{
			AccountID: clu.AWSAccountID,
			Cidr:      clu.Cidr,
		},
	})

	err := c.provider.Repository().SaveState()
	if err != nil {
		return err
	}

	data, err := clu.Config.YAML()
	if err != nil {
		return err
	}

	return c.provider.Repository().WriteToDefault("cluster_config", data)
}

// DeleteCluster knows how to delete cluster state
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

// GetCluster knows how to get cluster state
func (c *cluster) GetCluster(env string) (*api.Cluster, error) {
	data, err := c.provider.Repository().ReadFromDefault("cluster_config")
	if err != nil {
		return nil, err
	}

	cfg := v1alpha1.NewClusterConfig()

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	clusters := c.provider.Repository().State().Clusters

	for _, cluster := range clusters {
		if cluster.Environment == env {
			return &api.Cluster{
				Environment:  cluster.Environment,
				AWSAccountID: cluster.AWS.AccountID,
				Cidr:         cluster.AWS.Cidr,
				Config:       cfg,
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
