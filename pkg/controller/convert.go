package controller

import (
	"path"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/spf13/afero"
)

// ExistingServices contains information about what services already exists in a cluster
type ExistingServices struct {
	hasALBIngressController           bool
	hasAWSLoadBalancerController      bool
	hasCluster                        bool
	hasExternalDNS                    bool
	hasExternalSecrets                bool
	hasAutoscaler                     bool
	hasBlockstorage                   bool
	hasKubePromStack                  bool
	hasGithubSetup                    bool
	hasIdentityManager                bool
	hasArgoCD                         bool
	hasPrimaryHostedZone              bool
	hasVPC                            bool
	hasDelegatedHostedZoneNameservers bool
}

// NewCreateCurrentStateTreeOpts creates an initialized ExistingServices struct
func NewCreateCurrentStateTreeOpts(fs *afero.Afero, outputDir string, githubGetter GithubGetter, hzFetcher HostedZoneFetcher) (*ExistingServices, error) {
	hz := hzFetcher()

	return &ExistingServices{
		hasGithubSetup:                    githubTester(githubGetter()),
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

// CreateCurrentStateTree knows how to generate a ResourceNode tree based on the current state
func CreateCurrentStateTree(opts *ExistingServices) (root *resourcetree.ResourceNode) {
	root = createNode(nil, resourcetree.ResourceNodeTypeGroup, true)

	var vpcNode,
		clusterNode *resourcetree.ResourceNode

	createNode(root, resourcetree.ResourceNodeTypeGithub, opts.hasGithubSetup)

	primaryHostedZoneNode := createNode(root, resourcetree.ResourceNodeTypeZone, opts.hasPrimaryHostedZone)
	createNode(primaryHostedZoneNode,
		resourcetree.ResourceNodeTypeNameserverDelegator,
		opts.hasDelegatedHostedZoneNameservers,
	)

	vpcNode = createNode(primaryHostedZoneNode, resourcetree.ResourceNodeTypeVPC, opts.hasVPC)

	clusterNode = createNode(vpcNode, resourcetree.ResourceNodeTypeCluster, opts.hasCluster)

	createNode(clusterNode, resourcetree.ResourceNodeTypeExternalSecrets, opts.hasExternalSecrets)
	createNode(clusterNode, resourcetree.ResourceNodeTypeAutoscaler, opts.hasAutoscaler)
	createNode(clusterNode, resourcetree.ResourceNodeTypeKubePromStack, opts.hasKubePromStack)
	createNode(clusterNode, resourcetree.ResourceNodeTypeBlockstorage, opts.hasBlockstorage)
	createNode(clusterNode, resourcetree.ResourceNodeTypeALBIngress, opts.hasALBIngressController)
	createNode(clusterNode, resourcetree.ResourceNodeTypeAWSLoadBalancerController, opts.hasAWSLoadBalancerController)
	createNode(clusterNode, resourcetree.ResourceNodeTypeExternalDNS, opts.hasExternalDNS)

	identityProviderNode := createNode(clusterNode, resourcetree.ResourceNodeTypeIdentityManager, opts.hasIdentityManager)
	createNode(identityProviderNode, resourcetree.ResourceNodeTypeArgoCD, opts.hasArgoCD)

	return root
}

// CreateDesiredStateTree knows how to create a ResourceNode tree based on a cluster declaration
func CreateDesiredStateTree(cluster *v1alpha1.Cluster) (root *resourcetree.ResourceNode) {
	root = createNode(nil, resourcetree.ResourceNodeTypeGroup, true)

	var vpcNode,
		clusterNode *resourcetree.ResourceNode

	createNode(root, resourcetree.ResourceNodeTypeGithub, true)

	primaryHostedZoneNode := createNode(root, resourcetree.ResourceNodeTypeZone, true)
	createNode(primaryHostedZoneNode, resourcetree.ResourceNodeTypeNameserverDelegator, true)

	vpcNode = createNode(primaryHostedZoneNode, resourcetree.ResourceNodeTypeVPC, true)

	clusterNode = createNode(vpcNode, resourcetree.ResourceNodeTypeCluster, true)

	createNode(clusterNode, resourcetree.ResourceNodeTypeExternalSecrets, cluster.Integrations.ExternalSecrets)
	createNode(clusterNode, resourcetree.ResourceNodeTypeAutoscaler, cluster.Integrations.Autoscaler)
	createNode(clusterNode, resourcetree.ResourceNodeTypeKubePromStack, cluster.Integrations.KubePromStack)
	createNode(clusterNode, resourcetree.ResourceNodeTypeBlockstorage, cluster.Integrations.Blockstorage)
	createNode(clusterNode, resourcetree.ResourceNodeTypeALBIngress, cluster.Integrations.ALBIngressController)
	createNode(clusterNode, resourcetree.ResourceNodeTypeAWSLoadBalancerController, cluster.Integrations.AWSLoadBalancerController)
	createNode(clusterNode, resourcetree.ResourceNodeTypeExternalDNS, cluster.Integrations.ExternalDNS)

	identityProviderNode := createNode(clusterNode, resourcetree.ResourceNodeTypeIdentityManager, cluster.Integrations.Cognito)
	createNode(identityProviderNode, resourcetree.ResourceNodeTypeArgoCD, cluster.Integrations.ArgoCD)

	return root
}

func createNode(parent *resourcetree.ResourceNode, nodeType resourcetree.ResourceNodeType, present bool) (child *resourcetree.ResourceNode) {
	child = &resourcetree.ResourceNode{
		Type:     nodeType,
		Children: make([]*resourcetree.ResourceNode, 0),
	}

	if present {
		child.State = resourcetree.ResourceNodeStatePresent
	} else {
		child.State = resourcetree.ResourceNodeStateAbsent
	}

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

func githubTester(github state.Github) bool {
	if len(github.Repositories) == 0 {
		return false
	}

	for _, repo := range github.Repositories {
		err := repo.Validate()
		if err != nil {
			return false
		}

		break
	}

	return true
}

// GithubGetter knows how to get the current state Github
type GithubGetter func() state.Github
