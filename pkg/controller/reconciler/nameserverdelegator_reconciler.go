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

const nsRecordValidationIntervalSeconds = 15

// NameserverHandlerReconcilerResourceState contains data extracted from the desired state
type NameserverHandlerReconcilerResourceState struct {
	PrimaryHostedZoneFQDN string
	Nameservers           []string
}

// nameserverDelegationReconciler handles creation (later edit and deletion) of nameserver delegation resources.
// A nameserver delegation consists of creating a request to add a NS record to the top level domain and verifying
// that the delegation has happened.
type nameserverDelegationReconciler struct {
	commonMetadata *resourcetree.CommonMetadata

	client        client.NameserverRecordDelegationService
	domainService client.DomainService
}

// NodeType returns the relevant ResourceNodeType for this reconciler
func (z *nameserverDelegationReconciler) NodeType() resourcetree.ResourceNodeType {
	return resourcetree.ResourceNodeTypeNameserverDelegator
}

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *nameserverDelegationReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *nameserverDelegationReconciler) Reconcile(node *resourcetree.ResourceNode) (result ReconcilationResult, err error) {
	resourceState, ok := node.ResourceState.(NameserverHandlerReconcilerResourceState)
	if !ok {
		return result, errors.New("casting nameserverhandler state")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		var record *client.NameserverRecord

		primaryHostedZoneFQDN := dns.Fqdn(z.commonMetadata.Declaration.PrimaryDNSZone.ParentDomain)

		record, err = z.client.CreateNameserverRecordDelegationRequest(&client.CreateNameserverDelegationRequestOpts{
			ClusterID:             z.commonMetadata.ClusterID,
			PrimaryHostedZoneFQDN: primaryHostedZoneFQDN,
			Nameservers:           resourceState.Nameservers,
		})
		if err != nil {
			return result, fmt.Errorf("handling nameservers: %w", err)
		}

		fmt.Fprintf(
			z.commonMetadata.Out,
			delegationRequestMessage,
			aurora.Green("nameserver delegation request"),
			aurora.Bold("#kjøremiljø-support"),
		)

		waitForNameserverDelegation(nsRecordValidationIntervalSeconds, primaryHostedZoneFQDN)

		err = z.domainService.SetHostedZoneDelegation(z.commonMetadata.Ctx, domain.EnsureNotFQDN(record.FQDN), true)
		if err != nil {
			return result, fmt.Errorf("setting hosted zone delegation status: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return result, errors.New("deletion of the hosted zone delegation is not implemented")
	}

	return result, nil
}

// NewNameserverDelegationReconciler creates a new reconciler for the nameserver record delegation resource
func NewNameserverDelegationReconciler(client client.NameserverRecordDelegationService, domainService client.DomainService) Reconciler {
	return &nameserverDelegationReconciler{
		client:        client,
		domainService: domainService,
	}
}

func waitForNameserverDelegation(interval time.Duration, fqdn string) {
	var err error

	for {
		err = domain.ShouldHaveNameServers(fqdn)

		if err == nil {
			break
		}

		time.Sleep(interval * time.Second)
	}
}

const delegationRequestMessage = `


A %s has been submitted. We'll process this request as soon as possible.
Let us know in %s if it takes too long


`
