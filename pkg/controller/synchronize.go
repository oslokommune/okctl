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
// nolint: maligned
type SynchronizeOpts struct {
	ID        api.ID
	Debug     bool
	Out       io.Writer
	DeleteAll bool

	ClusterDeclaration    *v1alpha1.Cluster
	ReconciliationManager reconciler.Reconciler
	State                 *clientCore.StateHandlers
}

// DeleteCluster sets all resources to absent and initiates a reversed order reconciliation
func DeleteCluster(opts *SynchronizeOpts) error {
	desiredTree := CreateResourceDependencyTree()

	desiredTree.ApplyFunction(setStateAbsent(), &resourcetree.ResourceNode{})

	if opts.Debug {
		_, _ = fmt.Fprintf(opts.Out, "Resources to be deleted: \n%s\n\n", desiredTree.String())
	}

	return Process(opts.ReconciliationManager, opts.State, FlattenTreeReverse(desiredTree, []*resourcetree.ResourceNode{}))
}

// Synchronize knows how to discover differences between desired and actual state and rectify them
func Synchronize(opts *SynchronizeOpts) error {
	desiredTree := CreateResourceDependencyTree()
	currentStateTree := CreateResourceDependencyTree()
	diffTree := CreateResourceDependencyTree()

	existingResources, err := IdentifyResourcePresence(opts.ID, opts.State)
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

	return Process(opts.ReconciliationManager, opts.State, FlattenTree(diffTree, []*resourcetree.ResourceNode{}))
}

// FlattenTree flattens the tree to an execution order
func FlattenTree(current *resourcetree.ResourceNode, order []*resourcetree.ResourceNode) []*resourcetree.ResourceNode {
	cpy := *current
	cpy.Children = nil

	order = append(order, &cpy)

	for _, node := range current.Children {
		order = FlattenTree(node, order)
	}

	return order
}

// FlattenTreeReverse flattens the tree to a reverse execution order
func FlattenTreeReverse(current *resourcetree.ResourceNode, order []*resourcetree.ResourceNode) []*resourcetree.ResourceNode {
	order = FlattenTree(current, order)

	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	return order
}

// Process knows how to run Reconcile() on every node of a ResourceNode tree
//goland:noinspection GoNilness
func Process(reconcilerManager reconciler.Reconciler, state *clientCore.StateHandlers, order []*resourcetree.ResourceNode) (err error) {
	for _, node := range order {
		result := reconciler.ReconcilationResult{
			Requeue:      true,
			RequeueAfter: 0 * time.Second,
		}

		for requeues := 0; result.Requeue; requeues++ {
			if requeues == constant.DefaultMaxReconciliationRequeues {
				return fmt.Errorf("maximum allowed reconciliation requeues reached: %w", err)
			}

			time.Sleep(result.RequeueAfter)

			result, err = reconcilerManager.Reconcile(node, state)
			if err != nil && !result.Requeue {
				return fmt.Errorf("reconciling node (%s): %w", node.Type.String(), err)
			}
		}
	}

	return nil
}

// nolint: gocyclo
func applyDeclaration(declaration *v1alpha1.Cluster) resourcetree.ApplyFn {
	return func(desiredTreeNode *resourcetree.ResourceNode, _ *resourcetree.ResourceNode) {
		switch desiredTreeNode.Type {
		// Mandatory
		case resourcetree.ResourceNodeTypeServiceQuota:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeZone:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeNameserverDelegator:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeNameserversDelegatedTest:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeCleanupALB:
			desiredTreeNode.State = resourcetree.ResourceNodeStateNoop
		case resourcetree.ResourceNodeTypeCleanupSG:
			desiredTreeNode.State = resourcetree.ResourceNodeStateNoop
		case resourcetree.ResourceNodeTypeVPC:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeCluster:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		// Integrations
		case resourcetree.ResourceNodeTypeAutoscaler:
			desiredTreeNode.State = BoolToState(declaration.Integrations.Autoscaler)
		case resourcetree.ResourceNodeTypeAWSLoadBalancerController:
			desiredTreeNode.State = BoolToState(declaration.Integrations.AWSLoadBalancerController)
		case resourcetree.ResourceNodeTypeBlockstorage:
			desiredTreeNode.State = BoolToState(declaration.Integrations.Blockstorage)
		case resourcetree.ResourceNodeTypeExternalDNS:
			desiredTreeNode.State = BoolToState(declaration.Integrations.ExternalDNS)
		case resourcetree.ResourceNodeTypeExternalSecrets:
			desiredTreeNode.State = BoolToState(declaration.Integrations.ExternalSecrets)
		case resourcetree.ResourceNodeTypeIdentityManager:
			desiredTreeNode.State = BoolToState(declaration.Integrations.Cognito)
		case resourcetree.ResourceNodeTypeKubePromStack:
			desiredTreeNode.State = BoolToState(declaration.Integrations.KubePromStack)
		case resourcetree.ResourceNodeTypeLoki:
			desiredTreeNode.State = BoolToState(declaration.Integrations.Loki)
		case resourcetree.ResourceNodeTypePromtail:
			desiredTreeNode.State = BoolToState(declaration.Integrations.Promtail)
		case resourcetree.ResourceNodeTypeTempo:
			desiredTreeNode.State = BoolToState(declaration.Integrations.Tempo)
		case resourcetree.ResourceNodeTypeArgoCD:
			desiredTreeNode.State = BoolToState(declaration.Integrations.ArgoCD)
		case resourcetree.ResourceNodeTypeUsers:
			desiredTreeNode.State = BoolToState(len(declaration.Users) > 0)
		case resourcetree.ResourceNodeTypePostgres:
			if declaration.Databases != nil {
				for _, db := range declaration.Databases.Postgres {
					node := createNode(
						desiredTreeNode,
						resourcetree.ResourceNodeType(fmt.Sprintf("%s-%s", resourcetree.ResourceNodeTypePostgresInstance, db.Name)),
					)
					node.Data = &reconciler.PostgresReconcilerState{
						DB: db,
					}
				}
			}

			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		}
	}
}

//nolint:gocyclo,funlen
func applyExistingState(existingResources ExistingResources) resourcetree.ApplyFn {
	return func(receiver *resourcetree.ResourceNode, _ *resourcetree.ResourceNode) {
		switch receiver.Type {
		// Mandatory
		case resourcetree.ResourceNodeTypeServiceQuota:
			receiver.State = BoolToState(existingResources.hasServiceQuotaCheck)
		case resourcetree.ResourceNodeTypeZone:
			receiver.State = BoolToState(existingResources.hasPrimaryHostedZone)
		case resourcetree.ResourceNodeTypeNameserverDelegator:
			receiver.State = BoolToState(existingResources.hasDelegatedHostedZoneNameservers)
		case resourcetree.ResourceNodeTypeNameserversDelegatedTest:
			receiver.State = BoolToState(existingResources.hasDelegatedHostedZoneNameserversTest)
		case resourcetree.ResourceNodeTypeCleanupALB:
			receiver.State = resourcetree.ResourceNodeStateNoop
		case resourcetree.ResourceNodeTypeCleanupSG:
			receiver.State = resourcetree.ResourceNodeStateNoop
		case resourcetree.ResourceNodeTypeVPC:
			receiver.State = BoolToState(existingResources.hasVPC)
		case resourcetree.ResourceNodeTypeCluster:
			receiver.State = BoolToState(existingResources.hasCluster)
		// Integrations
		case resourcetree.ResourceNodeTypeAutoscaler:
			receiver.State = BoolToState(existingResources.hasAutoscaler)
		case resourcetree.ResourceNodeTypeAWSLoadBalancerController:
			receiver.State = BoolToState(existingResources.hasAWSLoadBalancerController)
		case resourcetree.ResourceNodeTypeBlockstorage:
			receiver.State = BoolToState(existingResources.hasBlockstorage)
		case resourcetree.ResourceNodeTypeExternalDNS:
			receiver.State = BoolToState(existingResources.hasExternalDNS)
		case resourcetree.ResourceNodeTypeExternalSecrets:
			receiver.State = BoolToState(existingResources.hasExternalSecrets)
		case resourcetree.ResourceNodeTypeIdentityManager:
			receiver.State = BoolToState(existingResources.hasIdentityManager)
		case resourcetree.ResourceNodeTypeKubePromStack:
			receiver.State = BoolToState(existingResources.hasKubePromStack)
		case resourcetree.ResourceNodeTypeLoki:
			receiver.State = BoolToState(existingResources.hasLoki)
		case resourcetree.ResourceNodeTypePromtail:
			receiver.State = BoolToState(existingResources.hasPromtail)
		case resourcetree.ResourceNodeTypeTempo:
			receiver.State = BoolToState(existingResources.hasTempo)
		case resourcetree.ResourceNodeTypeArgoCD:
			receiver.State = BoolToState(existingResources.hasArgoCD)
		case resourcetree.ResourceNodeTypeUsers:
			receiver.State = BoolToState(existingResources.hasUsers)
		case resourcetree.ResourceNodeTypePostgres:
		NextDb:
			for _, existingDB := range existingResources.hasPostgres {
				for _, declaredDB := range receiver.Children {
					data, ok := declaredDB.Data.(reconciler.PostgresReconcilerState)
					if !ok {
						panic("could not cast to database state")
					}

					if data.DB.Name == existingDB.Name {
						declaredDB.State = resourcetree.ResourceNodeStatePresent
						continue NextDb
					}
				}

				node := createNode(
					receiver,
					resourcetree.ResourceNodeType(fmt.Sprintf("%s-%s", resourcetree.ResourceNodeTypePostgresInstance, existingDB.Name)),
				)
				node.Data = &reconciler.PostgresReconcilerState{
					DB: *existingDB,
				}
				node.State = resourcetree.ResourceNodeStatePresent
			}

			receiver.State = resourcetree.ResourceNodeStatePresent
		}
	}
}

// setStateAbsent sets the state as absent
func setStateAbsent() resourcetree.ApplyFn {
	return func(receiver *resourcetree.ResourceNode, _ *resourcetree.ResourceNode) {
		receiver.State = resourcetree.ResourceNodeStateAbsent
	}
}

// applyCurrentState knows how to apply the current state on a desired state ResourceNode tree to produce a diff that
// knows which resources to create, and which resources is already existing
func applyCurrentState(receiver *resourcetree.ResourceNode, target *resourcetree.ResourceNode) {
	if target == nil {
		return
	}

	if receiver.State == target.State {
		receiver.State = resourcetree.ResourceNodeStateNoop
	}
}

// BoolToState converts a boolean to a resourcetree.ResourceNodeState
func BoolToState(present bool) resourcetree.ResourceNodeState {
	if present {
		return resourcetree.ResourceNodeStatePresent
	}

	return resourcetree.ResourceNodeStateAbsent
}
