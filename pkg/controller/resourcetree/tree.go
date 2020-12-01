package resourcetree

import (
	"context"
	"github.com/oslokommune/okctl/pkg/api"
)

// ResourceNodeType defines what type of resource a ResourceNode represents
type ResourceNodeType int
const (
	// ResourceNodeTypeGroup represents a node that has no actions associated with it. For now, only the root node
	ResourceNodeTypeGroup ResourceNodeType = iota
	// ResourceNodeTypeZone represents a HostedZone resource
	ResourceNodeTypeZone
	// ResourceNodeTypeVPC represents a VPC resource
	ResourceNodeTypeVPC
	// ResourceNodeTypeCluster represents a EKS cluster resource
	ResourceNodeTypeCluster
	// ResourceNodeTypeExternalSecrets represents an External Secrets resource
	ResourceNodeTypeExternalSecrets
	// ResourceNodeTypeALBIngress represents an ALB Ingress resource
	ResourceNodeTypeALBIngress
	// ResourceNodeTypeExternalDNS represents an External DNS resource
	ResourceNodeTypeExternalDNS
	// ResourceNodeTypeGithub represents a Github setup
	ResourceNodeTypeGithub
	// ResourceNodeTypeIdentityManager represents a Identity Manager resource
	ResourceNodeTypeIdentityManager
	// ResourceNodeTypeArgoCD represents an ArgoCD resource
	ResourceNodeTypeArgoCD
)

// ResourceNodeState defines what state the resource is in, used to infer what action to take
type ResourceNodeState int
const (
	// ResourceNodeStateNoop represents a state where no action is needed. E.g.: if the desired state of the
	// resource conforms with the actual state
	ResourceNodeStateNoop ResourceNodeState = iota
	// ResourceNodeStatePresent represents the state where the resource exists
	ResourceNodeStatePresent
	// ResourceNodeStateAbsent represents the state where the resource does not exist
	ResourceNodeStateAbsent
)

// CommonMetadata represents metadata required by most if not all operations on services
type CommonMetadata struct {
	Ctx       context.Context
	ClusterId api.ID
}

// StateRefreshFn is a function that attempts to retrieve state potentially can only be retrieved at runtime. E.g.:
// state that can only exist after an external resource has been created
type StateRefreshFn func(node *ResourceNode)

// ResourceNode represents a component of the cluster and its dependencies
type ResourceNode struct {
	Type  ResourceNodeType
	State ResourceNodeState

	// Contains metadata regarding the resource supplied by the desired state definition
	Metadata             interface{}

	StateRefresher 		 StateRefreshFn
	// ResourceState contains data that needs to be retrieved runtime. In other words, data that potentially only exist
	// after an external resource has been created
	ResourceState 		 interface{}

	Children []*ResourceNode
}

// RefreshState calls the stored StateRefreshFn if it exists
func (node *ResourceNode) RefreshState() {
	if node.StateRefresher == nil {
		return
	}

	node.StateRefresher(node)
}

// SetStateRefresher stores a StateRefreshFn on the node to be used to retrieve runtime state later
func (node *ResourceNode) SetStateRefresher(nodeType ResourceNodeType, refresher StateRefreshFn) {
	targetNode := node.GetNode(&ResourceNode{Type: nodeType})

	if targetNode == nil {
		return
	}

	targetNode.StateRefresher = refresher
}

// Equals knows how to compare two ResourceNodes and determine equality
func (node *ResourceNode) Equals(targetNode *ResourceNode) bool {
	if targetNode == nil {
		return false
	}

	// P.S.: For now, this is good enough due to all types of resources only existing once
	return node.Type == targetNode.Type
}

// GetNode returns an identical node as targetNode from the receiver's tree
func (node *ResourceNode) GetNode(targetNode *ResourceNode) *ResourceNode {
	if node.Equals(targetNode) {
		return node
	}
	
	for _, child := range node.Children {
		result := child.GetNode(targetNode)

		if result != nil {
			return result
		}
	}
	
	return nil
}

// ApplyFn is a kind of function we can run on all the nodes in a ResourceNode tree with the ApplyFunction() function
type ApplyFn func(receiver *ResourceNode, target *ResourceNode)

// ApplyFunction will use the supplied ApplyFn on all the nodes in the receiver tree, with an equal node from the target
// tree
func (node *ResourceNode) ApplyFunction(fn ApplyFn, targetGraph *ResourceNode) {
	for _, child := range node.Children {
		child.ApplyFunction(fn, targetGraph)
	}
	
	targetNode := targetGraph.GetNode(node)
	fn(node, targetNode)
}
