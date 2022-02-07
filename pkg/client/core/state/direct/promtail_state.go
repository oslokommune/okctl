package direct

import (
	"context"
	"fmt"
	"time"

	"github.com/oslokommune/okctl/pkg/helm/charts/promtail"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
)

type promtailState struct {
	helm        api.HelmService
	clusterMeta v1alpha1.ClusterMeta
}

func (l *promtailState) HasPromtail() (bool, error) {
	ch := promtail.New(promtail.NewDefaultValues(), 0*time.Minute)

	_, err := l.helm.GetHelmRelease(context.Background(), api.GetHelmReleaseOpts{
		ClusterID: api.ID{
			Region:       l.clusterMeta.Region,
			AWSAccountID: l.clusterMeta.AccountID,
			ClusterName:  l.clusterMeta.Name,
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

// NewPromtailState returns an initialized state client
func NewPromtailState(clusterMeta v1alpha1.ClusterMeta, helmService api.HelmService) client.PromtailState {
	return &promtailState{
		clusterMeta: clusterMeta,
		helm:        helmService,
	}
}
