package rest

import (
	route "github.com/aws/aws-sdk-go/service/route53"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/route53"
)

// TargetHostedZone matches the REST API route
const TargetHostedZone = "domains/hostedzones/"

type domainAPI struct {
	client *HTTPClient
}

func (a *domainAPI) DeleteHostedZoneRecords(provider v1alpha1.CloudProvider, hostedZoneID string) (*route.ChangeResourceRecordSetsOutput, error) {
	return route53.New(provider).DeleteHostedZoneRecordSets(hostedZoneID)
}

func (a *domainAPI) DeletePrimaryHostedZone(domain string, opts client.DeletePrimaryHostedZoneOpts) error {
	err := a.client.DoDelete(TargetHostedZone, &api.DeleteHostedZoneOpts{
		ID:           opts.ID,
		Domain:       domain,
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *domainAPI) CreatePrimaryHostedZone(opts client.CreatePrimaryHostedZoneOpts) (*client.HostedZone, error) {
	into := &api.HostedZone{}

	err := a.client.DoPost(TargetHostedZone, &api.CreateHostedZoneOpts{
		ID:     opts.ID,
		Domain: opts.Domain,
		FQDN:   opts.FQDN,
		NSTTL:  opts.NSTTL,
	}, into)
	if err != nil {
		return nil, err
	}

	return &client.HostedZone{
		IsDelegated: false,
		Primary:     true,
		HostedZone:  into,
	}, nil
}

// NewDomainAPI returns an initialised REST API client
func NewDomainAPI(client *HTTPClient) client.DomainAPI {
	return &domainAPI{
		client: client,
	}
}
