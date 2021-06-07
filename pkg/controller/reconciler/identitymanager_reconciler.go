package reconciler

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/miekg/dns"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// identityManagerReconciler contains service and metadata for the relevant resource
type identityManagerReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.IdentityManagerService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *identityManagerReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeIdentityManager
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
func (z *identityManagerReconciler) Reconcile(node *resourcetree.ResourceNode, state *clientCore.StateHandlers) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		authDomain := fmt.Sprintf("auth.%s", z.commonMetadata.Declaration.ClusterRootDomain)
		authFQDN := dns.Fqdn(authDomain)

		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return result, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		_, err = z.client.CreateIdentityPool(z.commonMetadata.Ctx, client.CreateIdentityPoolOpts{
			ID:           z.commonMetadata.ClusterID,
			AuthDomain:   authDomain,
			AuthFQDN:     authFQDN,
			HostedZoneID: hz.HostedZoneID,
		})
		if err != nil {
			result.Requeue = errors.IsKind(err, errors.Timeout)

			return result, fmt.Errorf("creating identity manager resource: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err := z.client.DeleteIdentityPool(z.commonMetadata.Ctx, z.commonMetadata.ClusterID)
		if err != nil {
			return result, fmt.Errorf("deleting identity manager: %w", err)
		}

		return result, nil
	}

	return result, nil
}

// NewIdentityManagerReconciler creates a new reconciler for the Identity Manager resource
func NewIdentityManagerReconciler(client client.IdentityManagerService) Reconciler {
	return &identityManagerReconciler{
		client: client,
	}
}
