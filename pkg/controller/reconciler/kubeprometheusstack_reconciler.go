package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/state"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// KubePromStackState contains runtime data needed in Reconcile()
type KubePromStackState struct {
	HostedZone *state.HostedZone

	UserPoolID string
	AuthDomain string
}

type kubePrometheusStackReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.MonitoringService
}

// NodeType returns the resource node type
func (z *kubePrometheusStackReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeKubePromStack
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *kubePrometheusStackReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *kubePrometheusStackReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	resourceState, ok := node.ResourceState.(KubePromStackState)
	if !ok {
		return ReconcilationResult{}, fmt.Errorf("casting KubePromStackState resource resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreateKubePromStack(z.commonMetadata.Ctx, client.CreateKubePromStackOpts{
			ID:           z.commonMetadata.ClusterID,
			Domain:       z.commonMetadata.Declaration.ClusterRootURL,
			HostedZoneID: resourceState.HostedZone.ID,
			AuthDomain:   resourceState.AuthDomain,
			UserPoolID:   resourceState.UserPoolID,
		})
		if err != nil {
			return result, fmt.Errorf("creating kubepromstack: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteKubePromStack(z.commonMetadata.Ctx, client.DeleteKubePromStackOpts{
			ID:     z.commonMetadata.ClusterID,
			Domain: z.commonMetadata.Declaration.ClusterRootURL,
		})
		if err != nil {
			return result, fmt.Errorf("deleting kubepromstack: %w", err)
		}
	}

	return result, nil
}

// NewKubePrometheusStackReconciler creates a new reconciler for the KubePromStack resource
func NewKubePrometheusStackReconciler(client client.MonitoringService) Reconciler {
	return &kubePrometheusStackReconciler{
		client: client,
	}
}
