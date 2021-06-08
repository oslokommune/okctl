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
type ResourceNodeType string

const (
	// ResourceNodeTypeGroup represents a node that has no actions associated with it. For now, only the root node
	ResourceNodeTypeGroup ResourceNodeType = "group"
	// ResourceNodeTypeZone represents a HostedZone resource
	ResourceNodeTypeZone ResourceNodeType = "hosted-zone"
	// ResourceNodeTypeVPC represents a VPC resource
	ResourceNodeTypeVPC ResourceNodeType = "vpc"
	// ResourceNodeTypeCluster represents a EKS cluster resource
	ResourceNodeTypeCluster ResourceNodeType = "cluster"
	// ResourceNodeTypeExternalSecrets represents an External Secrets resource
	ResourceNodeTypeExternalSecrets ResourceNodeType = "external-secrets"
	// ResourceNodeTypeAutoscaler represents an autoscaler resource
	ResourceNodeTypeAutoscaler ResourceNodeType = "autoscaler"
	// ResourceNodeTypeBlockstorage represents a blockstorage resource
	ResourceNodeTypeBlockstorage ResourceNodeType = "blockstorage"
	// ResourceNodeTypeKubePromStack represents a kubernetes-prometheus-stack resource
	ResourceNodeTypeKubePromStack ResourceNodeType = "kubernetes-prometheus-stack"
	// ResourceNodeTypeLoki represents a loki resource
	ResourceNodeTypeLoki ResourceNodeType = "loki"
	// ResourceNodeTypePromtail represents a promtail deployment
	ResourceNodeTypePromtail ResourceNodeType = "promtail"
	// ResourceNodeTypeTempo represents a Tempo deployment
	ResourceNodeTypeTempo ResourceNodeType = "tempo"
	// ResourceNodeTypeAWSLoadBalancerController represents an AWS load balancer controller resource
	ResourceNodeTypeAWSLoadBalancerController ResourceNodeType = "aws-load-balancer-controller"
	// ResourceNodeTypeExternalDNS represents an External DNS resource
	ResourceNodeTypeExternalDNS ResourceNodeType = "external-dns"
	// ResourceNodeTypeIdentityManager represents a Identity Manager resource
	ResourceNodeTypeIdentityManager ResourceNodeType = "identity-manager"
	// ResourceNodeTypeArgoCD represents an ArgoCD resource
	ResourceNodeTypeArgoCD ResourceNodeType = "argocd"
	// ResourceNodeTypeNameserverDelegator represents delegation of nameservers for a HostedZone
	ResourceNodeTypeNameserverDelegator ResourceNodeType = "nameserver-delegator"
	// ResourceNodeTypeNameserversDelegatedTest represents testing if nameservers has been successfully delegated
	ResourceNodeTypeNameserversDelegatedTest ResourceNodeType = "nameserver-delegator-test"
	// ResourceNodeTypeUsers represents the users we want to add to the cognito user pool
	ResourceNodeTypeUsers ResourceNodeType = "users"
	// ResourceNodeTypePostgres represents the postgres databases we want to add to the cluster
	ResourceNodeTypePostgres ResourceNodeType = "postgres"
	// ResourceNodeTypePostgresInstance represents a postgres instance
	ResourceNodeTypePostgresInstance ResourceNodeType = "postgres-instance"
	// ResourceNodeTypeApplication represents an okctl application resource
	ResourceNodeTypeApplication ResourceNodeType = "application"
	// ResourceNodeTypeCleanupALB represents a cleanup of ALBs
	ResourceNodeTypeCleanupALB ResourceNodeType = "cleanup-alb"
	// ResourceNodeTypeCleanupSG represents a cleanup of SecurityGroups
	ResourceNodeTypeCleanupSG ResourceNodeType = "cleanup-sg"
	// ResourceNodeTypeServiceQuota represents a service quota check
	ResourceNodeTypeServiceQuota ResourceNodeType = "service-quota"
	// ResourceNodeTypeContainerRepository represents a container repository
	ResourceNodeTypeContainerRepository ResourceNodeType = "container-repository"
)

// String representation of the ResourceNodeType
func (t ResourceNodeType) String() string {
	return string(t)
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

	result += "- " + node.Type.String()

	if node.State == ResourceNodeStatePresent {
		result += " (add)" + "\n"
		_, _ = writer.Write([]byte(result))
	}

	if node.State == ResourceNodeStateAbsent {
		result += " (remove)" + "\n"
		_, _ = writer.Write([]byte(result))
	}

	if node.State == ResourceNodeStateNoop {
		result += " (noop)" + "\n"
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
type ApplyFn func(receiver *ResourceNode)

// ApplyFnWithTarget is a kind of function we can run on all the nodes in a ResourceNode tree with the ApplyFunction() function
type ApplyFnWithTarget func(receiver *ResourceNode, target *ResourceNode)

// ApplyFunction will use the supplied ApplyFn on all the nodes in the receiver tree
func (node *ResourceNode) ApplyFunction(fn ApplyFn) {
	for _, child := range node.Children {
		child.ApplyFunction(fn)
	}

	fn(node)
}

// ApplyFunctionWithTarget will use the supplied ApplyFnWithTarget on all the nodes in the receiver tree, with an equal
// node from the target tree
func (node *ResourceNode) ApplyFunctionWithTarget(fn ApplyFnWithTarget, targetTree *ResourceNode) {
	for _, child := range node.Children {
		child.ApplyFunctionWithTarget(fn, targetTree)
	}

	targetNode := targetTree.GetNode(node)
	fn(node, targetNode)
}
