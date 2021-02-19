// Package reconciler contains different reconcilers for each of the necessary resources
package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ReconcilationResult contains information about the result of a Reconcile() call
type ReconcilationResult struct {
	// Requeue indicates if this Reconciliation must be run again
	Requeue bool
}

// Reconciler defines functions needed for the controller to use a reconciler
type Reconciler interface {
	// Reconcile knows how to do what is necessary to ensure the desired state is achieved
	Reconcile(*resourcetree.ResourceNode) (*ReconcilationResult, error)
	// SetCommonMetadata knows how to store metadata needed by the reconciler for later use
	SetCommonMetadata(metadata *resourcetree.CommonMetadata)
}

// Manager provides a simpler way to organize reconcilers
type Manager struct {
	commonMetadata *resourcetree.CommonMetadata
	Reconcilers    map[resourcetree.ResourceNodeType]Reconciler
}

// AddReconciler makes a Reconciler available in the Manager
func (manager *Manager) AddReconciler(key resourcetree.ResourceNodeType, reconciler Reconciler) {
	reconciler.SetCommonMetadata(manager.commonMetadata)

	manager.Reconcilers[key] = reconciler
}

// Reconcile chooses the correct reconciler to use based on a nodes type
func (manager *Manager) Reconcile(node *resourcetree.ResourceNode) (result *ReconcilationResult, err error) {
	spinner := manager.commonMetadata.Spin.SubSpinner()

	err = spinner.Start(resourcetree.ResourceNodeTypeToString(node.Type))
	if err != nil {
		return nil, fmt.Errorf("starting subspinner: %w", err)
	}

	defer func() {
		_ = spinner.Stop()
	}()

	node.RefreshState()

	return manager.Reconcilers[node.Type].Reconcile(node)
}

// NewReconcilerManager creates a new Manager with a NoopReconciler already installed
func NewReconcilerManager(metadata *resourcetree.CommonMetadata) *Manager {
	return &Manager{
		commonMetadata: metadata,
		Reconcilers: map[resourcetree.ResourceNodeType]Reconciler{
			resourcetree.ResourceNodeTypeGroup: &NoopReconciler{},
		},
	}
}
