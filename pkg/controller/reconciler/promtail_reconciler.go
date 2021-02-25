package reconciler

import (
	"fmt"

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
func (z *promtailReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreatePromtail(z.commonMetadata.Ctx, client.CreatePromtailOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return result, fmt.Errorf("creating promtail: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeletePromtail(z.commonMetadata.Ctx, client.DeletePromtailOpts{
			ID: z.commonMetadata.ClusterID,
		})
		if err != nil {
			return result, fmt.Errorf("deleting promtail: %w", err)
		}
	}

	return result, nil
}

// NewPromtailReconciler creates a new reconciler for the Promtail resource
func NewPromtailReconciler(client client.MonitoringService) Reconciler {
	return &promtailReconciler{
		client: client,
	}
}
