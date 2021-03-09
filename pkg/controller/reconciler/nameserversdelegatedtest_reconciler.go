package reconciler

import (
	"fmt"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/miekg/dns"
	"github.com/mishudark/errors"
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
func (n *nameserversDelegatedTestReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		fmt.Fprintf(
			n.commonMetadata.Out,
			delegationRequestMessage,
			aurora.Green("nameserver delegation request"),
			aurora.Bold("#kjøremiljø-support"),
		)

		primaryHostedZoneFQDN := dns.Fqdn(n.commonMetadata.Declaration.ClusterRootURL)

		err = domain.ShouldHaveNameServers(primaryHostedZoneFQDN)
		if err != nil {
			result.Requeue = true
			result.RequeueAfter = defaultTestingIntervalMinutes

			return result, nil
		}

		err = n.domainService.SetHostedZoneDelegation(
			n.commonMetadata.Ctx,
			n.commonMetadata.Declaration.ClusterRootURL,
			true,
		)
		if err != nil {
			return result, fmt.Errorf("setting hosted zone delegation status: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return result, errors.New("removing nameservers delegated test is not implemented")
	}

	return result, nil
}

// NewNameserverDelegatedTestReconciler creates a new reconciler for the nameserver record delegation test resource
func NewNameserverDelegatedTestReconciler(domainService client.DomainService) Reconciler {
	return &nameserversDelegatedTestReconciler{domainService: domainService}
}

const delegationRequestMessage = `


A %s has been submitted. We'll process this request as soon as possible.
Let us know in %s if it takes too long


`
