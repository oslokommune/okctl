package controller

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/stretchr/testify/assert"
)

type dummyReconciler struct {
	ReconcileCounter     int
	ReconciliationResult []reconciler.ReconcilationResult
}

func (d *dummyReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeGroup
}
func (d *dummyReconciler) SetCommonMetadata(_ *resourcetree.CommonMetadata) {}

func (d *dummyReconciler) Reconcile(_ *resourcetree.ResourceNode) (reconciler.ReconcilationResult, error) {
	d.ReconcileCounter++

	return d.ReconciliationResult[d.ReconcileCounter-1], nil
}

// createRequeues creates reconciler.Reconciliation slice where all results except the last one has Requeue set to true
// Input: numberOfRequeues = 0
// Output: []reconciler.ReconciliationResult{
// 	{Requeue: false},
// }
//
// Input: numberOfRequeues = 2
// Output: []reconciler.ReconciliationResult{
// 	{Requeue: true},
// 	{Requeue: true},
//  {Requeue: false},
// }
func createRequeues(numberOfRequeues int) []reconciler.ReconcilationResult {
	requeues := make([]reconciler.ReconcilationResult, numberOfRequeues+1)

	for i := range requeues {
		requeues[i] = reconciler.ReconcilationResult{Requeue: true}
	}

	requeues[numberOfRequeues] = reconciler.ReconcilationResult{Requeue: false}

	return requeues
}

func TestHandleNode(t *testing.T) {
	testCases := []struct {
		name string

		withNumberOfRequeues int

		expectReconcileCallCount int
		expectErr                bool
	}{
		{
			name: "Should call reconcile function once on a node without requeues",

			withNumberOfRequeues: 0,

			expectReconcileCallCount: 1,
		},
		{
			name: "Should call reconcile function on a node 3 times due to a requeue",

			withNumberOfRequeues: 2,

			expectReconcileCallCount: 3,
		},
		{
			name: "Should call reconcile function on a node config.DefaultMax times with 'eternal' requeues",

			withNumberOfRequeues: constant.DefaultMaxReconciliationRequeues + 5,

			expectReconcileCallCount: constant.DefaultMaxReconciliationRequeues,
			expectErr:                true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			dummy := &dummyReconciler{
				ReconcileCounter:     0,
				ReconciliationResult: createRequeues(tc.withNumberOfRequeues),
			}

			node := &resourcetree.ResourceNode{}

			err := handleNode(dummy, node)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tc.expectReconcileCallCount, dummy.ReconcileCounter)
		})
	}
}
