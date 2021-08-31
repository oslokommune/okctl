// Package hostedzone provides some functionality for interacting
// with the AWS API
package hostedzone

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// NameServersFunc defines the method signature for retrieving name servers
type NameServersFunc func(domain string) (nameservers []string, err error)

// Client contains state for interacting with required
// services
type Client struct {
	provider v1alpha1.CloudProvider
}

// New returns an initialised hostedzone client
func New(provider v1alpha1.CloudProvider) *Client {
	return &Client{
		provider: provider,
	}
}

// UndelegatedHostedZone contains information about
// a hosted zone that is present, but the NS records
// do not resolve.
type UndelegatedHostedZone struct {
	Name        string
	NameServers []string
}

// UndelegatedZonesInHostedZones will find all NS record delegations for a HostedZone and determine
// if any of them are delegated without responding to a NS record query
func (c *Client) UndelegatedZonesInHostedZones(hostedZoneID string, fn NameServersFunc) ([]*UndelegatedHostedZone, error) {
	var undelegated []*UndelegatedHostedZone

	nextRecordName := "."

	for {
		res, err := c.provider.Route53().ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
			HostedZoneId:    aws.String(hostedZoneID),
			StartRecordName: aws.String(nextRecordName),
			StartRecordType: aws.String("NS"),
		})
		if err != nil {
			return nil, fmt.Errorf(constant.GetRecordsForHostedZoneError, err)
		}

		for _, record := range res.ResourceRecordSets {
			if *record.Type != "NS" {
				continue
			}

			var expectedNameServers []string

			for _, ns := range record.ResourceRecords {
				expectedNameServers = append(expectedNameServers, *ns.Value)
			}

			gotNameServers, err := fn(*record.Name)
			if err != nil && !strings.Contains(err.Error(), "Name servers refused query (lame delegation?)") {
				return nil, fmt.Errorf(constant.GetNameServersForDNSLookup, err)
			}

			missing := compare(expectedNameServers, gotNameServers)

			if len(missing) == len(expectedNameServers) {
				undelegated = append(undelegated, &UndelegatedHostedZone{
					Name:        *record.Name,
					NameServers: expectedNameServers,
				})
			}
		}

		if !*res.IsTruncated {
			break
		}

		nextRecordName = *res.NextRecordName
	}

	return undelegated, nil
}

// compare is copied from:
// https://gist.github.com/arxdsilva/7392013cbba7a7090cbcd120b7f5ca31
func compare(a, b []string) []string {
	for i := len(a) - 1; i >= 0; i-- {
		for _, vD := range b {
			if a[i] == vD {
				a = append(a[:i], a[i+1:]...)
				break
			}
		}
	}

	return a
}
