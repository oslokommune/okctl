package reconciliation

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

type lokiReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.MonitoringService
}

// NodeType returns the relevant NodeType for this reconciler
func (z *lokiReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeLoki
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *lokiReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *lokiReconciler) Reconcile(node *dependencytree.Node, _ *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		_, err = z.client.CreateLoki(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("creating Loki: %w", err)
		}
	case dependencytree.NodeStateAbsent:
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
