// Package controller knows how to ensure desired state and current state matches
package controller

import (
	"fmt"
	"io"
	"time"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/config"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/spf13/afero"
)

// SynchronizeOpts contains the necessary information that Synchronize() needs to do its work
type SynchronizeOpts struct {
	Debug bool
	Out   io.Writer

	ClusterID api.ID

	ClusterDeclaration *v1alpha1.Cluster

	ReconciliationManager reconciler.Reconciler

	Fs        *afero.Afero
	OutputDir string

	GithubGetter GithubGetter

	CIDRGetter              StringFetcher
	IdentityPoolFetcher     IdentityPoolFetcher
	PrimaryHostedZoneGetter HostedZoneFetcher
}

// Synchronize knows how to discover differences between desired and actual state and rectify them
func Synchronize(opts *SynchronizeOpts) error {
	desiredTree := CreateResourceDependencyTree()
	currentStateTree := CreateResourceDependencyTree()

	setRefreshers(desiredTree, opts)

	existingResources, err := IdentifyResourcePresence(opts.Fs, opts.OutputDir, opts.GithubGetter, opts.PrimaryHostedZoneGetter)
	if err != nil {
		return fmt.Errorf("getting existing integrations: %w", err)
	}

	desiredTree.ApplyFunction(applyDeclaration(opts.ClusterDeclaration), &resourcetree.ResourceNode{})
	currentStateTree.ApplyFunction(applyExistingState(existingResources), &resourcetree.ResourceNode{})

	diffTree := *desiredTree

	diffTree.ApplyFunction(applyCurrentState, currentStateTree)

	if opts.Debug {
		fmt.Fprintf(opts.Out, "Present resources in desired tree (what is desired): \n%s\n\n", desiredTree.String())
		fmt.Fprintf(opts.Out, "Present resources in current state tree (what is currently): \n%s\n\n", currentStateTree.String())
		fmt.Fprintf(opts.Out, "Present resources in difference tree (what should be generated): \n%s\n\n", diffTree.String())
	}

	return handleNode(opts.ReconciliationManager, &diffTree)
}

// handleNode knows how to run Reconcile() on every node of a ResourceNode tree
func handleNode(reconcilerManager reconciler.Reconciler, currentNode *resourcetree.ResourceNode) (err error) {
	reconciliationResult := reconciler.ReconcilationResult{Requeue: true, RequeueAfter: 0 * time.Second}

	for requeues := 0; reconciliationResult.Requeue; requeues++ {
		if requeues == config.DefaultMaxReconciliationRequeues {
			return errors.New("maximum allowed reconciliation requeues reached")
		}

		time.Sleep(reconciliationResult.RequeueAfter)

		reconciliationResult, err = reconcilerManager.Reconcile(currentNode)
		if err != nil {
			return fmt.Errorf("reconciling node: %w", err)
		}
	}

	for _, node := range currentNode.Children {
		err = handleNode(reconcilerManager, node)
		if err != nil {
			return fmt.Errorf("handling node: %w", err)
		}
	}

	return nil
}

func applyDeclaration(declaration *v1alpha1.Cluster) resourcetree.ApplyFn {
	return func(desiredTreeNode *resourcetree.ResourceNode, _ *resourcetree.ResourceNode) {
		switch desiredTreeNode.Type {
		// Mandatory
		case resourcetree.ResourceNodeTypeGithub:
			receiver.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeZone:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeNameserverDelegator:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeVPC:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		case resourcetree.ResourceNodeTypeCluster:
			desiredTreeNode.State = resourcetree.ResourceNodeStatePresent
		// Integrations
		case resourcetree.ResourceNodeTypeALBIngress:
			desiredTreeNode.State = boolToState(declaration.Integrations.ALBIngressController)
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
		case resourcetree.ResourceNodeTypeArgoCD:
			desiredTreeNode.State = boolToState(declaration.Integrations.ArgoCD)
		}
	}
}

func applyExistingState(existingResources ExistingResources) resourcetree.ApplyFn {
	return func(receiver *resourcetree.ResourceNode, _ *resourcetree.ResourceNode) {
		switch receiver.Type {
		// Mandatory
		case resourcetree.ResourceNodeTypeGithub:
			receiver.State = boolToState(existingResources.hasGithubSetup)
		case resourcetree.ResourceNodeTypeZone:
			receiver.State = boolToState(existingResources.hasPrimaryHostedZone)
		case resourcetree.ResourceNodeTypeNameserverDelegator:
			receiver.State = boolToState(existingResources.hasDelegatedHostedZoneNameservers)
		case resourcetree.ResourceNodeTypeVPC:
			receiver.State = boolToState(existingResources.hasVPC)
		case resourcetree.ResourceNodeTypeCluster:
			receiver.State = boolToState(existingResources.hasCluster)
		// Integrations
		case resourcetree.ResourceNodeTypeALBIngress:
			receiver.State = boolToState(existingResources.hasALBIngressController)
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
		case resourcetree.ResourceNodeTypeArgoCD:
			receiver.State = boolToState(existingResources.hasArgoCD)
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

// setRefreshers sets a refresher on each node of a tree
func setRefreshers(desiredTree *resourcetree.ResourceNode, opts *SynchronizeOpts) {
	desiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeCluster, CreateClusterStateRefresher(
		opts.Fs,
		opts.OutputDir,
		opts.CIDRGetter,
	))

	desiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeALBIngress, CreateALBIngressControllerRefresher(
		opts.Fs,
		opts.OutputDir,
	))

	desiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeAWSLoadBalancerController, CreateAWSLoadBalancerControllerRefresher(
		opts.Fs,
		opts.OutputDir,
	))

	desiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeExternalDNS, CreateExternalDNSStateRefresher(
		opts.PrimaryHostedZoneGetter,
	))

	desiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeIdentityManager, CreateIdentityManagerRefresher(
		opts.PrimaryHostedZoneGetter,
	))

	desiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeArgoCD, CreateArgocdStateRefresher(
		opts.IdentityPoolFetcher,
		opts.PrimaryHostedZoneGetter,
	))

	desiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeKubePromStack, CreateKubePromStackRefresher(
		opts.IdentityPoolFetcher,
		opts.PrimaryHostedZoneGetter,
	))

	desiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeNameserverDelegator, CreateNameserverDelegationStateRefresher(
		opts.PrimaryHostedZoneGetter,
	))
}

// boolToState converts a boolean to a resourcetree.ResourceNodeState
func boolToState(present bool) resourcetree.ResourceNodeState {
	if present {
		return resourcetree.ResourceNodeStatePresent
	}

	return resourcetree.ResourceNodeStateAbsent
}
