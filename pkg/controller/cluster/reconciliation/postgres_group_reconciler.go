package reconciliation

import (
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// PostgresGroupReconciler handles reconciliation for dummy nodes (e.g. the root node) and acts as a template for other
// reconcilers
type PostgresGroupReconciler struct{}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (receiver *PostgresGroupReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypePostgres
}

// SetCommonMetadata knows how to store common metadata on the reconciler. This should do nothing if common metadata is
// not needed
func (receiver *PostgresGroupReconciler) SetCommonMetadata(_ *resourcetree.CommonMetadata) {
	// Do nothing because a PostgresGroupReconciler does nothing and therefore does need to store any common metadata
}

// Reconcile knows how to create, update and delete the relevant resource
func (receiver *PostgresGroupReconciler) Reconcile(_ *resourcetree.ResourceNode, _ *clientCore.StateHandlers) (reconciliation.Result, error) {
	return reconciliation.Result{
		Requeue: false,
	}, nil
}
