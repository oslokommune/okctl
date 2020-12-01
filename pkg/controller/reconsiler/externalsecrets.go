package reconsiler

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

type externalSecretsReconsiler struct {
	commonMetadata *resourcetree.CommonMetadata
	
	client client.ExternalSecretsService
}

// SetCommonMetadata saves common metadata for use in Reconsile()
func (z *externalSecretsReconsiler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconsile knows how to do what is necessary to ensure the desired state is achieved
func (z *externalSecretsReconsiler) Reconsile(node *resourcetree.ResourceNode) (*ReconsilationResult, error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateExternalSecrets(z.commonMetadata.Ctx, client.CreateExternalSecretsOpts{ID: z.commonMetadata.Id})
		
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error creating external secrets: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteExternalSecrets(z.commonMetadata.Ctx, z.commonMetadata.Id)

		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error deleting external secrets: %w", err)
		}
	}

	return &ReconsilationResult{Requeue: false}, nil
}

// NewExternalSecretsReconsiler creates a new reconsiler for the ExternalSecrets resource
func NewExternalSecretsReconsiler(client client.ExternalSecretsService) *externalSecretsReconsiler {
	return &externalSecretsReconsiler{
		client: client,
	}
}

