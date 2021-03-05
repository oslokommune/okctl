package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

type tempoReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.MonitoringService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *tempoReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeTempo
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *tempoReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *tempoReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreateTempo(z.commonMetadata.Ctx, client.CreateTempoOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return result, fmt.Errorf("creating tempo: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteTempo(z.commonMetadata.Ctx, client.DeleteTempoOpts{
			ID: z.commonMetadata.ClusterID,
		})
		if err != nil {
			return result, fmt.Errorf("deleting tempo: %w", err)
		}
	}

	return result, nil
}

// NewTempoReconciler creates a new reconciler for the tempo resource
func NewTempoReconciler(client client.MonitoringService) Reconciler {
	return &tempoReconciler{
		client: client,
	}
}
