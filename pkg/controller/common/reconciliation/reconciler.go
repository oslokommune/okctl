package reconciliation

import (
	"time"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/controller/common/resourcetree"
)

// Result contains information about the result of a Reconcile() call
type Result struct {
	// Requeue indicates if this Reconciliation must be run again
	Requeue bool
	// RequeueAfter sets the amount of delay before the requeued reconciliation should be done
	RequeueAfter time.Duration
}

// Reconciler defines functions needed for the controller to use a reconciler
type Reconciler interface {
	NodeType() resourcetree.ResourceNodeType
	// Reconcile knows how to do what is necessary to ensure the desired state is achieved
	Reconcile(node *resourcetree.ResourceNode, state *clientCore.StateHandlers) (Result, error)
	// SetCommonMetadata knows how to store metadata needed by the reconciler for later use
	SetCommonMetadata(metadata *resourcetree.CommonMetadata)
}
