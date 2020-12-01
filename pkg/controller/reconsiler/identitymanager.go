package reconsiler

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// IdentityManagerResourceState contains runtime data needed in Reconsile()
type IdentityManagerResourceState struct {
	HostedZoneID string
	Domain string
}

type identityManagerReconsiler struct {
	commonMetadata *resourcetree.CommonMetadata

	client client.IdentityManagerService
}

// SetCommonMetadata saves common metadata for use in Reconsile()
func (z *identityManagerReconsiler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

/*
Reconsile knows how to do what is necessary to ensure the desired state is achieved
Requires:
- Hosted Zone
- Nameservers setup
 */
func (z *identityManagerReconsiler) Reconsile(node *resourcetree.ResourceNode) (*ReconsilationResult, error) {
	resourceState, ok := node.ResourceState.(IdentityManagerResourceState)
	if !ok {
		return nil, errors.New("unable to cast identity manager resourceState")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		authDomain := fmt.Sprintf("auth.%s", resourceState.Domain)
		authFQDN := dns.Fqdn(authDomain)
		
		_, err := z.client.CreateIdentityPool(z.commonMetadata.Ctx, api.CreateIdentityPoolOpts{
			ID:           z.commonMetadata.Id,
			AuthDomain:   authDomain,
			AuthFQDN:     authFQDN,
			HostedZoneID: resourceState.HostedZoneID,
		})
		if err != nil {
			return &ReconsilationResult{Requeue: true}, fmt.Errorf("error creating identity manager resource: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return nil, errors.New("deleting identity manager resource is not implemented")
	}

	return &ReconsilationResult{Requeue: false}, nil
}

// NewIdentityManagerReconsiler creates a new reconsiler for the Identity Manager resource
func NewIdentityManagerReconsiler(client client.IdentityManagerService) *identityManagerReconsiler {
	return &identityManagerReconsiler{
		client: client,
	}
}
