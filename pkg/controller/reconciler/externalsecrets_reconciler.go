package reconciler

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// externalSecretsReconciler contains service and metadata for the relevant resource
type externalSecretsReconciler struct {
	commonMetadata *resourcetree.CommonMetadata
	stateHandlers  *clientCore.StateHandlers

	client client.ExternalSecretsService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *externalSecretsReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeExternalSecrets
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *externalSecretsReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// SetStateHandlers sets the state handlers
func (z *externalSecretsReconciler) SetStateHandlers(handlers *clientCore.StateHandlers) {
	z.stateHandlers = handlers
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *externalSecretsReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreateExternalSecrets(z.commonMetadata.Ctx, client.CreateExternalSecretsOpts{
			ID: z.commonMetadata.ClusterID,
		})
		if err != nil {
			return result, fmt.Errorf("creating external secrets: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.DeleteExternalSecrets(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting external secrets: %w", err)
		}
	}

	return result, nil
}

// NewExternalSecretsReconciler creates a new reconciler for the ExternalSecrets resource
func NewExternalSecretsReconciler(client client.ExternalSecretsService) Reconciler {
	return &externalSecretsReconciler{
		client: client,
	}
}
