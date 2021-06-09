package reconciliation

import (
	"fmt"
	"strings"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
	"github.com/oslokommune/okctl/pkg/spinner"
)

// compositeReconciler simplifies reconciliation by containing a collection of reconcilers and chooses the correct one
// for the provided Node based on its type
type compositeReconciler struct {
	reconcilers map[dependencytree.NodeType]Reconciler
	spinner     spinner.Spinner
}

// NodeType returns the relevant NodeType for this reconciler
func (c *compositeReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeGroup
}

// Reconcile knows what reconciler to use for the provided Node
func (c *compositeReconciler) Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (result Result, err error) {
	err = c.spinner.Start(node.Type.String())
	if err != nil {
		return result, fmt.Errorf("starting subspinner: %w", err)
	}

	defer func() {
		_ = c.spinner.Stop()
	}()

	t := node.Type

	if strings.HasPrefix(t.String(), dependencytree.NodeTypePostgresInstance.String()) {
		t = dependencytree.NodeTypePostgresInstance
	}

	_, ok := c.reconcilers[t]
	if !ok {
		return result, fmt.Errorf("no reconciler for type exists: %s", node.Type.String())
	}

	return c.reconcilers[t].Reconcile(node, state)
}

// SetCommonMetadata sets commonMetadata for all reconcilers
func (c *compositeReconciler) SetCommonMetadata(commonMetadata *CommonMetadata) {
	for _, reconciler := range c.reconcilers {
		reconciler.SetCommonMetadata(commonMetadata)
	}
}

// NewCompositeReconciler initializes a compositeReconciler
func NewCompositeReconciler(spin spinner.Spinner, reconcilers ...Reconciler) Reconciler {
	reconcilerMap := map[dependencytree.NodeType]Reconciler{
		dependencytree.NodeTypeGroup: &NoopReconciler{},
	}

	for _, reconciler := range reconcilers {
		reconcilerMap[reconciler.NodeType()] = reconciler
	}

	return &compositeReconciler{spinner: spin, reconcilers: reconcilerMap}
}
