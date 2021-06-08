package reconciliation

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

type kubePrometheusStackReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.MonitoringService
}

// NodeType returns the resource node type
func (z *kubePrometheusStackReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeKubePromStack
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *kubePrometheusStackReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *kubePrometheusStackReconciler) Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return result, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		im, err := state.IdentityManager.GetIdentityPool(
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
	case dependencytree.NodeStateAbsent:
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
func NewKubePrometheusStackReconciler(client client.MonitoringService) reconciliation.Reconciler {
	return &kubePrometheusStackReconciler{
		client: client,
	}
}
