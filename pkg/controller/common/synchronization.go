package common

import (
	"fmt"
	"time"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
	"github.com/oslokommune/okctl/pkg/controller/common/resourcetree"
)

// FlattenTree flattens the tree to an execution order
func FlattenTree(current *resourcetree.ResourceNode, order []*resourcetree.ResourceNode) []*resourcetree.ResourceNode {
	cpy := *current
	cpy.Children = nil

	order = append(order, &cpy)

	for _, node := range current.Children {
		order = FlattenTree(node, order)
	}

	return order
}

// FlattenTreeReverse flattens the tree to a reverse execution order
func FlattenTreeReverse(current *resourcetree.ResourceNode, order []*resourcetree.ResourceNode) []*resourcetree.ResourceNode {
	order = FlattenTree(current, order)

	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}

	return order
}

// Process knows how to run Reconcile() on every node of a ResourceNode tree
//goland:noinspection GoNilness
func Process(reconcilerManager reconciliation.Reconciler, state *clientCore.StateHandlers, order []*resourcetree.ResourceNode) (err error) {
	for _, node := range order {
		result := reconciliation.Result{
			Requeue:      true,
			RequeueAfter: 0 * time.Second,
		}

		for requeues := 0; result.Requeue; requeues++ {
			if requeues == constant.DefaultMaxReconciliationRequeues {
				return fmt.Errorf("maximum allowed reconciliation requeues reached: %w", err)
			}

			time.Sleep(result.RequeueAfter)

			result, err = reconcilerManager.Reconcile(node, state)
			if err != nil && !result.Requeue {
				return fmt.Errorf("reconciling node (%s): %w", node.Type.String(), err)
			}
		}
	}

	return nil
}
