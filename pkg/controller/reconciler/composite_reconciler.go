package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/oslokommune/okctl/pkg/spinner"
)

// compositeReconciler simplifies reconciliation by containing a collection of reconcilers and chooses the correct one
// for the provided ResourceNode based on its type
type compositeReconciler struct {
	reconcilers map[resourcetree.ResourceNodeType]Reconciler
	spinner     spinner.Spinner
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (c *compositeReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeGroup
}

// Reconcile knows what reconciler to use for the provided ResourceNode
func (c *compositeReconciler) Reconcile(node *resourcetree.ResourceNode) (result *ReconcilationResult, err error) {
	err = c.spinner.Start(resourcetree.ResourceNodeTypeToString(node.Type))
	if err != nil {
		return nil, fmt.Errorf("starting subspinner: %w", err)
	}

	defer func() {
		_ = c.spinner.Stop()
	}()

	node.RefreshState()

	return c.reconcilers[node.Type].Reconcile(node)
}

// SetCommonMetadata sets commonMetadata for all reconcilers
func (c *compositeReconciler) SetCommonMetadata(commonMetadata *resourcetree.CommonMetadata) {
	for _, reconciler := range c.reconcilers {
		reconciler.SetCommonMetadata(commonMetadata)
	}
}

// NewCompositeReconciler initializes a compositeReconciler
func NewCompositeReconciler(spin spinner.Spinner, reconcilers ...Reconciler) Reconciler {
	reconcilerMap := map[resourcetree.ResourceNodeType]Reconciler{resourcetree.ResourceNodeTypeGroup: &NoopReconciler{}}

	for _, reconciler := range reconcilers {
		reconcilerMap[reconciler.NodeType()] = reconciler
	}

	return &compositeReconciler{spinner: spin, reconcilers: reconcilerMap}
}
