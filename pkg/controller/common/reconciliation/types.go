package reconciliation

import (
	"context"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"io"
	"time"

	"github.com/pkg/errors"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
)

// Result contains information about the result of a Reconcile() call
type Result struct {
	// Requeue indicates if this Reconciliation must be run again
	Requeue bool
	// RequeueAfter sets the amount of delay before the requeued reconciliation should be done
	RequeueAfter time.Duration
}

// Action represents actions a Reconciler can take
type Action string

const (
	// ActionCreate indicates creation
	ActionCreate = "create"
	// ActionDelete indicates deletion
	ActionDelete = "delete"
	// ActionNoop indicates no necessary action
	ActionNoop = "noop"
	// ActionWait indicates the need to wait
	ActionWait = "wait"
)

// Metadata represents metadata required by most if not all operations on services
type Metadata struct {
	Out io.Writer

	ClusterDeclaration     *v1alpha1.Cluster
	ApplicationDeclaration v1alpha1.Application

	Purge bool
}

// Reconciler defines functions needed for the controller to use a reconciler
type Reconciler interface {
	// Reconcile knows how to do what is necessary to ensure the desired state is achieved
	Reconcile(ctx context.Context, meta Metadata, state *clientCore.StateHandlers) (Result, error)
	// String returns a name that describes the Reconciler
	String() string
}

var (
	// ErrMaximumReconciliationRequeues represents the reconciler trying a single reconciler too many times
	ErrMaximumReconciliationRequeues = errors.New(constant.MaxReconciliationReqeueusError)
	// ErrIndecisive represents the situation where the reconciler can't figure out what to do
	ErrIndecisive = errors.New(constant.IndescisiveError)
)
