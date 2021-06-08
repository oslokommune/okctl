package reconciliation

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

type tempoReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.MonitoringService
}

// NodeType returns the relevant NodeType for this reconciler
func (z *tempoReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeTempo
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *tempoReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *tempoReconciler) Reconcile(node *dependencytree.Node, _ *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		_, err = z.client.CreateTempo(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("creating tempo: %w", err)
		}
	case dependencytree.NodeStateAbsent:
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
