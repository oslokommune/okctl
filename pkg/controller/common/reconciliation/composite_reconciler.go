package reconciliation

import (
	"fmt"
	"strings"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

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
func (c *compositeReconciler) Reconcile(node *resourcetree.ResourceNode, state *clientCore.StateHandlers) (result Result, err error) {
	err = c.spinner.Start(node.Type.String())
	if err != nil {
		return result, fmt.Errorf("starting subspinner: %w", err)
	}

	defer func() {
		_ = c.spinner.Stop()
	}()

	t := node.Type

	if strings.HasPrefix(t.String(), resourcetree.ResourceNodeTypePostgresInstance.String()) {
		t = resourcetree.ResourceNodeTypePostgresInstance
	}

	_, ok := c.reconcilers[t]
	if !ok {
		return result, fmt.Errorf("no reconciler for type exists: %s", node.Type.String())
	}

	return c.reconcilers[t].Reconcile(node, state)
}

// SetCommonMetadata sets commonMetadata for all reconcilers
func (c *compositeReconciler) SetCommonMetadata(commonMetadata *resourcetree.CommonMetadata) {
	for _, reconciler := range c.reconcilers {
		reconciler.SetCommonMetadata(commonMetadata)
	}
}

// NewCompositeReconciler initializes a compositeReconciler
func NewCompositeReconciler(spin spinner.Spinner, reconcilers ...Reconciler) Reconciler {
	reconcilerMap := map[resourcetree.ResourceNodeType]Reconciler{
		resourcetree.ResourceNodeTypeGroup: &NoopReconciler{},
	}

	for _, reconciler := range reconcilers {
		reconcilerMap[reconciler.NodeType()] = reconciler
	}

	return &compositeReconciler{spinner: spin, reconcilers: reconcilerMap}
}
