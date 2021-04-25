// Package resourcetree contains Resource Node specification and common operations done on them
package resourcetree

import (
	"bytes"
	"context"
	"io"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

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
	// ResourceNodeTypeAutoscaler represents an autoscaler resource
	ResourceNodeTypeAutoscaler
	// ResourceNodeTypeBlockstorage represents a blockstorage resource
	ResourceNodeTypeBlockstorage
	// ResourceNodeTypeKubePromStack represents a kubernetes-prometheus-stack resource
	ResourceNodeTypeKubePromStack
	// ResourceNodeTypeLoki represents a loki resource
	ResourceNodeTypeLoki
	// ResourceNodeTypePromtail represents a promtail deployment
	ResourceNodeTypePromtail
	// ResourceNodeTypeTempo represents a Tempo deployment
	ResourceNodeTypeTempo
	// ResourceNodeTypeAWSLoadBalancerController represents an AWS load balancer controller resource
	ResourceNodeTypeAWSLoadBalancerController
	// ResourceNodeTypeExternalDNS represents an External DNS resource
	ResourceNodeTypeExternalDNS
	// ResourceNodeTypeIdentityManager represents a Identity Manager resource
	ResourceNodeTypeIdentityManager
	// ResourceNodeTypeArgoCD represents an ArgoCD resource
	ResourceNodeTypeArgoCD
	// ResourceNodeTypeNameserverDelegator represents delegation of nameservers for a HostedZone
	ResourceNodeTypeNameserverDelegator
	// ResourceNodeTypeNameserversDelegatedTest represents testing if nameservers has been successfully delegated
	ResourceNodeTypeNameserversDelegatedTest
	// ResourceNodeTypeUsers represents the users we want to add to the cognito user pool
	ResourceNodeTypeUsers
	// ResourceNodeTypePostgres represents the postgres databases we want to add to the cluster
	ResourceNodeTypePostgres
	// ResourceNodeTypePostgresInstance represents a postgres instance
	ResourceNodeTypePostgresInstance
	// ResourceNodeTypeApplication represents an okctl application resource
	ResourceNodeTypeApplication
	// ResourceNodeTypeCleanupALB represents a cleanup of ALBs
	ResourceNodeTypeCleanupALB
	// ResourceNodeTypeCleanupSG represents a cleanup of SecurityGroups
	ResourceNodeTypeCleanupSG
	// ResourceNodeTypeServiceQuota represents a service quota check
	ResourceNodeTypeServiceQuota
)

// ResourceNodeTypeToString knows how to convert a Resource Node type to a human readable string
// nolint: gocyclo funlen
func ResourceNodeTypeToString(nodeType ResourceNodeType) string {
	switch nodeType {
	case ResourceNodeTypeGroup:
		return "Group"
	case ResourceNodeTypeZone:
		return "Hosted Zone"
	case ResourceNodeTypeVPC:
		return "VPC"
	case ResourceNodeTypeCluster:
		return "K8s Cluster"
	case ResourceNodeTypeExternalSecrets:
		return "External Secrets"
	case ResourceNodeTypeAutoscaler:
		return "Autoscaler"
	case ResourceNodeTypeKubePromStack:
		return "KubePromStack"
	case ResourceNodeTypeLoki:
		return "Loki"
	case ResourceNodeTypePromtail:
		return "Promtail"
	case ResourceNodeTypeTempo:
		return "Tempo"
	case ResourceNodeTypeBlockstorage:
		return "Blockstorage"
	case ResourceNodeTypeAWSLoadBalancerController:
		return "AWS Load Balancer Controller"
	case ResourceNodeTypeExternalDNS:
		return "External DNS"
	case ResourceNodeTypeIdentityManager:
		return "Identity Manager"
	case ResourceNodeTypeArgoCD:
		return "ArgoCD controller"
	case ResourceNodeTypeNameserverDelegator:
		return "Nameserver Delegation"
	case ResourceNodeTypeNameserversDelegatedTest:
		return "Nameservers Delegated Test"
	case ResourceNodeTypeUsers:
		return "Users"
	case ResourceNodeTypePostgres:
		return "Postgres"
	case ResourceNodeTypePostgresInstance:
		return "Postgres Instance"
	case ResourceNodeTypeCleanupALB:
		return "Cleanup ALBs"
	case ResourceNodeTypeCleanupSG:
		return "Cleanup SGs"
	case ResourceNodeTypeServiceQuota:
		return "Service Quota"
	default:
		return "N/A"
	}
}

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
	Ctx context.Context

	Out io.Writer

	ClusterID              api.ID
	Declaration            *v1alpha1.Cluster
	ApplicationDeclaration v1alpha1.Application
}

// ResourceNode represents a component of the cluster and its dependencies
type ResourceNode struct {
	Type     ResourceNodeType
	State    ResourceNodeState
	Data     interface{}
	Children []*ResourceNode
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

func (node *ResourceNode) treeViewGenerator(writer io.Writer, tabs int) {
	result := ""
	for index := 0; index < tabs; index++ {
		result += "\t"
	}

	result += "- " + ResourceNodeTypeToString(node.Type) + "\n"

	if node.State == ResourceNodeStatePresent {
		_, _ = writer.Write([]byte(result))
	}

	for _, child := range node.Children {
		child.treeViewGenerator(writer, tabs+1)
	}
}

func (node *ResourceNode) String() string {
	var buf bytes.Buffer

	node.treeViewGenerator(&buf, 0)

	return buf.String()
}

// ApplyFn is a kind of function we can run on all the nodes in a ResourceNode tree with the ApplyFunction() function
type ApplyFn func(receiver *ResourceNode, target *ResourceNode)

// ApplyFunction will use the supplied ApplyFn on all the nodes in the receiver tree, with an equal node from the target
// tree
func (node *ResourceNode) ApplyFunction(fn ApplyFn, targetTree *ResourceNode) {
	for _, child := range node.Children {
		child.ApplyFunction(fn, targetTree)
	}

	targetNode := targetTree.GetNode(node)
	fn(node, targetNode)
}
