// Package controller knows how to ensure desired state and current state matches
package controller

import (
	"fmt"
	"io"
	"time"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// SynchronizeOpts contains the necessary information that Synchronize() needs to do its work
type SynchronizeOpts struct {
	ID    api.ID
	Debug bool
	Out   io.Writer

	ClusterDeclaration    *v1alpha1.Cluster
	ReconciliationManager reconciler.Reconciler
	StateHandlers         *clientCore.StateHandlers
}

// Synchronize knows how to discover differences between desired and actual state and rectify them
func Synchronize(opts *SynchronizeOpts) error {
	desiredTree := CreateResourceDependencyTree()
	currentStateTree := CreateResourceDependencyTree()
	diffTree := CreateResourceDependencyTree()

	existingResources, err := IdentifyResourcePresence(opts.ID, opts.StateHandlers)
	if err != nil {
		return fmt.Errorf("getting existing integrations: %w", err)
	}

	desiredTree.ApplyFunction(applyDeclaration(opts.ClusterDeclaration), &resourcetree.ResourceNode{})
	currentStateTree.ApplyFunction(applyExistingState(existingResources), &resourcetree.ResourceNode{})

	diffTree.ApplyFunction(applyDeclaration(opts.ClusterDeclaration), &resourcetree.ResourceNode{})
	diffTree.ApplyFunction(applyCurrentState, currentStateTree)

	if opts.Debug {
		_, _ = fmt.Fprintf(opts.Out, "Present resources in desired tree (what is desired): \n%s\n\n", desiredTree.String())
		_, _ = fmt.Fprintf(opts.Out, "Present resources in current state tree (what is currently): \n%s\n\n", currentStateTree.String())
		_, _ = fmt.Fprintf(opts.Out, "Present resources in difference tree (what should be generated): \n%s\n\n", diffTree.String())
	}

	return HandleNode(opts.ReconciliationManager, diffTree)
}

// HandleNode knows how to run Reconcile() on every node of a ResourceNode tree
//goland:noinspection GoNilness
func HandleNode(reconcilerManager reconciler.Reconciler, currentNode *resourcetree.ResourceNode) (err error) {
	reconciliationResult := reconciler.ReconcilationResult{Requeue: true, RequeueAfter: 0 * time.Second}

	for requeues := 0; reconciliationResult.Requeue; requeues++ {
		if requeues == constant.DefaultMaxReconciliationRequeues {
			return fmt.Errorf("maximum allowed reconciliation requeues reached: %w", err)
		}

		time.Sleep(reconciliationResult.RequeueAfter)

		reconciliationResult, err = reconcilerManager.Reconcile(currentNode)
		if err != nil && !reconciliationResult.Requeue {
			return fmt.Errorf("reconciling node: %w", err)
		}
	}

	for _, node := range currentNode.Children {
		err = HandleNode(reconcilerManager, node)
		if err != nil {
			return err
		}
	}

	return nil
}

// nolint: gocyclo
func applyDeclaration(declaration *v1alpha1.Cluster) resourcetree.ApplyFn {
	return func(desiredTreeNode *resourcetree.ResourceNode, _ *resourcetree.ResourceNode) {
		switch desiredTreeNode.Type {
		// Mandatory
		case resourcetree.ResourceNodeTypeZone:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeNameserverDelegator:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeVPC:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeCluster:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		// Integrations
		case resourcetree.ResourceNodeTypeAutoscaler:
			desiredTreeNode.State = boolToState(declaration.Integrations.Autoscaler)
		case resourcetree.ResourceNodeTypeAWSLoadBalancerController:
			desiredTreeNode.State = boolToState(declaration.Integrations.AWSLoadBalancerController)
		case resourcetree.ResourceNodeTypeBlockstorage:
			desiredTreeNode.State = boolToState(declaration.Integrations.Blockstorage)
		case resourcetree.ResourceNodeTypeExternalDNS:
			desiredTreeNode.State = boolToState(declaration.Integrations.ExternalDNS)
		case resourcetree.ResourceNodeTypeExternalSecrets:
			desiredTreeNode.State = boolToState(declaration.Integrations.ExternalSecrets)
		case resourcetree.ResourceNodeTypeIdentityManager:
			desiredTreeNode.State = boolToState(declaration.Integrations.Cognito)
		case resourcetree.ResourceNodeTypeKubePromStack:
			desiredTreeNode.State = boolToState(declaration.Integrations.KubePromStack)
		case resourcetree.ResourceNodeTypeLoki:
			desiredTreeNode.State = boolToState(declaration.Integrations.Loki)
		case resourcetree.ResourceNodeTypePromtail:
			desiredTreeNode.State = boolToState(declaration.Integrations.Promtail)
		case resourcetree.ResourceNodeTypeTempo:
			desiredTreeNode.State = boolToState(declaration.Integrations.Tempo)
		case resourcetree.ResourceNodeTypeArgoCD:
			desiredTreeNode.State = boolToState(declaration.Integrations.ArgoCD)
		case resourcetree.ResourceNodeTypeUsers:
			desiredTreeNode.State = boolToState(len(declaration.Users) > 0)
		case resourcetree.ResourceNodeTypePostgres:
			desiredTreeNode.State = boolToState(false)

			if declaration.Databases != nil {
				desiredTreeNode.State = boolToState(len(declaration.Databases.Postgres) > 0)
			}
		}
	}
}

// nolint: gocyclo
func applyExistingState(existingResources ExistingResources) resourcetree.ApplyFn {
	return func(receiver *resourcetree.ResourceNode, _ *resourcetree.ResourceNode) {
		switch receiver.Type {
		// Mandatory
		case resourcetree.ResourceNodeTypeZone:
			receiver.State = boolToState(existingResources.hasPrimaryHostedZone)
		case resourcetree.ResourceNodeTypeNameserverDelegator:
			receiver.State = boolToState(existingResources.hasDelegatedHostedZoneNameservers)
		case resourcetree.ResourceNodeTypeVPC:
			receiver.State = boolToState(existingResources.hasVPC)
		case resourcetree.ResourceNodeTypeCluster:
			receiver.State = boolToState(existingResources.hasCluster)
		// Integrations
		case resourcetree.ResourceNodeTypeAutoscaler:
			receiver.State = boolToState(existingResources.hasAutoscaler)
		case resourcetree.ResourceNodeTypeAWSLoadBalancerController:
			receiver.State = boolToState(existingResources.hasAWSLoadBalancerController)
		case resourcetree.ResourceNodeTypeBlockstorage:
			receiver.State = boolToState(existingResources.hasBlockstorage)
		case resourcetree.ResourceNodeTypeExternalDNS:
			receiver.State = boolToState(existingResources.hasExternalDNS)
		case resourcetree.ResourceNodeTypeExternalSecrets:
			receiver.State = boolToState(existingResources.hasExternalSecrets)
		case resourcetree.ResourceNodeTypeIdentityManager:
			receiver.State = boolToState(existingResources.hasIdentityManager)
		case resourcetree.ResourceNodeTypeKubePromStack:
			receiver.State = boolToState(existingResources.hasKubePromStack)
		case resourcetree.ResourceNodeTypeLoki:
			receiver.State = boolToState(existingResources.hasLoki)
		case resourcetree.ResourceNodeTypePromtail:
			receiver.State = boolToState(existingResources.hasPromtail)
		case resourcetree.ResourceNodeTypeTempo:
			receiver.State = boolToState(existingResources.hasTempo)
		case resourcetree.ResourceNodeTypeArgoCD:
			receiver.State = boolToState(existingResources.hasArgoCD)
		case resourcetree.ResourceNodeTypeUsers:
			receiver.State = boolToState(existingResources.hasUsers)
		case resourcetree.ResourceNodeTypePostgres:
			receiver.State = boolToState(existingResources.hasPostgres)
		}
	}
}

// applyCurrentState knows how to apply the current state on a desired state ResourceNode tree to produce a diff that
// knows which resources to create, and which resources is already existing
func applyCurrentState(receiver *resourcetree.ResourceNode, target *resourcetree.ResourceNode) {
	if receiver.State == target.State {
		receiver.State = resourcetree.ResourceNodeStateNoop
	}
}

// boolToState converts a boolean to a resourcetree.ResourceNodeState
func boolToState(present bool) resourcetree.ResourceNodeState {
	if present {
		return resourcetree.ResourceNodeStatePresent
	}

	return resourcetree.ResourceNodeStateAbsent
}
