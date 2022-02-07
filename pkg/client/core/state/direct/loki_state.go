package direct

import (
	"context"
	"fmt"
	"time"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/helm/charts/loki"
)

type lokiState struct {
	helm        api.HelmService
	clusterMeta v1alpha1.ClusterMeta
}

func (l *lokiState) HasLoki() (bool, error) {
	ch := loki.New(loki.NewDefaultValues(), 0*time.Minute)

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

// NewLokiState returns an initialized state client
func NewLokiState(clusterMeta v1alpha1.ClusterMeta, helmService api.HelmService) client.LokiState {
	return &lokiState{
		clusterMeta: clusterMeta,
		helm:        helmService,
	}
}
