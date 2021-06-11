package direct

import (
	"fmt"
	"time"

	"github.com/oslokommune/okctl/pkg/helm/charts/tempo"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
)

type tempoState struct {
	helm        client.HelmAPI
	clusterMeta v1alpha1.ClusterMeta
}

func (l *tempoState) HasTempo() (bool, error) {
	ch := tempo.New(tempo.NewDefaultValues(), 0*time.Minute)

	_, err := l.helm.GetHelmRelease(api.GetHelmReleaseOpts{
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

// NewTempoState returns an initialized state client
func NewTempoState(clusterMeta v1alpha1.ClusterMeta, helmClient client.HelmAPI) client.TempoState {
	return &tempoState{
		clusterMeta: clusterMeta,
		helm:        helmClient,
	}
}
