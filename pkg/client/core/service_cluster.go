package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type clusterService struct {
	api    client.ClusterAPI
	store  client.ClusterStore
	report client.ClusterReport
}

func (c *clusterService) CreateCluster(_ context.Context, opts api.ClusterCreateOpts) (*api.Cluster, error) {
	cluster, err := c.api.CreateCluster(opts)
	if err != nil {
		return nil, err
	}

	report, err := c.store.SaveCluster(cluster)
	if err != nil {
		return nil, err
	}

	err = c.report.ReportCreateCluster(cluster, report)
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

	return nil
}

// NewClusterService returns an initialised cluster service
func NewClusterService(api client.ClusterAPI, store client.ClusterStore, report client.ClusterReport) client.ClusterService {
	return &clusterService{
		api:    api,
		store:  store,
		report: report,
	}
}
