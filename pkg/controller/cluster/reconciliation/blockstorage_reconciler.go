package reconciliation

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// blockstorageReconciler contains service and metadata for the relevant resource
type blockstorageReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.BlockstorageService
}

// NodeType returns the relevant NodeType for this reconciler
func (z *blockstorageReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeBlockstorage
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *blockstorageReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *blockstorageReconciler) Reconcile(node *dependencytree.Node, _ *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		_, err = z.client.CreateBlockstorage(z.commonMetadata.Ctx, client.CreateBlockstorageOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return result, fmt.Errorf("creating blockstorage: %w", err)
		}
	case dependencytree.NodeStateAbsent:
		err = z.client.DeleteBlockstorage(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting blockstorage: %w", err)
		}
	}

	return result, nil
}

// NewBlockstorageReconciler creates a new reconciler for the Blockstorage resource
func NewBlockstorageReconciler(client client.BlockstorageService) reconciliation.Reconciler {
	return &blockstorageReconciler{
		client: client,
	}
}
