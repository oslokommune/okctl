package reconciliation

import (
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// NoopReconciler handles reconciliation for dummy nodes (e.g. the root node) and acts as a template for other
// reconcilers
type NoopReconciler struct{}

// NodeType returns the relevant NodeType for this reconciler
func (receiver *NoopReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeGroup
}

// SetCommonMetadata knows how to store common metadata on the reconciler. This should do nothing if common metadata is
// not needed
func (receiver *NoopReconciler) SetCommonMetadata(_ *CommonMetadata) {
	// Do nothing because a NoopReconciler does nothing and therefore does need to store any common metadata
}

// Reconcile knows how to create, update and delete the relevant resource
func (receiver *NoopReconciler) Reconcile(_ *dependencytree.Node, _ *clientCore.StateHandlers) (Result, error) {
	return Result{Requeue: false}, nil
}
