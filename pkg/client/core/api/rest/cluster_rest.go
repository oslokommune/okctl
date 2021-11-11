package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// ClusterTarget is the API route for HTTP requests
const ClusterTarget = "clusters/"

// ClusterSecurityGroupIDTarget is the API route for HTTP requests
const ClusterSecurityGroupIDTarget = "clustersecuritygroupid/"

type clusterAPI struct {
	client *HTTPClient
}

func (a *clusterAPI) CreateCluster(opts api.ClusterCreateOpts) (*api.Cluster, error) {
	into := &api.Cluster{}
	return into, a.client.DoPost(ClusterTarget, &opts, into)
}

func (a *clusterAPI) DeleteCluster(opts api.ClusterDeleteOpts) error {
	return a.client.DoDelete(ClusterTarget, &opts)
}

func (a *clusterAPI) GetClusterSecurityGroupID(opts api.ClusterSecurityGroupIDGetOpts) (*api.ClusterSecurityGroupID, error) {
	id := &api.ClusterSecurityGroupID{}
	return id, a.client.DoGet(ClusterSecurityGroupIDTarget, &opts, id)
}

// NewClusterAPI returns an initialised cluster API
func NewClusterAPI(client *HTTPClient) client.ClusterAPI {
	return &clusterAPI{
		client: client,
	}
}
