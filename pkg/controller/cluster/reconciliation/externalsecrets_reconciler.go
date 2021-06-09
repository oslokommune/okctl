package reconciliation

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// externalSecretsReconciler contains service and metadata for the relevant resource
type externalSecretsReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.ExternalSecretsService
}

// NodeType returns the relevant NodeType for this reconciler
func (z *externalSecretsReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeExternalSecrets
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *externalSecretsReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *externalSecretsReconciler) Reconcile(node *dependencytree.Node, _ *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		_, err = z.client.CreateExternalSecrets(z.commonMetadata.Ctx, client.CreateExternalSecretsOpts{
			ID: z.commonMetadata.ClusterID,
		})
		if err != nil {
			return result, fmt.Errorf("creating external secrets: %w", err)
		}
	case dependencytree.NodeStateAbsent:
		err = z.client.DeleteExternalSecrets(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting external secrets: %w", err)
		}
	}

	return result, nil
}

// NewExternalSecretsReconciler creates a new reconciler for the ExternalSecrets resource
func NewExternalSecretsReconciler(client client.ExternalSecretsService) reconciliation.Reconciler {
	return &externalSecretsReconciler{
		client: client,
	}
}
