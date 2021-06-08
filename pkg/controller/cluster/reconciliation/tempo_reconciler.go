package reconciliation

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/resourcetree"
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
func (z *tempoReconciler) Reconcile(node *resourcetree.ResourceNode, _ *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreateTempo(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("creating tempo: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteTempo(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting tempo: %w", err)
		}
	}

	return result, nil
}

// NewTempoReconciler creates a new reconciler for the tempo resource
func NewTempoReconciler(client client.MonitoringService) reconciliation.Reconciler {
	return &tempoReconciler{
		client: client,
	}
}
