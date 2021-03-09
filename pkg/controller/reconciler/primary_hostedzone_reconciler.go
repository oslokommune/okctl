package reconciler

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// zoneReconciler contains service and metadata for the relevant resource
type zoneReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.DomainService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *zoneReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeZone
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *zoneReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *zoneReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, err = z.client.CreatePrimaryHostedZoneWithoutUserinput(z.commonMetadata.Ctx, client.CreatePrimaryHostedZoneOpts{
			ID:     z.commonMetadata.ClusterID,
			Domain: z.commonMetadata.Declaration.ClusterRootURL,
			FQDN:   dns.Fqdn(z.commonMetadata.Declaration.ClusterRootURL),
		})
		if err != nil {
			return result, fmt.Errorf("creating hosted zone: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return result, errors.New("deletion of the hosted zone resource is not implemented")
	}

	return result, nil
}

// NewZoneReconciler creates a new reconciler for the Hosted Zone resource
func NewZoneReconciler(client client.DomainService) Reconciler {
	return &zoneReconciler{
		client: client,
	}
}
