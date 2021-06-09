package reconciliation

import (
	"context"
	"io"
	"time"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// Result contains information about the result of a Reconcile() call
type Result struct {
	// Requeue indicates if this Reconciliation must be run again
	Requeue bool
	// RequeueAfter sets the amount of delay before the requeued reconciliation should be done
	RequeueAfter time.Duration
}

// CommonMetadata represents metadata required by most if not all operations on services
type CommonMetadata struct {
	Ctx context.Context

	Out io.Writer

	ClusterID              api.ID
	Declaration            *v1alpha1.Cluster
	ApplicationDeclaration v1alpha1.Application
}

// Reconciler defines functions needed for the controller to use a reconciler
type Reconciler interface {
	NodeType() dependencytree.NodeType
	// Reconcile knows how to do what is necessary to ensure the desired state is achieved
	Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (Result, error)
	// SetCommonMetadata knows how to store metadata needed by the reconciler for later use
	SetCommonMetadata(metadata *CommonMetadata)
}
