package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// blockstorageReconciler contains service and metadata for the relevant resource
type blockstorageReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.BlockstorageService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *blockstorageReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeBlockstorage
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *blockstorageReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *blockstorageReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateBlockstorage(z.commonMetadata.Ctx, client.CreateBlockstorageOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("creating blockstorage: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteBlockstorage(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("deleting blockstorage: %w", err)
		}
	}

	return &ReconcilationResult{Requeue: false}, nil
}

// NewBlockstorageReconciler creates a new reconciler for the Blockstorage resource
func NewBlockstorageReconciler(client client.BlockstorageService) Reconciler {
	return &blockstorageReconciler{
		client: client,
	}
}
