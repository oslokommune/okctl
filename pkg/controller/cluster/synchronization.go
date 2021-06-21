package cluster

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	clusterrec "github.com/oslokommune/okctl/pkg/controller/cluster/reconciliation"
	"github.com/oslokommune/okctl/pkg/controller/common"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
)

// SynchronizeOpts contains the necessary information that Synchronize() needs to do its work
// nolint: maligned
type SynchronizeOpts struct {
	ID        api.ID
	Debug     bool
	Out       io.Writer
	DeleteAll bool

	ClusterDeclaration    *v1alpha1.Cluster
	ReconciliationManager reconciliation.Reconciler
	State                 *clientCore.StateHandlers
}

// DeleteCluster sets all resources to absent and initiates a reversed order reconciliation
func DeleteCluster(opts *SynchronizeOpts) error {
	desiredTree := CreateResourceDependencyTree()

	desiredTree.ApplyFunction(common.SetAllNodesAbsent)

	if opts.Debug {
		_, _ = fmt.Fprintf(opts.Out, "Resources to be deleted: \n%s\n\n", desiredTree.String())
	}

	return common.Process(opts.ReconciliationManager, opts.State, common.FlattenTreeReverse(desiredTree, []*dependencytree.Node{}))
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

	desiredTree.ApplyFunction(applyDeclaration(opts.ClusterDeclaration))

	currentStateTree.ApplyFunction(applyExistingState(existingResources))

	diffTree.ApplyFunction(applyDeclaration(opts.ClusterDeclaration))
	diffTree.ApplyFunctionWithTarget(applyCurrentState, currentStateTree)

	if opts.Debug {
		_, _ = fmt.Fprintf(opts.Out, "Present resources in desired tree (what is desired): \n%s\n\n", desiredTree.String())
		_, _ = fmt.Fprintf(opts.Out, "Present resources in current state tree (what is currently): \n%s\n\n", currentStateTree.String())
		_, _ = fmt.Fprintf(opts.Out, "Present resources in difference tree (what should be generated): \n%s\n\n", diffTree.String())
	}

	return common.Process(opts.ReconciliationManager, opts.State, common.FlattenTree(diffTree, []*dependencytree.Node{}))
}

//nolint:gocyclo,funlen
func applyDeclaration(declaration *v1alpha1.Cluster) dependencytree.ApplyFn {
	return func(desiredTreeNode *dependencytree.Node) {
		switch desiredTreeNode.Type {
		// Mandatory
		case dependencytree.NodeTypeServiceQuota:
			desiredTreeNode.State = dependencytree.NodeStatePresent
		case dependencytree.NodeTypeZone:
			desiredTreeNode.State = dependencytree.NodeStatePresent
		case dependencytree.NodeTypeNameserverDelegator:
			desiredTreeNode.State = dependencytree.NodeStatePresent
		case dependencytree.NodeTypeNameserversDelegatedTest:
			desiredTreeNode.State = dependencytree.NodeStatePresent
		case dependencytree.NodeTypeCleanupALB:
			desiredTreeNode.State = dependencytree.NodeStateNoop
		case dependencytree.NodeTypeCleanupSG:
			desiredTreeNode.State = dependencytree.NodeStateNoop
		case dependencytree.NodeTypeVPC:
			desiredTreeNode.State = dependencytree.NodeStatePresent
		case dependencytree.NodeTypeCluster:
			desiredTreeNode.State = dependencytree.NodeStatePresent
		// Integrations
		case dependencytree.NodeTypeAutoscaler:
			desiredTreeNode.State = common.BoolToState(declaration.Integrations.Autoscaler)
		case dependencytree.NodeTypeAWSLoadBalancerController:
			desiredTreeNode.State = common.BoolToState(declaration.Integrations.AWSLoadBalancerController)
		case dependencytree.NodeTypeBlockstorage:
			desiredTreeNode.State = common.BoolToState(declaration.Integrations.Blockstorage)
		case dependencytree.NodeTypeExternalDNS:
			desiredTreeNode.State = common.BoolToState(declaration.Integrations.ExternalDNS)
		case dependencytree.NodeTypeExternalSecrets:
			desiredTreeNode.State = common.BoolToState(declaration.Integrations.ExternalSecrets)
		case dependencytree.NodeTypeIdentityManager:
			desiredTreeNode.State = common.BoolToState(declaration.Integrations.Cognito)
		case dependencytree.NodeTypeKubePromStack:
			desiredTreeNode.State = common.BoolToState(declaration.Integrations.KubePromStack)
		case dependencytree.NodeTypeLoki:
			desiredTreeNode.State = common.BoolToState(declaration.Integrations.Loki)
		case dependencytree.NodeTypePromtail:
			desiredTreeNode.State = common.BoolToState(declaration.Integrations.Promtail)
		case dependencytree.NodeTypeTempo:
			desiredTreeNode.State = common.BoolToState(declaration.Integrations.Tempo)
		case dependencytree.NodeTypeArgoCD:
			desiredTreeNode.State = common.BoolToState(declaration.Integrations.ArgoCD)
		case dependencytree.NodeTypeUsers:
			desiredTreeNode.State = common.BoolToState(len(declaration.Users) > 0)
		case dependencytree.NodeTypePostgres:
			if declaration.Databases != nil {
				for _, db := range declaration.Databases.Postgres {
					node := dependencytree.NewNode(
						dependencytree.NodeType(fmt.Sprintf("%s-%s", dependencytree.NodeTypePostgresInstance, db.Name)),
					)
					desiredTreeNode.AppendChild(node)

					node.Data = &clusterrec.PostgresReconcilerState{
						DB: db,
					}
				}
			}

			desiredTreeNode.State = dependencytree.NodeStatePresent
		}
	}
}

//nolint:gocyclo,funlen
func applyExistingState(existingResources ExistingResources) dependencytree.ApplyFn {
	return func(receiver *dependencytree.Node) {
		switch receiver.Type {
		// Mandatory
		case dependencytree.NodeTypeServiceQuota:
			receiver.State = common.BoolToState(existingResources.hasServiceQuotaCheck)
		case dependencytree.NodeTypeZone:
			receiver.State = common.BoolToState(existingResources.hasPrimaryHostedZone)
		case dependencytree.NodeTypeNameserverDelegator:
			receiver.State = common.BoolToState(existingResources.hasDelegatedHostedZoneNameservers)
		case dependencytree.NodeTypeNameserversDelegatedTest:
			receiver.State = common.BoolToState(existingResources.hasDelegatedHostedZoneNameserversTest)
		case dependencytree.NodeTypeCleanupALB:
			receiver.State = dependencytree.NodeStateNoop
		case dependencytree.NodeTypeCleanupSG:
			receiver.State = dependencytree.NodeStateNoop
		case dependencytree.NodeTypeVPC:
			receiver.State = common.BoolToState(existingResources.hasVPC)
		case dependencytree.NodeTypeCluster:
			receiver.State = common.BoolToState(existingResources.hasCluster)
		// Integrations
		case dependencytree.NodeTypeAutoscaler:
			receiver.State = common.BoolToState(existingResources.hasAutoscaler)
		case dependencytree.NodeTypeAWSLoadBalancerController:
			receiver.State = common.BoolToState(existingResources.hasAWSLoadBalancerController)
		case dependencytree.NodeTypeBlockstorage:
			receiver.State = common.BoolToState(existingResources.hasBlockstorage)
		case dependencytree.NodeTypeExternalDNS:
			receiver.State = common.BoolToState(existingResources.hasExternalDNS)
		case dependencytree.NodeTypeExternalSecrets:
			receiver.State = common.BoolToState(existingResources.hasExternalSecrets)
		case dependencytree.NodeTypeIdentityManager:
			receiver.State = common.BoolToState(existingResources.hasIdentityManager)
		case dependencytree.NodeTypeKubePromStack:
			receiver.State = common.BoolToState(existingResources.hasKubePromStack)
		case dependencytree.NodeTypeLoki:
			receiver.State = common.BoolToState(existingResources.hasLoki)
		case dependencytree.NodeTypePromtail:
			receiver.State = common.BoolToState(existingResources.hasPromtail)
		case dependencytree.NodeTypeTempo:
			receiver.State = common.BoolToState(existingResources.hasTempo)
		case dependencytree.NodeTypeArgoCD:
			receiver.State = common.BoolToState(existingResources.hasArgoCD)
		case dependencytree.NodeTypeUsers:
			receiver.State = common.BoolToState(existingResources.hasUsers)
		case dependencytree.NodeTypePostgres:
		NextDb:
			for _, existingDB := range existingResources.hasPostgres {
				for _, declaredDB := range receiver.Children {
					data := declaredDB.Data.(*clusterrec.PostgresReconcilerState)

					if data.DB.Name == existingDB.Name {
						declaredDB.State = dependencytree.NodeStatePresent
						continue NextDb
					}
				}

				node := dependencytree.NewNode(
					dependencytree.NodeType(fmt.Sprintf("%s-%s", dependencytree.NodeTypePostgresInstance, existingDB.Name)),
				)
				receiver.AppendChild(node)
				node.Data = &clusterrec.PostgresReconcilerState{
					DB: *existingDB,
				}
				node.State = dependencytree.NodeStatePresent
			}

			receiver.State = dependencytree.NodeStatePresent
		}
	}
}

// applyCurrentState knows how to apply the current state on a desired state Node tree to produce a diff that
// knows which resources to create, and which resources is already existing
func applyCurrentState(receiver *dependencytree.Node, target *dependencytree.Node) {
	if target == nil {
		return
	}

	if receiver.State == target.State {
		receiver.State = dependencytree.NodeStateNoop
	}
}
