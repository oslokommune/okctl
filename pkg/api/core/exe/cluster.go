// Package exe implements the exe layer
package exe

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/binaries/run/eksctl"
)

type cluster struct {
	provider binaries.Provider
}

// CreateCluster invokes a CLI for performing create
func (c *cluster) CreateCluster(config *api.ClusterConfig) error {
	cli, err := c.provider.Eksctl(eksctl.Version)
	if err != nil {
		return err
	}

	_, err = cli.CreateCluster(config)

	return err
}

// DeleteClusterConfig invokes a CLI for performing delete
func (c *cluster) DeleteCluster(config *api.ClusterConfig) error {
	cli, err := c.provider.Eksctl(eksctl.Version)
	if err != nil {
		return err
	}

	_, err = cli.DeleteCluster(config)

	return err
}

// NewClusterExe returns a executor for cluster
func NewClusterExe(provider binaries.Provider) api.ClusterExe {
	return &cluster{
		provider: provider,
	}
}
