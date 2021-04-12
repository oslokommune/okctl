package controller

import (
	"time"

	"github.com/oslokommune/okctl/pkg/helm/charts/externalsecrets"

	"github.com/oslokommune/okctl/pkg/helm/charts/autoscaler"
	"github.com/oslokommune/okctl/pkg/helm/charts/awslbc"
	"github.com/oslokommune/okctl/pkg/helm/charts/blockstorage"
	"github.com/oslokommune/okctl/pkg/helm/charts/kubepromstack"
	lokipkg "github.com/oslokommune/okctl/pkg/helm/charts/loki"
	"github.com/oslokommune/okctl/pkg/helm/charts/promtail"
	"github.com/oslokommune/okctl/pkg/helm/charts/tempo"

	"github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/pkg/errors"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ExistingResources contains information about what services already exists in a cluster
type ExistingResources struct {
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
	hasPostgres                           bool
}

func isNotFound(_ interface{}, err error) bool {
	return errors.Is(err, storm.ErrNotFound)
}

// IdentifyResourcePresence creates an initialized ExistingResources struct
func IdentifyResourcePresence(id api.ID, handlers *clientCore.StateHandlers) (ExistingResources, error) {
	hz, err := handlers.Domain.GetPrimaryHostedZone()
	if err != nil && !errors.Is(err, storm.ErrNotFound) {
		return ExistingResources{}, err
	}

	return ExistingResources{
		hasPrimaryHostedZone:                  !isNotFound(handlers.Domain.GetPrimaryHostedZone()),
		hasVPC:                                !isNotFound(handlers.Vpc.GetVpc(cfn.NewStackNamer().Vpc(id.ClusterName))),
		hasCluster:                            !isNotFound(handlers.Cluster.GetCluster(id.ClusterName)),
		hasExternalSecrets:                    !isNotFound(handlers.Helm.GetHelmRelease(externalsecrets.ExternalSecrets(nil).ReleaseName)),
		hasAutoscaler:                         !isNotFound(handlers.Helm.GetHelmRelease(autoscaler.New(nil).ReleaseName)),
		hasKubePromStack:                      !isNotFound(handlers.Helm.GetHelmRelease(kubepromstack.New(0*time.Second, nil).ReleaseName)),
		hasLoki:                               !isNotFound(handlers.Helm.GetHelmRelease(lokipkg.New(nil).ReleaseName)),
		hasPromtail:                           !isNotFound(handlers.Helm.GetHelmRelease(promtail.New(nil).ReleaseName)),
		hasTempo:                              !isNotFound(handlers.Helm.GetHelmRelease(tempo.New(nil).ReleaseName)),
		hasBlockstorage:                       !isNotFound(handlers.Helm.GetHelmRelease(blockstorage.New(nil).ReleaseName)),
		hasAWSLoadBalancerController:          !isNotFound(handlers.Helm.GetHelmRelease(awslbc.New(nil).ReleaseName)),
		hasExternalDNS:                        !isNotFound(handlers.ExternalDNS.GetExternalDNS()),
		hasIdentityManager:                    !isNotFound(handlers.IdentityManager.GetIdentityPool(cfn.NewStackNamer().IdentityPool(id.ClusterName))),
		hasArgoCD:                             !isNotFound(handlers.ArgoCD.GetArgoCD()),
		hasDelegatedHostedZoneNameservers:     hz != nil && hz.IsDelegated,
		hasDelegatedHostedZoneNameserversTest: false,
		hasUsers:                              false, // For now we will always check if there are missing users
		hasPostgres:                           false, // For now we will always check if there are missing postgres databases
	}, nil
}

// CreateResourceDependencyTree creates a tree
func CreateResourceDependencyTree() (root *resourcetree.ResourceNode) {
	root = createNode(nil, resourcetree.ResourceNodeTypeGroup)

	var vpcNode,
		clusterNode,
		primaryHostedZoneNode *resourcetree.ResourceNode

	primaryHostedZoneNode = createNode(root, resourcetree.ResourceNodeTypeZone)
	createNode(primaryHostedZoneNode, resourcetree.ResourceNodeTypeNameserverDelegator)

	vpcNode = createNode(primaryHostedZoneNode, resourcetree.ResourceNodeTypeVPC)

	clusterNode = createNode(vpcNode, resourcetree.ResourceNodeTypeCluster)
	createNode(clusterNode, resourcetree.ResourceNodeTypeExternalSecrets)
	createNode(clusterNode, resourcetree.ResourceNodeTypeAutoscaler)
	createNode(clusterNode, resourcetree.ResourceNodeTypeBlockstorage)
	createNode(clusterNode, resourcetree.ResourceNodeTypeAWSLoadBalancerController)
	createNode(clusterNode, resourcetree.ResourceNodeTypeExternalDNS)
	createNode(clusterNode, resourcetree.ResourceNodeTypePostgres)

	// All resources that requires SSL / a certificate needs the delegatedNameserversConfirmedNode as a dependency
	delegatedNameserversConfirmedNode := createNode(clusterNode, resourcetree.ResourceNodeTypeNameserversDelegatedTest)

	identityProviderNode := createNode(delegatedNameserversConfirmedNode, resourcetree.ResourceNodeTypeIdentityManager)
	createNode(identityProviderNode, resourcetree.ResourceNodeTypeArgoCD)
	createNode(identityProviderNode, resourcetree.ResourceNodeTypeUsers)
	kubePromStack := createNode(identityProviderNode, resourcetree.ResourceNodeTypeKubePromStack)
	// This is not strictly required, but to a large extent it doesn't make much sense to setup Loki before
	// we have setup grafana.
	loki := createNode(kubePromStack, resourcetree.ResourceNodeTypeLoki)
	createNode(kubePromStack, resourcetree.ResourceNodeTypeTempo)
	// Similarly, it doesn't make sense to install promtail without loki
	createNode(loki, resourcetree.ResourceNodeTypePromtail)

	return root
}

// CreateApplicationResourceDependencyTree creates a dependency tree for applications
func CreateApplicationResourceDependencyTree() (root *resourcetree.ResourceNode) {
	root = createNode(nil, resourcetree.ResourceNodeTypeGroup)

	createNode(root, resourcetree.ResourceNodeTypeApplication)

	return root
}

func createNode(parent *resourcetree.ResourceNode, nodeType resourcetree.ResourceNodeType) (child *resourcetree.ResourceNode) {
	child = &resourcetree.ResourceNode{
		Type:     nodeType,
		Children: make([]*resourcetree.ResourceNode, 0),
	}

	child.State = resourcetree.ResourceNodeStatePresent

	if parent != nil {
		parent.Children = append(parent.Children, child)
	}

	return child
}
