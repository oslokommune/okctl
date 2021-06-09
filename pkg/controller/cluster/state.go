package cluster

import (
	"errors"

	"github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common"
	"github.com/oslokommune/okctl/pkg/helm/charts/autoscaler"
	"github.com/oslokommune/okctl/pkg/helm/charts/awslbc"
	"github.com/oslokommune/okctl/pkg/helm/charts/blockstorage"
	"github.com/oslokommune/okctl/pkg/helm/charts/externalsecrets"
	"github.com/oslokommune/okctl/pkg/helm/charts/kubepromstack"
	"github.com/oslokommune/okctl/pkg/helm/charts/loki"
	"github.com/oslokommune/okctl/pkg/helm/charts/promtail"
	"github.com/oslokommune/okctl/pkg/helm/charts/tempo"
)

// ExistingResources contains information about what services already exists in a cluster
type ExistingResources struct {
	hasServiceQuotaCheck                  bool
	hasAWSLoadBalancerController          bool
	hasCluster                            bool
	hasExternalDNS                        bool
	hasExternalSecrets                    bool
	hasAutoscaler                         bool
	hasBlockstorage                       bool
	hasKubePromStack                      bool
	hasLoki                               bool
	hasPromtail                           bool
	hasTempo                              bool
	hasIdentityManager                    bool
	hasArgoCD                             bool
	hasPrimaryHostedZone                  bool
	hasVPC                                bool
	hasDelegatedHostedZoneNameservers     bool
	hasDelegatedHostedZoneNameserversTest bool
	hasUsers                              bool
	hasPostgres                           map[string]*v1alpha1.ClusterDatabasesPostgres
}

// IdentifyResourcePresence creates an initialized ExistingResources struct
func IdentifyResourcePresence(id api.ID, handlers *clientCore.StateHandlers) (ExistingResources, error) {
	hz, err := handlers.Domain.GetPrimaryHostedZone()
	if err != nil && !errors.Is(err, storm.ErrNotFound) {
		return ExistingResources{}, err
	}

	dbs, err := handlers.Component.GetPostgresDatabases()
	if err != nil {
		return ExistingResources{}, nil
	}

	haveDBs := map[string]*v1alpha1.ClusterDatabasesPostgres{}

	for _, db := range dbs {
		haveDBs[db.ApplicationName] = &v1alpha1.ClusterDatabasesPostgres{
			Name:      db.ApplicationName,
			User:      db.UserName,
			Namespace: db.Namespace,
		}
	}

	return ExistingResources{
		hasServiceQuotaCheck:                  false,
		hasAWSLoadBalancerController:          !common.IsNotFound(handlers.Helm.GetHelmRelease(awslbc.ReleaseName)),
		hasCluster:                            !common.IsNotFound(handlers.Cluster.GetCluster(id.ClusterName)),
		hasExternalDNS:                        !common.IsNotFound(handlers.ExternalDNS.GetExternalDNS()),
		hasExternalSecrets:                    !common.IsNotFound(handlers.Helm.GetHelmRelease(externalsecrets.ReleaseName)),
		hasAutoscaler:                         !common.IsNotFound(handlers.Helm.GetHelmRelease(autoscaler.ReleaseName)),
		hasBlockstorage:                       !common.IsNotFound(handlers.Helm.GetHelmRelease(blockstorage.ReleaseName)),
		hasKubePromStack:                      !common.IsNotFound(handlers.Helm.GetHelmRelease(kubepromstack.ReleaseName)),
		hasLoki:                               !common.IsNotFound(handlers.Helm.GetHelmRelease(loki.ReleaseName)),
		hasPromtail:                           !common.IsNotFound(handlers.Helm.GetHelmRelease(promtail.ReleaseName)),
		hasTempo:                              !common.IsNotFound(handlers.Helm.GetHelmRelease(tempo.ReleaseName)),
		hasIdentityManager:                    !common.IsNotFound(handlers.IdentityManager.GetIdentityPool(cfn.NewStackNamer().IdentityPool(id.ClusterName))),
		hasArgoCD:                             !common.IsNotFound(handlers.ArgoCD.GetArgoCD()),
		hasPrimaryHostedZone:                  !common.IsNotFound(handlers.Domain.GetPrimaryHostedZone()),
		hasVPC:                                !common.IsNotFound(handlers.Vpc.GetVpc(cfn.NewStackNamer().Vpc(id.ClusterName))),
		hasDelegatedHostedZoneNameservers:     hz != nil && hz.IsDelegated,
		hasDelegatedHostedZoneNameserversTest: false,
		hasUsers:                              false, // For now we will always check if there are missing users
		hasPostgres:                           haveDBs,
	}, nil
}
