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
func (z *identityManagerReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	resourceState, ok := node.ResourceState.(IdentityManagerResourceState)
	if !ok {
		return result, errors.New("casting identity manager resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		authDomain := fmt.Sprintf("auth.%s", z.commonMetadata.Declaration.ClusterRootURL)
		authFQDN := dns.Fqdn(authDomain)

		_, err = z.client.CreateIdentityPool(z.commonMetadata.Ctx, api.CreateIdentityPoolOpts{
			ID:           z.commonMetadata.ClusterID,
			AuthDomain:   authDomain,
			AuthFQDN:     authFQDN,
			HostedZoneID: resourceState.HostedZoneID,
		})
		if err != nil {
			return result, fmt.Errorf("creating identity manager resource: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return result, errors.New("deleting identity manager resource is not implemented")
	}

	return result, nil
}

// NewIdentityManagerReconciler creates a new reconciler for the Identity Manager resource
func NewIdentityManagerReconciler(client client.IdentityManagerService) Reconciler {
	return &identityManagerReconciler{
		client: client,
	}
}
