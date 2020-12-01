package reconsiler

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// HostedZoneMetadata contains data extracted from the desired state
type HostedZoneMetadata struct {
	Domain string
}

type zoneReconsiler struct {
	commonMetadata *resourcetree.CommonMetadata
	
	client client.DomainService
}

// SetCommonMetadata saves common metadata for use in Reconsile()
func (z *zoneReconsiler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconsile knows how to do what is necessary to ensure the desired state is achieved
func (z *zoneReconsiler) Reconsile(node *resourcetree.ResourceNode) (*ReconsilationResult, error) {
	metadata, ok := node.Metadata.(HostedZoneMetadata)
	if !ok {
		return nil, errors.New("error casting HostedZone metadata")
	}
	
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		fqdn := dns.Fqdn(metadata.Domain)

		_, err := z.client.CreatePrimaryHostedZoneWithoutUserinput(z.commonMetadata.Ctx, client.CreatePrimaryHostedZoneOpts{
			ID:     z.commonMetadata.Id,
			Domain: metadata.Domain,
			FQDN: fqdn,
		})
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error creating hosted zone: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeletePrimaryHostedZone(z.commonMetadata.Ctx, client.DeletePrimaryHostedZoneOpts{ID: z.commonMetadata.Id})
		if err != nil {
		    return &ReconsilationResult{Requeue: true}, fmt.Errorf("error deleting hosted zone: %w", err)
		}
	}

	return &ReconsilationResult{Requeue: false}, nil
}

// NewZoneReconsiler creates a new reconsiler for the Hosted Zone resource
func NewZoneReconsiler(client client.DomainService) *zoneReconsiler {
	return &zoneReconsiler{
		client: client,
	}
}

