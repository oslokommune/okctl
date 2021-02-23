package controller

import (
	"path"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/spf13/afero"
)

// ExistingResources contains information about what services already exists in a cluster
type ExistingResources struct {
	hasALBIngressController           bool
	hasAWSLoadBalancerController      bool
	hasCluster                        bool
	hasExternalDNS                    bool
	hasExternalSecrets                bool
	hasAutoscaler                     bool
	hasBlockstorage                   bool
	hasKubePromStack                  bool
	hasIdentityManager                bool
	hasArgoCD                         bool
	hasPrimaryHostedZone              bool
	hasVPC                            bool
	hasDelegatedHostedZoneNameservers bool
}

// IdentifyResourcePresence creates an initialized ExistingResources struct
func IdentifyResourcePresence(fs *afero.Afero, outputDir string, hzFetcher HostedZoneFetcher) (ExistingResources, error) {
	hz := hzFetcher()

	return ExistingResources{
		hasPrimaryHostedZone:              hz != nil,
		hasVPC:                            directoryTester(fs, outputDir, config.DefaultVpcBaseDir),
		hasCluster:                        directoryTester(fs, outputDir, config.DefaultClusterBaseDir),
		hasExternalSecrets:                directoryTester(fs, outputDir, config.DefaultExternalSecretsBaseDir),
		hasAutoscaler:                     directoryTester(fs, outputDir, config.DefaultAutoscalerBaseDir),
		hasKubePromStack:                  directoryTester(fs, outputDir, config.DefaultKubePromStackBaseDir),
		hasBlockstorage:                   directoryTester(fs, outputDir, config.DefaultBlockstorageBaseDir),
		hasALBIngressController:           directoryTester(fs, outputDir, config.DefaultAlbIngressControllerBaseDir),
		hasAWSLoadBalancerController:      directoryTester(fs, outputDir, config.DefaultAWSLoadBalancerControllerBaseDir),
		hasExternalDNS:                    directoryTester(fs, outputDir, config.DefaultExternalDNSBaseDir),
		hasIdentityManager:                directoryTester(fs, outputDir, config.DefaultIdentityPoolBaseDir),
		hasArgoCD:                         directoryTester(fs, outputDir, config.DefaultArgoCDBaseDir),
		hasDelegatedHostedZoneNameservers: hz != nil && hz.IsDelegated,
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

	identityProviderNode := createNode(clusterNode, resourcetree.ResourceNodeTypeIdentityManager)
	createNode(identityProviderNode, resourcetree.ResourceNodeTypeKubePromStack)
	createNode(identityProviderNode, resourcetree.ResourceNodeTypeArgoCD)

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
