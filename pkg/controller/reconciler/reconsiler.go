// Package reconciler contains different reconcilers for each of the necessary resources
package reconciler

import (
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

/*
 * Reconciler
 */

// ReconcilationResult contains information about the result of a Reconcile() call
type ReconcilationResult struct {
	// Requeue indicates if this Reconciliation must be run again
	Requeue bool
}

type Reconciler interface {
	// Reconcile knows how to do what is necessary to ensure the desired state is achieved
	Reconcile(*resourcetree.ResourceNode) (*ReconcilationResult, error)
	SetCommonMetadata(metadata *resourcetree.CommonMetadata)
}

/*
ReconcilerManager provides a simpler way to organize reconcilers
*/
type ReconcilerManager struct {
	commonMetadata *resourcetree.CommonMetadata
	Reconcilers    map[resourcetree.ResourceNodeType]Reconciler
}

// AddReconciler makes a Reconciler available in the ReconcilerManager
func (manager *ReconcilerManager) AddReconciler(key resourcetree.ResourceNodeType, reconciler Reconciler) {
	reconciler.SetCommonMetadata(manager.commonMetadata)

	manager.Reconcilers[key] = reconciler
}

// Reconcile chooses the correct reconciler to use based on a nodes type
func (manager *ReconcilerManager) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	node.RefreshState()

	return manager.Reconcilers[node.Type].Reconcile(node)
}

// NewReconcilerManager creates a new ReconcilerManager with a NoopReconciler already installed
func NewReconcilerManager(metadata *resourcetree.CommonMetadata) *ReconcilerManager {
	return &ReconcilerManager{
		commonMetadata: metadata,
		Reconcilers: map[resourcetree.ResourceNodeType]Reconciler{
			resourcetree.ResourceNodeTypeGroup: &NoopReconciler{},
		},
	}
}
