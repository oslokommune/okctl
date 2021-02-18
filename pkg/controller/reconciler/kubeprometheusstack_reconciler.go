package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

type kubePrometheusStackReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.KubePromStackService
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *kubePrometheusStackReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *kubePrometheusStackReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateKubePromStack(z.commonMetadata.Ctx, client.CreateKubePromStackOpts{ID: z.commonMetadata.ClusterID})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error creating KubePromStack: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return nil, fmt.Errorf("deletion of KubePromStack not implemented")
	}

	return &ReconcilationResult{Requeue: false}, nil
}

// NewKubePrometheusStackReconciler creates a new reconciler for the KubePromStack resource
func NewKubePrometheusStackReconciler(client client.KubePromStackService) Reconciler {
	return &kubePrometheusStackReconciler{
		client: client,
	}
}
