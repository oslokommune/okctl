package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/oslokommune/okctl/pkg/api"
)

// Cluster contains the core state for a cluster
type Cluster struct {
	ID     api.ID
	Name   string
	Config *v1alpha5.ClusterConfig
}

// ClusterCreateOpts specifies the required inputs for creating a cluster
type ClusterCreateOpts struct {
	ID                api.ID
	Cidr              string
	Version           string
	VpcID             string
	VpcPrivateSubnets []VpcSubnet
	VpcPublicSubnets  []VpcSubnet
}

// ClusterDeleteOpts specifies the required inputs for deleting a cluster
type ClusterDeleteOpts struct {
	ID                 api.ID
	FargateProfileName string
}

// ClusterService orchestrates the creation of a cluster
type ClusterService interface {
	CreateCluster(context.Context, ClusterCreateOpts) (*Cluster, error)
	DeleteCluster(context.Context, ClusterDeleteOpts) error
}

// ClusterAPI invokes the API calls for creating a cluster
type ClusterAPI interface {
	CreateCluster(opts api.ClusterCreateOpts) (*api.Cluster, error)
	DeleteCluster(opts api.ClusterDeleteOpts) error
}

// ClusterState implements the state layer
type ClusterState interface {
	SaveCluster(cluster *Cluster) error
	GetCluster(name string) (*Cluster, error)
	RemoveCluster(name string) error
}
