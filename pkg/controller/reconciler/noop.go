package reconciler

import (
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// NoopReconciler handles reconciliation for dummy nodes (e.g. the root node) and acts as a template for other
// reconcilers
type NoopReconciler struct{}

// SetCommonMetadata knows how to store common metadata on the reconciler. This should do nothing if common metadata is
// not needed
func (receiver *NoopReconciler) SetCommonMetadata(_ *resourcetree.CommonMetadata) {
	// Do nothing because a NoopReconciler does nothing and therefore does need to store any common metadata
}

// Reconcile knows how to create, update and delete the relevant resource
func (receiver *NoopReconciler) Reconcile(_ *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	return &ReconcilationResult{Requeue: false}, nil
}
