package reconciler

import (
	"fmt"
	"github.com/miekg/dns"
	"time"

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

// SetCommonMetadata saves common metadata for use in Reconcile()
func (z *nameserverDelegationReconciler) SetCommonMetadata(metadata *resourcetree.CommonMetadata) {
	z.commonMetadata = metadata
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *nameserverDelegationReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	resourceState, ok := node.ResourceState.(NameserverHandlerReconcilerResourceState)
	if !ok {
		return nil, errors.New("error casting nameserverhandler state")
	}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		err := z.commonMetadata.Spin.Start("Nameserver delegation")
		if err != nil {
			return nil, err
		}

		defer func() {
			err = z.commonMetadata.Spin.Stop()
		}()

		primaryHostedZoneFQDN := dns.Fqdn(z.commonMetadata.Declaration.PrimaryDNSZone.ParentDomain)
		record, err := z.client.CreateNameserverRecordDelegationRequest(&client.CreateNameserverDelegationRequestOpts{
			ClusterID:             z.commonMetadata.ClusterID,
			PrimaryHostedZoneFQDN: primaryHostedZoneFQDN,
			Nameservers:           resourceState.Nameservers,
		})
		if err != nil {
			return &ReconcilationResult{Requeue: true}, fmt.Errorf("error handling nameservers: %w", err)
		}

		fmt.Fprint(z.commonMetadata.Out, delegationRequestMessage)

		waitForNameserverDelegation(nsRecordValidationIntervalSeconds, primaryHostedZoneFQDN)

		err = z.domainService.SetHostedZoneDelegation(z.commonMetadata.Ctx, domain.EnsureNotFQDN(record.FQDN), true)
		if err != nil {
			return nil, fmt.Errorf("error setting hosted zone delegation status: %w", err)
		}
	case resourcetree.ResourceNodeStateAbsent:
		return nil, errors.New("deletion of the hosted zone delegation is not implemented")
	}

	return &ReconcilationResult{Requeue: false}, nil
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
A nameserver delegation request has been submitted. We'll process this request as soon as possible.
Let us know in #kjøremiljø-support if it takes too long


`
