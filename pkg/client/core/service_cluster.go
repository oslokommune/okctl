package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type clusterService struct {
	api    client.ClusterAPI
	store  client.ClusterStore
	report client.ClusterReport
	state  client.ClusterState
}

func (c *clusterService) CreateCluster(_ context.Context, opts api.ClusterCreateOpts) (*api.Cluster, error) {
	cluster, err := c.api.CreateCluster(opts)
	if err != nil {
		return nil, err
	}

	r1, err := c.store.SaveCluster(cluster)
	if err != nil {
		return nil, err
	}

	r2, err := c.state.SaveCluster(cluster)
	if err != nil {
		return nil, err
	}

	err = c.report.ReportCreateCluster(cluster, []*store.Report{r1, r2})
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (c *clusterService) DeleteCluster(_ context.Context, opts api.ClusterDeleteOpts) error {
	err := c.api.DeleteCluster(opts)
	if err != nil {
		return err
	}

	_, err = c.store.DeleteCluster(opts.ID)
	if err != nil {
		return err
	}

	_, err = c.state.DeleteCluster(opts.ID)
	if err != nil {
		return err
	}

	return nil
}

// NewClusterService returns an initialised cluster service
func NewClusterService(api client.ClusterAPI, store client.ClusterStore, report client.ClusterReport, state client.ClusterState) client.ClusterService {
	return &clusterService{
		api:    api,
		store:  store,
		report: report,
		state:  state,
	}
}
