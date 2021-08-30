package direct

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"time"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/helm/charts/autoscaler"
)

type autoscalerState struct {
	clusterMeta v1alpha1.ClusterMeta
	helm        client.HelmAPI
}

func (a autoscalerState) HasAutoscaler() (bool, error) {
	ch := autoscaler.New(
		autoscaler.NewDefaultValues(a.clusterMeta.Region, a.clusterMeta.Name, ""),
		0*time.Minute,
	)

	_, err := a.helm.GetHelmRelease(api.GetHelmReleaseOpts{
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

		return false, fmt.Errorf(constant.GetHelmReleaseError, err)
	}

	return true, nil
}

// NewAutoscalerState returns an initialized state client
func NewAutoscalerState(clusterMeta v1alpha1.ClusterMeta, helmClient client.HelmAPI) client.AutoscalerState {
	return &autoscalerState{
		clusterMeta: clusterMeta,
		helm:        helmClient,
	}
}
