package reconciliation

import (
	"context"
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/spinner"
)

// Run initiates scheduling of reconcilers
func (c *Scheduler) Run(ctx context.Context, state *clientCore.StateHandlers) (Result, error) {
	err := c.spinner.Start("reconciling")
	if err != nil {
		return Result{}, fmt.Errorf("starting spinner: %w", err)
	}

	defer func() {
		_ = c.spinner.Stop()
	}()

	queue := NewQueue(c.reconcilers)
	metadata := c.metadata()

	for reconciler := queue.Pop(); reconciler != nil; reconciler = queue.Pop() {
		subSpinner := c.spinner.SubSpinner()

		err := subSpinner.Start(reconciler.String())
		if err != nil {
			return Result{}, fmt.Errorf("starting subspinner: %w", err)
		}

		result, err := reconciler.Reconcile(ctx, metadata, state)
		if err != nil {
			return Result{}, fmt.Errorf("reconciling %s: %w", reconciler.String(), err)
		}

		if result.Requeue {
			err = queue.Push(reconciler)
			if err != nil {
				return Result{}, fmt.Errorf("passing requeue check for %s: %w", reconciler.String(), err)
			}
		}

		c.reconciliationLoopDelayFunction()

		err = subSpinner.Stop()
		if err != nil {
			return Result{}, fmt.Errorf("stopping spinner: %w", err)
		}
	}

	return Result{}, nil
}

func (c *Scheduler) metadata() Metadata {
	return Metadata{
		Out:                    c.out,
		ClusterDeclaration:     &c.clusterDeclaration,
		ApplicationDeclaration: c.applicationDeclaration,
		Purge:                  c.purgeFlag,
	}
}

// NewScheduler initializes a Scheduler
func NewScheduler(opts SchedulerOpts, reconcilers ...Reconciler) Scheduler {
	return Scheduler{
		out:     opts.Out,
		spinner: opts.Spinner,

		purgeFlag:              opts.PurgeFlag,
		clusterDeclaration:     opts.ClusterDeclaration,
		applicationDeclaration: opts.ApplicationDeclaration,

		reconciliationLoopDelayFunction: opts.ReconciliationLoopDelayFunction,
		reconcilers:                     reconcilers,
	}
}

// SchedulerOpts contains required data for scheduling reconciliations
type SchedulerOpts struct {
	// Out provides reconcilers a way to express data
	Out io.Writer
	// Spinner provides the user some eye candy
	Spinner spinner.Spinner

	// Context of the scheduling. Signifies the intent of the user
	// PurgeFlag indicates if everything should be deleted
	PurgeFlag bool
	// ReconciliationLoopDelayFunction introduces delay to the reconciliation process
	ReconciliationLoopDelayFunction func()
	ClusterDeclaration              v1alpha1.Cluster
	ApplicationDeclaration          v1alpha1.Application
}

// Scheduler knows how to run reconcilers in a reasonable way
type Scheduler struct {
	out     io.Writer
	spinner spinner.Spinner

	purgeFlag              bool
	clusterDeclaration     v1alpha1.Cluster
	applicationDeclaration v1alpha1.Application

	reconciliationLoopDelayFunction func()
	reconcilers                     []Reconciler
}
