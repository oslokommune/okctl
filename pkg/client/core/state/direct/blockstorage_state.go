package direct

import (
	"context"
	"fmt"

	merrors "github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/helm/charts/blockstorage"
)

type blockstorageState struct {
	clusterMeta v1alpha1.ClusterMeta
	helm        api.HelmService
}

// HasBlockstorage returns a boolean indicating if the resource exists
func (b blockstorageState) HasBlockstorage() (bool, error) {
	ch := blockstorage.New(
		blockstorage.NewDefaultValues(b.clusterMeta.Region, b.clusterMeta.Name, ""),
		constant.DefaultChartApplyTimeout,
	)

	_, err := b.helm.GetHelmRelease(context.Background(), api.GetHelmReleaseOpts{
		ClusterID: api.ID{
			Region:       b.clusterMeta.Region,
			AWSAccountID: b.clusterMeta.AccountID,
			ClusterName:  b.clusterMeta.Name,
		},
		ReleaseName: ch.ReleaseName,
		Namespace:   ch.Namespace,
	})
	if err != nil {
		if merrors.IsKind(err, merrors.NotExist) {
			return false, nil
		}

		return false, fmt.Errorf("acquiring Helm release: %w", err)
	}

	return true, nil
}

// NewBlockstorageState returns an initialized state client
func NewBlockstorageState(clusterMeta v1alpha1.ClusterMeta, helmService api.HelmService) client.BlockstorageState {
	return &blockstorageState{
		clusterMeta: clusterMeta,
		helm:        helmService,
	}
}
