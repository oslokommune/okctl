package storm

import (
	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/oslokommune/okctl/pkg/client"
)

type clusterState struct {
	node stormpkg.Node
}

// Cluster contains storm compatible state
type Cluster struct {
	Metadata `storm:"inline"`

	ID     ID
	Name   string `storm:"unique,index"`
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
	return c.node.Save(NewCluster(cluster, NewMetadata()))
}

func (c *clusterState) GetCluster(name string) (*client.Cluster, error) {
	cluster := &Cluster{}

	err := c.node.One("Name", name, cluster)
	if err != nil {
		return nil, err
	}

	return cluster.Convert(), nil
}

func (c *clusterState) RemoveCluster(name string) error {
	cluster := &Cluster{}

	err := c.node.One("Name", name, cluster)
	if err != nil {
		return err
	}

	return c.node.DeleteStruct(cluster)
}

// NewClusterState returns an initialised state client
func NewClusterState(node stormpkg.Node) client.ClusterState {
	return &clusterState{
		node: node,
	}
}
