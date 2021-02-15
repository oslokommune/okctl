package controller

import (
	"fmt"
	"path"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/oslokommune/okctl/pkg/git"
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
	hasGithubSetup                    bool
	hasIdentityManager                bool
	hasPrimaryHostedZone              bool
	hasVPC                            bool
	hasDelegatedHostedZoneNameservers bool
}

// NewCreateCurrentStateTreeOpts creates an initialized ExistingServices struct
func NewCreateCurrentStateTreeOpts(fs *afero.Afero, outputDir string, githubGetter reconciler.GithubGetter, hzFetcher HostedZoneFetcher) (*ExistingServices, error) {
	hz := hzFetcher()

	return &ExistingServices{
		hasGithubSetup:                    githubTester(githubGetter()),
		hasPrimaryHostedZone:              hz != nil,
		hasVPC:                            directoryTester(fs, outputDir, config.DefaultVpcBaseDir),
		hasCluster:                        directoryTester(fs, outputDir, config.DefaultClusterBaseDir),
		hasExternalSecrets:                directoryTester(fs, outputDir, config.DefaultExternalSecretsBaseDir),
		hasAutoscaler:                     directoryTester(fs, outputDir, config.DefaultAutoscalerBaseDir),
		hasALBIngressController:           directoryTester(fs, outputDir, config.DefaultAlbIngressControllerBaseDir),
		hasAWSLoadBalancerController:      directoryTester(fs, outputDir, config.DefaultAWSLoadBalancerControllerBaseDir),
		hasExternalDNS:                    directoryTester(fs, outputDir, config.DefaultExternalDNSBaseDir),
		hasIdentityManager:                directoryTester(fs, outputDir, config.DefaultIdentityPoolBaseDir),
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
	createNode(clusterNode, resourcetree.ResourceNodeTypeALBIngress, opts.hasALBIngressController)
	createNode(clusterNode, resourcetree.ResourceNodeTypeAWSLoadBalancerController, opts.hasAWSLoadBalancerController)
	createNode(clusterNode, resourcetree.ResourceNodeTypeExternalDNS, opts.hasExternalDNS)

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
	createNode(clusterNode, resourcetree.ResourceNodeTypeALBIngress, cluster.Integrations.ALBIngressController)
	createNode(clusterNode, resourcetree.ResourceNodeTypeAWSLoadBalancerController, cluster.Integrations.AWSLoadBalancerController)
	createNode(clusterNode, resourcetree.ResourceNodeTypeExternalDNS, cluster.Integrations.ExternalDNS)

	return root
}

// ApplyDesiredStateMetadata applies metadata from a cluster definition to the nodes
func ApplyDesiredStateMetadata(tree *resourcetree.ResourceNode, cluster *v1alpha1.Cluster, repoDir string) error {
	primaryHostedZoneNode := tree.GetNode(&resourcetree.ResourceNode{Type: resourcetree.ResourceNodeTypeZone})
	if primaryHostedZoneNode == nil {
		return errors.New("expected primary hosted zone node was not found")
	}

	primaryHostedZoneNode.Metadata = reconciler.HostedZoneMetadata{Domain: cluster.PrimaryDNSZone.ParentDomain}

	vpcNode := tree.GetNode(&resourcetree.ResourceNode{Type: resourcetree.ResourceNodeTypeVPC})
	if vpcNode == nil {
		return errors.New("expected vpc node was not found")
	}

	vpcNode.Metadata = reconciler.VPCMetadata{
		Cidr:             cluster.VPC.CIDR,
		HighAvailability: cluster.VPC.HighAvailability,
	}

	githubNode := tree.GetNode(&resourcetree.ResourceNode{Type: resourcetree.ResourceNodeTypeGithub})
	if githubNode == nil {
		return errors.New("expected github node was not found")
	}

	repo, err := git.GithubRepoFullName(cluster.Github.Organisation, repoDir)
	if err != nil {
		return fmt.Errorf("error fetching full git repo name: %w", err)
	}

	githubNode.Metadata = reconciler.GithubMetadata{
		Organization: cluster.Github.Organisation,
		Repository:   repo,
	}

	argocdNode := tree.GetNode(&resourcetree.ResourceNode{Type: resourcetree.ResourceNodeTypeArgoCD})
	if argocdNode != nil {
		argocdNode.Metadata = reconciler.ArgocdMetadata{Organization: cluster.Github.Organisation}
	}

	return nil
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
