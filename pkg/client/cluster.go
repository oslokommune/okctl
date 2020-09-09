package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// ClusterService orchestrates the creation of a cluster
type ClusterService interface {
	CreateCluster(context.Context, api.ClusterCreateOpts) (*api.Cluster, error)
	DeleteCluster(context.Context, api.ClusterDeleteOpts) error
}

// ClusterAPI invokes the API calls for creating a cluster
type ClusterAPI interface {
	CreateCluster(opts api.ClusterCreateOpts) (*api.Cluster, error)
	DeleteCluster(opts api.ClusterDeleteOpts) error
}

// ClusterStore stores the data
type ClusterStore interface {
	SaveCluster(cluster *api.Cluster) (*store.Report, error)
	DeleteCluster(id api.ID) (*store.Report, error)
	GetCluster(id api.ID) (*api.Cluster, error)
}

// ClusterReport provides output of the actions
type ClusterReport interface {
	ReportCreateCluster(cluster *api.Cluster, report *store.Report) error
}
