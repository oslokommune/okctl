package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/credentials/aws"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/kube"
	"github.com/oslokommune/okctl/pkg/kube/manifests/awsnode"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type clusterService struct {
	spinner spinner.Spinner
	api     client.ClusterAPI
	store   client.ClusterStore
	report  client.ClusterReport
	state   client.ClusterState

	provider v1alpha1.CloudProvider
	auth     aws.Authenticator
}

func (s *clusterService) CreateCluster(_ context.Context, opts api.ClusterCreateOpts) (*api.Cluster, error) {
	err := s.spinner.Start("cluster")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	cluster, err := s.api.CreateCluster(opts)
	if err != nil {
		return nil, err
	}

	clientSet, _, err := kube.NewFromEKSCluster(cluster.ID.ClusterName, cluster.ID.Region, s.provider, s.auth).Get()
	if err != nil {
		return nil, err
	}

	err = awsnode.New(clientSet).EnablePodENI()
	if err != nil {
		return nil, err
	}

	r1, err := s.store.SaveCluster(cluster)
	if err != nil {
		return nil, err
	}

	r2, err := s.state.SaveCluster(cluster)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreateCluster(cluster, []*store.Report{r1, r2})
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (s *clusterService) DeleteCluster(_ context.Context, opts api.ClusterDeleteOpts) error {
	err := s.spinner.Start("cluster")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = s.api.DeleteCluster(opts)
	if err != nil {
		return err
	}

	_, err = s.store.DeleteCluster(opts.ID)
	if err != nil {
		return err
	}

	_, err = s.state.DeleteCluster(opts.ID)
	if err != nil {
		return err
	}

	return nil
}

// NewClusterService returns an initialised cluster service
func NewClusterService(
	spinner spinner.Spinner,
	api client.ClusterAPI,
	store client.ClusterStore,
	report client.ClusterReport,
	state client.ClusterState,
	provider v1alpha1.CloudProvider,
	auth aws.Authenticator,
) client.ClusterService {
	return &clusterService{
		spinner:  spinner,
		api:      api,
		store:    store,
		report:   report,
		state:    state,
		provider: provider,
		auth:     auth,
	}
}
