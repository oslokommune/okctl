package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/credentials/aws"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/kube"
	"github.com/oslokommune/okctl/pkg/kube/manifests/awsnode"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type clusterService struct {
	api   client.ClusterAPI
	state client.ClusterState

	provider v1alpha1.CloudProvider
	auth     aws.Authenticator
}

func (s *clusterService) CreateCluster(_ context.Context, opts client.ClusterCreateOpts) (*client.Cluster, error) {
	c, err := s.api.CreateCluster(api.ClusterCreateOpts{
		ID:      opts.ID,
		Cidr:    opts.Cidr,
		Version: opts.Version,
		VpcID:   opts.VpcID,
		VpcPrivateSubnets: func() (subs []api.VpcSubnet) {
			for _, sub := range opts.VpcPrivateSubnets {
				subs = append(subs, api.VpcSubnet{
					ID:               sub.ID,
					Cidr:             sub.Cidr,
					AvailabilityZone: sub.AvailabilityZone,
				})
			}
			return subs
		}(),
		VpcPublicSubnets: func() (subs []api.VpcSubnet) {
			for _, sub := range opts.VpcPublicSubnets {
				subs = append(subs, api.VpcSubnet{
					ID:               sub.ID,
					Cidr:             sub.Cidr,
					AvailabilityZone: sub.AvailabilityZone,
				})
			}
			return subs
		}(),
	})
	if err != nil {
		return nil, err
	}

	clientSet, _, err := kube.NewFromEKSCluster(c.ID.ClusterName, c.ID.Region, s.provider, s.auth).Get()
	if err != nil {
		return nil, err
	}

	err = awsnode.New(clientSet).EnablePodENI()
	if err != nil {
		return nil, err
	}

	cluster := &client.Cluster{
		ID:     c.ID,
		Name:   c.ID.ClusterName,
		Config: c.Config,
	}

	err = s.state.SaveCluster(cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (s *clusterService) DeleteCluster(_ context.Context, opts client.ClusterDeleteOpts) error {
	err := s.api.DeleteCluster(api.ClusterDeleteOpts{
		ID:                 opts.ID,
		FargateProfileName: opts.FargateProfileName,
	})
	if err != nil {
		return err
	}

	err = s.state.RemoveCluster(opts.ID.ClusterName)
	if err != nil {
		return err
	}

	return nil
}

// NewClusterService returns an initialised cluster service
func NewClusterService(
	api client.ClusterAPI,
	state client.ClusterState,
	provider v1alpha1.CloudProvider,
	auth aws.Authenticator,
) client.ClusterService {
	return &clusterService{
		api:      api,
		state:    state,
		provider: provider,
		auth:     auth,
	}
}
