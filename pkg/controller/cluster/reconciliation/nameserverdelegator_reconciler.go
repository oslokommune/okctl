package reconciliation

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// nameserverDelegationReconciler handles creation (later edit and deletion) of nameserver delegation resources.
// A nameserver delegation consists of creating a request to add a NS record to the top level domain and verifying
// that the delegation has happened.
type nameserverDelegationReconciler struct {
	commonMetadata *reconciliation.CommonMetadata

	client client.NSRecordDelegationService
}

// NodeType returns the relevant NodeType for this reconciler
func (z *nameserverDelegationReconciler) NodeType() dependencytree.NodeType {
	return dependencytree.NodeTypeNameserverDelegator
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *nameserverDelegationReconciler) SetCommonMetadata(metadata *reconciliation.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *nameserverDelegationReconciler) Reconcile(node *dependencytree.Node, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	hz, err := state.Domain.GetPrimaryHostedZone()
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
	case dependencytree.NodeStatePresent:
		err = z.client.InitiateDomainDelegation(client.InitiateDomainDelegationOpts{
			ClusterID:             z.commonMetadata.ClusterID,
			PrimaryHostedZoneFQDN: hz.FQDN,
			Nameservers:           hz.NameServers,
			Labels:                labels,
		})
		if err != nil {
			return result, fmt.Errorf("initiating dns zone delegation: %w", err)
		}
	case dependencytree.NodeStateAbsent:
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
func NewNameserverDelegationReconciler(client client.NSRecordDelegationService) reconciliation.Reconciler {
	return &nameserverDelegationReconciler{
		client: client,
	}
}
