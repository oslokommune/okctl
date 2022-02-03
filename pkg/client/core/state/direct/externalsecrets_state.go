package direct

import (
	"context"
	"fmt"

	merrors "github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/helm/charts/externalsecrets"
)

type externalSecretsState struct {
	clusterMeta v1alpha1.ClusterMeta
	helm        api.HelmService
}

// HasExternalSecrets returns a boolean indicating if the resource exists
func (e externalSecretsState) HasExternalSecrets() (bool, error) {
	ch := externalsecrets.New(
		externalsecrets.NewDefaultValues(e.clusterMeta.Region),
		constant.DefaultChartApplyTimeout,
	)

	_, err := e.helm.GetHelmRelease(context.Background(), api.GetHelmReleaseOpts{
		ClusterID: api.ID{
			Region:       e.clusterMeta.Region,
			AWSAccountID: e.clusterMeta.AccountID,
			ClusterName:  e.clusterMeta.Name,
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

// NewExternalSecretsState returns an initialized state client
func NewExternalSecretsState(clusterMeta v1alpha1.ClusterMeta, helmService api.HelmService) client.ExternalSecretsState {
	return &externalSecretsState{
		clusterMeta: clusterMeta,
		helm:        helmService,
	}
}
