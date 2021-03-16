package controller

import (
	"path"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/spf13/afero"
)

// ExistingResources contains information about what services already exists in a cluster
type ExistingResources struct {
	hasALBIngressController               bool
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
}

// IdentifyResourcePresence creates an initialized ExistingResources struct
func IdentifyResourcePresence(fs *afero.Afero, outputDir string, hzFetcher HostedZoneFetcher) (ExistingResources, error) {
	hz := hzFetcher()

	return ExistingResources{
		hasPrimaryHostedZone:                  hz != nil,
		hasVPC:                                directoryTester(fs, outputDir, constant.DefaultVpcBaseDir),
		hasCluster:                            directoryTester(fs, outputDir, constant.DefaultClusterBaseDir),
		hasExternalSecrets:                    directoryTester(fs, outputDir, constant.DefaultExternalSecretsBaseDir),
		hasAutoscaler:                         directoryTester(fs, outputDir, constant.DefaultAutoscalerBaseDir),
		hasKubePromStack:                      directoryTester(fs, outputDir, path.Join(constant.DefaultMonitoringBaseDir, constant.DefaultKubePromStackBaseDir)),
		hasLoki:                               directoryTester(fs, outputDir, path.Join(constant.DefaultMonitoringBaseDir, constant.DefaultLokiBaseDir)),
		hasPromtail:                           directoryTester(fs, outputDir, path.Join(constant.DefaultMonitoringBaseDir, constant.DefaultPromtailBaseDir)),
		hasTempo:                              directoryTester(fs, outputDir, path.Join(constant.DefaultMonitoringBaseDir, constant.DefaultTempoBaseDir)),
		hasBlockstorage:                       directoryTester(fs, outputDir, constant.DefaultBlockstorageBaseDir),
		hasALBIngressController:               directoryTester(fs, outputDir, constant.DefaultAlbIngressControllerBaseDir),
		hasAWSLoadBalancerController:          directoryTester(fs, outputDir, constant.DefaultAWSLoadBalancerControllerBaseDir),
		hasExternalDNS:                        directoryTester(fs, outputDir, constant.DefaultExternalDNSBaseDir),
		hasIdentityManager:                    directoryTester(fs, outputDir, constant.DefaultIdentityPoolBaseDir),
		hasArgoCD:                             directoryTester(fs, outputDir, constant.DefaultArgoCDBaseDir),
		hasDelegatedHostedZoneNameservers:     hz != nil && hz.IsDelegated,
		hasDelegatedHostedZoneNameserversTest: false,
		hasUsers:                              false, // this means we will always check? needs to be changed?
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
	createNode(clusterNode, resourcetree.ResourceNodeTypeALBIngress)
	createNode(clusterNode, resourcetree.ResourceNodeTypeAWSLoadBalancerController)
	createNode(clusterNode, resourcetree.ResourceNodeTypeExternalDNS)

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

func directoryTester(fs *afero.Afero, outputDir string, target string) bool {
	baseDir := path.Join(outputDir, target)

	exists, _ := fs.DirExists(baseDir)

	return exists
}

// GithubGetter knows how to get the current state Github
type GithubGetter func() state.Github
