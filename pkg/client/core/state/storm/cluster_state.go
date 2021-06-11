package storm

import (
	"errors"
	"fmt"
	"time"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/oslokommune/okctl/pkg/breeze"
	"github.com/oslokommune/okctl/pkg/client"
)

type clusterState struct {
	node breeze.Client
}

// Cluster contains storm compatible state
type Cluster struct {
	Metadata `storm:"inline"`

	ID     ID
	Name   string `storm:"unique"`
	Config *v1alpha5.ClusterConfig
}

// NewCluster returns storm compatible state
func NewCluster(c *client.Cluster, meta Metadata) *Cluster {
	return &Cluster{
		Metadata: meta,
		ID:       NewID(c.ID),
		Name:     c.Name,
		Config:   c.Config,
	}
}

// Convert to client.Cluster
func (c *Cluster) Convert() *client.Cluster {
	return &client.Cluster{
		ID:     c.ID.Convert(),
		Name:   c.Name,
		Config: c.Config,
	}
}

func (c *clusterState) SaveCluster(cluster *client.Cluster) error {
	existing, err := c.getCluster(cluster.Name)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return c.node.Save(NewCluster(cluster, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return c.node.Save(NewCluster(cluster, existing.Metadata))
}

func (c *clusterState) GetCluster(name string) (*client.Cluster, error) {
	cluster, err := c.getCluster(name)
	if err != nil {
		return nil, err
	}

	return cluster.Convert(), nil
}

func (c *clusterState) HasCluster(name string) (bool, error) {
	_, err := c.getCluster(name)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("querying state for cluster data: %w", err)
	}

	return true, nil
}

func (c *clusterState) getCluster(name string) (*Cluster, error) {
	cluster := &Cluster{}

	err := c.node.One("Name", name, cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (c *clusterState) RemoveCluster(name string) error {
	cluster := &Cluster{}

	err := c.node.One("Name", name, cluster)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return c.node.DeleteStruct(cluster)
}

// NewClusterState returns an initialised state client
func NewClusterState(node breeze.Client) client.ClusterState {
	return &clusterState{
		node: node,
	}
}
