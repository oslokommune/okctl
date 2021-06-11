package reconciliation

import (
	"context"
	"testing"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/stretchr/testify/assert"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
)

func TestBasicQueueFunctionality(t *testing.T) {
	testCases := []struct {
		name string
		with []Reconciler
	}{
		{
			name: "Should produce an identical list of elements after pushing and popping",
			with: []Reconciler{},
		},
		{
			name: "Should produce an identical list of elements after pushing and popping",
			with: []Reconciler{
				mockReconciler{identifier: "a"},
				mockReconciler{identifier: "b"},
				mockReconciler{identifier: "c"},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			q := NewQueue([]Reconciler{})

			for _, reconciler := range tc.with {
				_ = q.Push(reconciler)
			}

			result := make([]Reconciler, 0)

			for reconciler := q.Pop(); reconciler != nil; reconciler = q.Pop() {
				result = append(result, reconciler)
			}

			assert.Equal(t, tc.with, result)
		})
	}
}

func TestErrorUponMaximumRequeues(t *testing.T) {
	testCases := []struct {
		name         string
		withRequeues int
		expectError  error
	}{
		{
			name:         "Should not err with one requeue",
			withRequeues: 1,
			expectError:  nil,
		},
		{
			name:         "Should not err upon maximum reconciliation requeues reached but not exceeded for a single reconciler",
			withRequeues: constant.DefaultMaxReconciliationRequeues,
			expectError:  nil,
		},
		{
			name:         "Should err upon maximum reconciliation requeues exceeded for a single reconciler",
			withRequeues: constant.DefaultMaxReconciliationRequeues + 1,
			expectError:  ErrMaximumReconciliationRequeues,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			q := NewQueue([]Reconciler{})

			var err error
			for index := 0; index < tc.withRequeues; index++ {
				err = q.Push(mockReconciler{identifier: "dummy"})

				if err != nil {
					break
				}
			}

			assert.Equal(t, tc.expectError, err)
		})
	}
}

type mockReconciler struct {
	identifier string
}

func (m mockReconciler) String() string {
	return m.identifier
}

func (m mockReconciler) Reconcile(_ context.Context, _ Metadata, _ *clientCore.StateHandlers) (Result, error) {
	panic("implement me")
}
