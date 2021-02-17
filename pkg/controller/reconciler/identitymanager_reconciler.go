package reconciler

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// IdentityManagerResourceState contains runtime data needed in Reconcile()
type IdentityManagerResourceState struct {
	HostedZoneID string
}

// identityManagerReconciler contains service and metadata for the relevant resource
type identityManagerReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.IdentityManagerService
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *identityManagerReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

/*
Reconcile knows how to do what is necessary to ensure the desired state is achieved
Requires:
- Hosted Zone
- Nameservers setup
*/
func (z *identityManagerReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	resourceState, ok := node.ResourceState.(IdentityManagerResourceState)
	if !ok {
		return nil, errors.New("casting identity manager resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		authDomain := fmt.Sprintf("auth.%s", z.commonMetadata.Declaration.PrimaryDNSZone.ParentDomain)
		authFQDN := dns.Fqdn(authDomain)

		_, err := z.client.CreateIdentityPool(z.commonMetadata.Ctx, api.CreateIdentityPoolOpts{
			ID:           z.commonMetadata.ClusterID,
			AuthDomain:   authDomain,
			AuthFQDN:     authFQDN,
			HostedZoneID: resourceState.HostedZoneID,
		})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error creating identity manager resource: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return nil, errors.New("deleting identity manager resource is not implemented")
	}

	return &ReconcilationResult{Requeue: false}, nil
}

// NewIdentityManagerReconciler creates a new reconciler for the Identity Manager resource
func NewIdentityManagerReconciler(client client.IdentityManagerService) Reconciler {
	return &identityManagerReconciler{
		client: client,
	}
}
