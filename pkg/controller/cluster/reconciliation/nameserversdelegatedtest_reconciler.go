package reconciliation

import (
	"fmt"
	"time"

	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/oslokommune/okctl/pkg/domain"
)

const defaultTestingIntervalMinutes = 5 * time.Minute

type nameserversDelegatedTestReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	domainService client.DomainService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (n *nameserversDelegatedTestReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeNameserversDelegatedTest
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (n *nameserversDelegatedTestReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	n.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (n *nameserversDelegatedTestReconciler) Reconcile(node *resourcetree.ResourceNode, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		_, _ = fmt.Fprintf(
			n.commonMetadata.Out,
			delegationRequestMessage,
			aurora.Green("nameserver delegation request"),
			aurora.Bold("#kjøremiljø-support"),
		)

		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return result, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		err = domain.ShouldHaveNameServers(hz.FQDN, hz.NameServers)
		if err != nil {
			result.Requeue = true
			result.RequeueAfter = defaultTestingIntervalMinutes

			_, _ = fmt.Fprintf(n.commonMetadata.Out, "failed to validate nameservers: %s", err)

			return result, fmt.Errorf("validating nameservers: %w", err)
		}

		err = n.domainService.SetHostedZoneDelegation(
			n.commonMetadata.Ctx,
			n.commonMetadata.Declaration.ClusterRootDomain,
			true,
		)
		if err != nil {
			return result, fmt.Errorf("setting hosted zone delegation status: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		// Nothing to do on absent
		return result, nil
	}

	return result, nil
}

// NewNameserverDelegatedTestReconciler creates a new reconciler for the nameserver record delegation test resource
func NewNameserverDelegatedTestReconciler(domainService client.DomainService) reconciliation.Reconciler {
	return &nameserversDelegatedTestReconciler{domainService: domainService}
}

const delegationRequestMessage = `


A %s has been submitted. We'll process this request as soon as possible.
Let us know in %s if it takes too long


`
