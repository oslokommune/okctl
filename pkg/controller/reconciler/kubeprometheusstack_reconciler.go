package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

type kubePrometheusStackReconciler struct {
	commonMetadata *resourcetree.CommonMetadata
	stateHandlers  *clientCore.StateHandlers

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

// SetStateHandlers sets the state handlers
func (z *kubePrometheusStackReconciler) SetStateHandlers(handlers *clientCore.StateHandlers) {
	z.stateHandlers = handlers
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *kubePrometheusStackReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		hz, err := z.stateHandlers.Domain.GetPrimaryHostedZone()
		if err != nil {
			return result, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		im, err := z.stateHandlers.IdentityManager.GetIdentityPool(
			cfn.NewStackNamer().IdentityPool(z.commonMetadata.Declaration.Metadata.Name),
		)
		if err != nil {
			return result, fmt.Errorf("getting identity pool: %w", err)
		}

		_, err = z.client.CreateKubePromStack(z.commonMetadata.Ctx, client.CreateKubePromStackOpts{
			ID:           z.commonMetadata.ClusterID,
			Domain:       z.commonMetadata.Declaration.ClusterRootDomain,
			HostedZoneID: hz.HostedZoneID,
			AuthDomain:   im.AuthDomain,
			UserPoolID:   im.UserPoolID,
		})
		if err != nil {
			return result, fmt.Errorf("creating kubepromstack: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteKubePromStack(z.commonMetadata.Ctx, client.DeleteKubePromStackOpts{
			ID:     z.commonMetadata.ClusterID,
			Domain: z.commonMetadata.Declaration.ClusterRootDomain,
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
