package application

import "github.com/oslokommune/okctl/pkg/controller/common/dependencytree"

// CreateResourceDependencyTree creates a dependency tree for applications
func CreateResourceDependencyTree() *dependencytree.Node {
	root := dependencytree.NewNode(dependencytree.NodeTypeGroup)

	containerRepositoryNode := dependencytree.NewNode(dependencytree.NodeTypeContainerRepository)
	containerRepositoryNode.AppendChild(dependencytree.NewNode(dependencytree.NodeTypeApplication))

	root.AppendChild(containerRepositoryNode)

	return root
}
