package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

type externalSecretsReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.ExternalSecretsService
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *externalSecretsReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *externalSecretsReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateExternalSecrets(z.commonMetadata.Ctx, client.CreateExternalSecretsOpts{ID: z.commonMetadata.ClusterId})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error creating external secrets: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteExternalSecrets(z.commonMetadata.Ctx, z.commonMetadata.ClusterId)
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error deleting external secrets: %w", err)
		}
	}

	return &ReconcilationResult{Requeue: false}, nil
}

// NewExternalSecretsReconciler creates a new reconciler for the ExternalSecrets resource
func NewExternalSecretsReconciler(client client.ExternalSecretsService) *externalSecretsReconciler {
	return &externalSecretsReconciler{
		client: client,
	}
}
