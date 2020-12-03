package reconciler

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

type zoneReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.DomainService
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *zoneReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *zoneReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	metadata, ok := node.Metadata.(HostedZoneMetadata)
	if !ok {
		return nil, errors.New("error casting HostedZone metadata")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		fqdn := dns.Fqdn(metadata.Domain)

		_, err := z.client.CreatePrimaryHostedZoneWithoutUserinput(z.commonMetadata.Ctx, client.CreatePrimaryHostedZoneOpts{
			ID:     z.commonMetadata.ClusterId,
			Domain: metadata.Domain,
			FQDN:   fqdn,
		})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error creating hosted zone: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeletePrimaryHostedZone(z.commonMetadata.Ctx, client.DeletePrimaryHostedZoneOpts{ID: z.commonMetadata.ClusterId})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error deleting hosted zone: %w", err)
		}
	}

	return &ReconcilationResult{Requeue: false}, nil
}

// NewZoneReconciler creates a new reconciler for the Hosted Zone resource
func NewZoneReconciler(client client.DomainService) *zoneReconciler {
	return &zoneReconciler{
		client: client,
	}
}
