package reconciler

import (
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// nameserverDelegationReconciler handles creation (later edit and deletion) of nameserver delegation resources.
// A nameserver delegation consists of creating a request to add a NS record to the top level domain and verifying
// that the delegation has happened.
type nameserverDelegationReconciler struct {
	commonMetadata *resourcetree.CommonMetadata
	stateHandlers  *clientCore.StateHandlers

	client client.NameserverRecordDelegationService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *nameserverDelegationReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeNameserverDelegator
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *nameserverDelegationReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// SetStateHandlers sets the state handlers
func (z *nameserverDelegationReconciler) SetStateHandlers(handlers *clientCore.StateHandlers) {
	z.stateHandlers = handlers
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *nameserverDelegationReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		hz, err := z.stateHandlers.Domain.GetPrimaryHostedZone()
		if err != nil {
			return result, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		_, err = z.client.CreateNameserverRecordDelegationRequest(&client.CreateNameserverDelegationRequestOpts{
			ClusterID:             z.commonMetadata.ClusterID,
			PrimaryHostedZoneFQDN: hz.FQDN,
			Nameservers:           hz.NameServers,
		})
		if err != nil {
			return result, fmt.Errorf("handling nameservers: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		// Should create a PR for removing the delegation.
		return result, nil
	}

	return result, nil
}

// NewNameserverDelegationReconciler creates a new reconciler for the nameserver record delegation resource
func NewNameserverDelegationReconciler(client client.NameserverRecordDelegationService) Reconciler {
	return &nameserverDelegationReconciler{
		client: client,
	}
}
