package reconciliation

import (
	"context"
	"errors"
	"io"
	"testing"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/stretchr/testify/assert"
)

func generateDummySpinner() spinner.Spinner {
	spin, _ := spinner.New("", io.Discard)

	return spin
}

type requeueReconciler struct {
	getResult      func() Result
	bumpReconciles func()
}

func (r requeueReconciler) String() string { return "" }

func (r requeueReconciler) Reconcile(_ context.Context, _ Metadata, _ *clientCore.StateHandlers) (Result, error) {
	r.bumpReconciles()

	return r.getResult(), nil
}

//nolint:funlen
func TestRequeueing(t *testing.T) {
	testCases := []struct {
		name string

		withResults []Result
		expectRuns  int
	}{
		{
			name: "Should work with no requeues",

			withResults: []Result{
				{Requeue: false},
				{Requeue: false},
			},
			expectRuns: 1,
		},
		{
			name: "Should work with a single requeue",

			withResults: []Result{
				{Requeue: true},
				{Requeue: false},
			},
			expectRuns: 2,
		},
		{
			name: "Should work with multiple requeues",

			withResults: []Result{
				{Requeue: true},
				{Requeue: true},
				{Requeue: false},
			},
			expectRuns: 3,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			resultIndex := 0
			reconciles := 0

			getResult := func() Result {
				result := tc.withResults[resultIndex]

				resultIndex++

				return result
			}

			bumpReconciles := func() {
				reconciles++
			}

			scheduler := NewScheduler(SchedulerOpts{
				Spinner:                         generateDummySpinner(),
				ReconciliationLoopDelayFunction: func() {},
			}, requeueReconciler{
				getResult:      getResult,
				bumpReconciles: bumpReconciles,
			})

			_, err := scheduler.Run(context.Background(), nil)
			assert.Nil(t, err)

			assert.Equal(t, tc.expectRuns, reconciles)
		})
	}
}

type noopReconciler struct{}

func (n noopReconciler) Reconcile(_ context.Context, _ Metadata, _ *clientCore.StateHandlers) (Result, error) {
	panic("implement me")
}
func (n noopReconciler) String() string { return "" }

func TestQueueingPops(t *testing.T) {
	testCases := []struct {
		name string

		withReconcilers []Reconciler
		expectPops      int
	}{
		{
			name:            "Should work with zero elements",
			withReconcilers: []Reconciler{},
			expectPops:      0,
		},
		{
			name:            "Should work a single reconciler",
			withReconcilers: []Reconciler{noopReconciler{}},
			expectPops:      1,
		},
		{
			name:            "Should work several reconcilers",
			withReconcilers: []Reconciler{noopReconciler{}, noopReconciler{}, noopReconciler{}, noopReconciler{}},
			expectPops:      4,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			q := NewQueue(tc.withReconcilers)
			popsCounter := 0

			for item := q.Pop(); item != nil; item = q.Pop() {
				popsCounter++
			}

			assert.Equal(t, tc.expectPops, popsCounter)
		})
	}
}

func TestQueueingPush(t *testing.T) {
	testCases := []struct {
		name string

		withReconcilers []Reconciler
	}{
		{
			name:            "Should work with zero reconcilers",
			withReconcilers: []Reconciler{},
		},
		{
			name:            "Should work with a single reconciler",
			withReconcilers: []Reconciler{noopReconciler{}},
		},
		{
			name:            "Should work with multiple reconcilers",
			withReconcilers: []Reconciler{noopReconciler{}},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			q := NewQueue([]Reconciler{})

			for _, reconciler := range tc.withReconcilers {
				_ = q.Push(reconciler)
			}

			popsCounter := 0
			for item := q.Pop(); item != nil; item = q.Pop() {
				popsCounter++
			}

			assert.Equal(t, len(tc.withReconcilers), popsCounter)
		})
	}
}

func TestNoMutation(t *testing.T) {
	testCases := []struct {
		name string

		withReconcilers []Reconciler
	}{
		{
			name:            "Should work with zero elements",
			withReconcilers: []Reconciler{},
		},
		{
			name:            "Should work a single reconciler",
			withReconcilers: []Reconciler{noopReconciler{}},
		},
		{
			name:            "Should work several reconcilers",
			withReconcilers: []Reconciler{noopReconciler{}, noopReconciler{}, noopReconciler{}, noopReconciler{}},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			originalLength := len(tc.withReconcilers)

			q := NewQueue(tc.withReconcilers)

			for item := q.Pop(); item != nil; item = q.Pop() {
			}

			assert.Equal(t, originalLength, len(tc.withReconcilers))
		})
	}
}

type deadlockReconciler struct {
	reconcileFn func() (Result, error)
}

func (d deadlockReconciler) String() string { return "" }

func (d deadlockReconciler) Reconcile(_ context.Context, _ Metadata, _ *clientCore.StateHandlers) (Result, error) {
	if d.reconcileFn == nil {
		return Result{}, nil
	}

	return d.reconcileFn()
}

func TestPreflightDeadlock(t *testing.T) {
	testCases := []struct {
		name            string
		withReconcilers []Reconciler
		expectErr       error
	}{
		{
			name: "Should prevent deadlock when upon eternal requeue requests",
			withReconcilers: []Reconciler{
				deadlockReconciler{
					reconcileFn: func() (Result, error) {
						return Result{Requeue: true}, nil
					},
				},
			},
			expectErr: ErrMaximumReconciliationRequeues,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			scheduler := NewScheduler(SchedulerOpts{
				Spinner:                         generateDummySpinner(),
				ReconciliationLoopDelayFunction: func() {},
			}, tc.withReconcilers...)

			_, err := scheduler.Run(context.Background(), nil)

			assert.True(t, errors.Is(err, tc.expectErr))
		})
	}
}
