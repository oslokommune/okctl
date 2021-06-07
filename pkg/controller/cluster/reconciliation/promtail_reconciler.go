package reconciliation

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

type promtailReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.MonitoringService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *promtailReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypePromtail
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *promtailReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *promtailReconciler) Reconcile(node *resourcetree.ResourceNode, _ *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreatePromtail(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("creating promtail: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeletePromtail(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting promtail: %w", err)
		}
	}

	return result, nil
}

// NewPromtailReconciler creates a new reconciler for the Promtail resource
func NewPromtailReconciler(client client.MonitoringService) reconciliation.Reconciler {
	return &promtailReconciler{
		client: client,
	}
}
