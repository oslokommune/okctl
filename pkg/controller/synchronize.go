// Package controller knows how to ensure desired state and current state matches
package controller

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/spf13/afero"
)

// SynchronizeOpts contains the necessary information that Synchronize() needs to do its work
type SynchronizeOpts struct {
	ClusterID api.ID

	ClusterDeclaration v1alpha1.Cluster
	DesiredTree        *resourcetree.ResourceNode

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
	opts.DesiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeCluster, CreateClusterStateRefresher(
		opts.Fs,
		opts.OutputDir,
		opts.CIDRGetter,
	))

	opts.DesiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeALBIngress, CreateALBIngressControllerRefresher(
		opts.Fs,
		opts.OutputDir,
	))

	opts.DesiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeAWSLoadBalancerController, CreateAWSLoadBalancerControllerRefresher(
		opts.Fs,
		opts.OutputDir,
	))

	opts.DesiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeExternalDNS, CreateExternalDNSStateRefresher(
		opts.PrimaryHostedZoneGetter,
	))

	opts.DesiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeIdentityManager, CreateIdentityManagerRefresher(
		opts.PrimaryHostedZoneGetter,
	))

	opts.DesiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeArgoCD, CreateArgocdStateRefresher(
		opts.IdentityPoolFetcher,
		opts.PrimaryHostedZoneGetter,
	))

	opts.DesiredTree.SetStateRefresher(resourcetree.ResourceNodeTypeNameserverDelegator, CreateNameserverDelegationStateRefresher(
		opts.PrimaryHostedZoneGetter,
	))

	currentStateTreeOpts, err := NewCreateCurrentStateTreeOpts(opts.Fs, opts.OutputDir, opts.GithubGetter, opts.PrimaryHostedZoneGetter)
	if err != nil {
		return fmt.Errorf("unable to get existing services: %w", err)
	}

	currentStateTree := CreateCurrentStateTree(currentStateTreeOpts)

	diffTree := *opts.DesiredTree

	diffTree.ApplyFunction(applyCurrentState, currentStateTree)

	return handleNode(opts.ReconciliationManager, &diffTree)
}

// handleNode knows how to run Reconcile() on every node of a ResourceNode tree
func handleNode(reconcilerManager reconciler.Reconciler, currentNode *resourcetree.ResourceNode) error {
	_, err := reconcilerManager.Reconcile(currentNode)
	if err != nil {
		return fmt.Errorf("error reconciling node: %w", err)
	}

	for _, node := range currentNode.Children {
		err = handleNode(reconcilerManager, node)
		if err != nil {
			return fmt.Errorf("error handling node: %w", err)
		}
	}

	return nil
}

// applyCurrentState knows how to apply the current state on a desired state ResourceNode tree to produce a diff that
// knows which resources to create, and which resources is already existing
func applyCurrentState(receiver *resourcetree.ResourceNode, target *resourcetree.ResourceNode) {
	if receiver.State == target.State {
		receiver.State = resourcetree.ResourceNodeStateNoop
	}
}
