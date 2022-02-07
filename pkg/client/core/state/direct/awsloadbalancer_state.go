package direct

import (
	"context"
	"fmt"
	"time"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/helm/charts/awslbc"
)

type awsLoadBalancerState struct {
	clusterMeta v1alpha1.ClusterMeta
	helm        api.HelmService
}

// HasAWSLoadBalancerController returns a boolean indicating if the resource exists
func (a awsLoadBalancerState) HasAWSLoadBalancerController() (bool, error) {
	ch := awslbc.New(awslbc.NewDefaultValues(a.clusterMeta.Name, "", a.clusterMeta.Region), 0*time.Minute)

	_, err := a.helm.GetHelmRelease(context.Background(), api.GetHelmReleaseOpts{
		ClusterID: api.ID{
			Region:       a.clusterMeta.Region,
			AWSAccountID: a.clusterMeta.AccountID,
			ClusterName:  a.clusterMeta.Name,
		},
		ReleaseName: ch.ReleaseName,
		Namespace:   ch.Namespace,
	})
	if err != nil {
		if errors.IsKind(err, errors.NotExist) {
			return false, nil
		}

		return false, fmt.Errorf("getting helm release: %w", err)
	}

	return true, nil
}

// NewAWSLoadBalancerState returns an initialized state client
func NewAWSLoadBalancerState(clusterMeta v1alpha1.ClusterMeta, helmService api.HelmService) client.AWSLoadBalancerControllerState {
	return &awsLoadBalancerState{
		clusterMeta: clusterMeta,
		helm:        helmService,
	}
}
