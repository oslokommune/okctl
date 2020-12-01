package reconsiler

import (
	"errors"
	"fmt"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// ExternalDNSResourceState contains runtime data needed in Reconsile()
type ExternalDNSResourceState struct {
	HostedZoneID string
	Domain string
}

type externalDNSReconsiler struct {
	commonMetadata *resourcetree.CommonMetadata
	
	client client.ExternalDNSService
}

// SetCommonMetadata saves common metadata for use in Reconsile()
func (z *externalDNSReconsiler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconsile knows how to ensure the desired state is achieved
func (z *externalDNSReconsiler) Reconsile(node *resourcetree.ResourceNode) (*ReconsilationResult, error) {
	resourceState, ok := node.ResourceState.(ExternalDNSResourceState)
	if !ok {
		return nil, errors.New("error casting External DNS resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err := z.client.CreateExternalDNS(z.commonMetadata.Ctx, client.CreateExternalDNSOpts{
			ID:           z.commonMetadata.ClusterId,
			HostedZoneID: resourceState.HostedZoneID,
			Domain:       resourceState.Domain,
		})
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error creating external DNS: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteExternalDNS(z.commonMetadata.Ctx, z.commonMetadata.ClusterId)
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error deleting external DNS: %w", err)
		}
	}


	return &ReconsilationResult{Requeue: false}, nil
}

// NewExternalDNSReconsiler creates a new reconsiler for the ExternalDNS resource
func NewExternalDNSReconsiler(client client.ExternalDNSService) *externalDNSReconsiler {
	return &externalDNSReconsiler{
		client: client,
	}
}
