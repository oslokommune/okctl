package reconciliation

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

type lokiReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.MonitoringService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *lokiReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeLoki
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *lokiReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *lokiReconciler) Reconcile(node *resourcetree.ResourceNode, _ *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreateLoki(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("creating Loki: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteLoki(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting Loki: %w", err)
		}
	}

	return result, nil
}

// NewLokiReconciler creates a new reconciler for the loki resource
func NewLokiReconciler(client client.MonitoringService) reconciliation.Reconciler {
	return &lokiReconciler{
		client: client,
	}
}
