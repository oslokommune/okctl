package route53

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"

	route "github.com/aws/aws-sdk-go/service/route53"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// Route53er defines the available methods
type Route53er interface {
	HostedZones() ([]HostedZone, error)
}

// HostedZone contains information about an AWS Route53 HostedZone
type HostedZone struct {
	ID     string
	Domain string
	FQDN   string
	Public bool
}

// Route53 contains the state required for querying
// the Route53 API
type Route53 struct {
	provider v1alpha1.CloudProvider
}

// New returns an initialised route53 client
func New(provider v1alpha1.CloudProvider) *Route53 {
	return &Route53{
		provider: provider,
	}
}

// PublicHostedZones returns all the registered hosted zones
func (r *Route53) PublicHostedZones() ([]*HostedZone, error) {
	out, err := r.provider.Route53().ListHostedZones(&route.ListHostedZonesInput{})
	if err != nil {
		return nil, err
	}

	var zones []*HostedZone

	for _, hz := range out.HostedZones {
		if !*hz.Config.PrivateZone {
			if !*hz.Config.PrivateZone {
				zones = append(zones, &HostedZone{
					ID:     strings.TrimPrefix(*hz.Id, "/hostedzone/"),
					Domain: strings.TrimSuffix(*hz.Name, "."),
					FQDN:   *hz.Name,
					Public: true,
				})
			}
		}
	}

	return zones, nil
}

// DeleteHostedZoneRecordSets removes all recods in a hosted zone except NS and SOA records, to enable zone deletion
func (r *Route53) DeleteHostedZoneRecordSets(hostedZoneID string) (*route.ChangeResourceRecordSetsOutput, error) {
	ListResourceRecordSetsInput := route.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
	}

	got, err := r.provider.Route53().ListResourceRecordSets(&ListResourceRecordSetsInput)
	if err != nil {
		return nil, err
	}

	sets := got.ResourceRecordSets

	var changes []*route.Change

	for _, set := range sets {
		if *set.Type != "SOA" && *set.Type != "NS" {
			changes = append(changes, &route.Change{
				Action: aws.String("DELETE"),
				ResourceRecordSet: &route.ResourceRecordSet{
					AliasTarget:             set.AliasTarget,
					Name:                    set.Name,
					TrafficPolicyInstanceId: set.TrafficPolicyInstanceId,
					ResourceRecords:         set.ResourceRecords,
					SetIdentifier:           set.SetIdentifier,
					TTL:                     set.TTL,
					Type:                    set.Type,
					Weight:                  set.Weight,
				},
			})
		}
	}

	if len(changes) == 0 {
		return nil, nil
	}

	request := &route.ChangeResourceRecordSetsInput{
		ChangeBatch: &route.ChangeBatch{
			Changes: changes,
		},
		HostedZoneId: aws.String(hostedZoneID),
	}

	resp, err := r.provider.Route53().ChangeResourceRecordSets(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetNSRecordTTL finds the NS record in a hosted zone, and sets the TTL for it. If there are zere or more than one
// NS records, an error is thrown.
func (r *Route53) SetNSRecordTTL(hostedZoneID string, nsTTLSecs int64) error {
	ListResourceRecordSetsInput := route.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
	}

	got, err := r.provider.Route53().ListResourceRecordSets(&ListResourceRecordSetsInput)
	if err != nil {
		return err
	}

	sets := got.ResourceRecordSets
	nsCount := 0

	var nsSet *route.ResourceRecordSet

	for _, set := range sets {
		if *set.Type == "NS" {
			nsCount++

			nsSet = set
		}
	}

	if nsCount != 1 {
		return errors.Errorf("expected 1 NS records in hosted zone, but found %d", nsCount)
	}

	var changes []*route.Change
	changes = append(changes, &route.Change{
		Action: aws.String("UPSERT"),
		ResourceRecordSet: &route.ResourceRecordSet{
			Name:            nsSet.Name,
			ResourceRecords: nsSet.ResourceRecords,
			TTL:             &nsTTLSecs,
			Type:            nsSet.Type,
		},
	})

	comment := "Update NS record TTL"

	request := &route.ChangeResourceRecordSetsInput{
		ChangeBatch: &route.ChangeBatch{
			Changes: changes,
			Comment: &comment,
		},
		HostedZoneId: aws.String(hostedZoneID),
	}

	_, err = r.provider.Route53().ChangeResourceRecordSets(request)
	if err != nil {
		return err
	}

	return nil
}
