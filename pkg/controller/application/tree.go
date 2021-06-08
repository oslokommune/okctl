package application

import "github.com/oslokommune/okctl/pkg/controller/common/resourcetree"

// CreateResourceDependencyTree creates a dependency tree for applications
func CreateResourceDependencyTree() *resourcetree.ResourceNode {
	root := resourcetree.NewNode(resourcetree.ResourceNodeTypeGroup)

	containerRepositoryNode := resourcetree.NewNode(resourcetree.ResourceNodeTypeContainerRepository)
	containerRepositoryNode.AppendChild(resourcetree.NewNode(resourcetree.ResourceNodeTypeApplication))

	root.AppendChild(containerRepositoryNode)

	return root
}
