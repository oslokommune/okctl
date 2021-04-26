package reconciler

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

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

	client client.NSRecordDelegationService
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
	hz, err := z.stateHandlers.Domain.GetPrimaryHostedZone()
	if err != nil {
		return result, fmt.Errorf("getting primary hosted zone: %w", err)
	}

	var labels []string

	if z.commonMetadata.Declaration.Experimental != nil {
		if z.commonMetadata.Declaration.Experimental.AutomatizeZoneDelegation {
			labels = append(labels, constant.DefaultAutomaticPullRequestMergeLabel)
		}
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		err = z.client.InitiateDomainDelegation(client.InitiateDomainDelegationOpts{
			ClusterID:             z.commonMetadata.ClusterID,
			PrimaryHostedZoneFQDN: hz.FQDN,
			Nameservers:           hz.NameServers,
			Labels:                labels,
		})
		if err != nil {
			return result, fmt.Errorf("initiating dns zone delegation: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		err = z.client.RevokeDomainDelegation(client.RevokeDomainDelegationOpts{
			ClusterID:             z.commonMetadata.ClusterID,
			PrimaryHostedZoneFQDN: hz.FQDN,
			Labels:                labels,
		})
		if err != nil {
			return result, fmt.Errorf("revoking dns zone delegation: %w", err)
		}
	}

	return result, nil
}

// NewNameserverDelegationReconciler creates a new reconciler for the nameserver record delegation resource
func NewNameserverDelegationReconciler(client client.NSRecordDelegationService) Reconciler {
	return &nameserverDelegationReconciler{
		client: client,
	}
}
