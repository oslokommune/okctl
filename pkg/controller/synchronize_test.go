package controller

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/sebdah/goldie/v2"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/controller/common/resourcetree"
	"github.com/stretchr/testify/assert"
)

type dummyReconciler struct {
	ReconcileCounter     int
	ReconciliationResult []reconciliation.Result
}

func (d *dummyReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeGroup
}
func (d *dummyReconciler) SetCommonMetadata(_ *resourcetree.CommonMetadata) {}

func (d *dummyReconciler) Reconcile(_ *resourcetree.ResourceNode, _ *clientCore.StateHandlers) (reconciliation.Result, error) {
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
func createRequeues(numberOfRequeues int) []reconciliation.Result {
	requeues := make([]reconciliation.Result, numberOfRequeues+1)

	for i := range requeues {
		requeues[i] = reconciliation.Result{Requeue: true}
	}

	requeues[numberOfRequeues] = reconciliation.Result{Requeue: false}

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

			err := Process(dummy, nil, FlattenTree(node, []*resourcetree.ResourceNode{}))
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tc.expectReconcileCallCount, dummy.ReconcileCounter)
		})
	}
}

type mockAlwaysErrorReconciler struct {
	iteration            int
	ReconciliationResult []reconciliation.Result
}

func (m *mockAlwaysErrorReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeGroup
}

func (m *mockAlwaysErrorReconciler) Reconcile(_ *resourcetree.ResourceNode, _ *clientCore.StateHandlers) (reconciliation.Result, error) {
	m.iteration++

	return m.ReconciliationResult[m.iteration-1], errors.New("dummy err")
}

func (m *mockAlwaysErrorReconciler) SetCommonMetadata(_ *resourcetree.CommonMetadata) {}

func TestReceivedErrorAfterRequeues(t *testing.T) {
	testCases := []struct {
		name string

		withResults []reconciliation.Result

		expectErrorAfterIterations int
		expectError                error
	}{
		{
			name: "Should break out of Process immediately when requeue is false",

			withResults: []reconciliation.Result{{Requeue: false}},

			expectErrorAfterIterations: 1,
			expectError:                errors.New("reconciling node (group): dummy err"),
		},
		{
			name: "Should break out of Process after second reconciliation when requeues are true, false",

			withResults: []reconciliation.Result{{Requeue: true}, {Requeue: false}},

			expectErrorAfterIterations: 2,
			expectError:                errors.New("reconciling node (group): dummy err"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			r := &mockAlwaysErrorReconciler{ReconciliationResult: tc.withResults}

			err := Process(r, nil, FlattenTree(&resourcetree.ResourceNode{Type: resourcetree.ResourceNodeTypeGroup}, []*resourcetree.ResourceNode{}))
			assert.NotNil(t, err)

			assert.Equal(t, tc.expectErrorAfterIterations, r.iteration)
			assert.Equal(t, tc.expectError.Error(), err.Error())
		})
	}
}

func TestOrderTree(t *testing.T) {
	tree := CreateResourceDependencyTree()

	order := FlattenTree(tree, []*resourcetree.ResourceNode{})
	reverse := FlattenTreeReverse(tree, []*resourcetree.ResourceNode{})

	g := goldie.New(t)

	g.AssertJson(t, "tree-order.json", struct {
		Normal  []*resourcetree.ResourceNode
		Reverse []*resourcetree.ResourceNode
	}{
		order,
		reverse,
	})
}
