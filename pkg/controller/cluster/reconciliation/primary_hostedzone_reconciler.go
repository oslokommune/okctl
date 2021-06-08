package reconciliation

import (
	"errors"
	"fmt"

	"github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/miekg/dns"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// zoneReconciler contains service and metadata for the relevant resource
type zoneReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.DomainService
}

// NodeType returns the relevant NodeType for this reconciler
func (z *zoneReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeZone
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *zoneReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *zoneReconciler) Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case dependencytree.NodeStatePresent:
		_, err = z.client.CreatePrimaryHostedZone(z.commonMetadata.Ctx, client.CreatePrimaryHostedZoneOpts{
			ID:     z.commonMetadata.ClusterID,
			Domain: z.commonMetadata.Declaration.ClusterRootDomain,
			FQDN:   dns.Fqdn(z.commonMetadata.Declaration.ClusterRootDomain),
		})
		if err != nil {
			return result, fmt.Errorf("creating hosted zone: %w", err)
		}
	case dependencytree.NodeStateAbsent:
		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			// Already removed, moving on
			if errors.Is(err, storm.ErrNotFound) {
				return result, nil
			}

			return result, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		err = z.client.DeletePrimaryHostedZone(z.commonMetadata.Ctx, client.DeletePrimaryHostedZoneOpts{
			ID:           z.commonMetadata.ClusterID,
			HostedZoneID: hz.HostedZoneID,
		})
		if err != nil {
			return result, fmt.Errorf("deleting primary hosted zone: %w", err)
		}

		return result, nil
	}

	return result, nil
}

// NewZoneReconciler creates a new reconciler for the Hosted Zone resource
func NewZoneReconciler(client client.DomainService) reconciliation.Reconciler {
	return &zoneReconciler{
		client: client,
	}
}
