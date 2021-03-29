// Package mock contains some convenience functions for getting data
package mock

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// nolint: golint
const (
	DefaultRegion         = "eu-west-1"
	DefaultAWSAccountID   = "123456789012"
	DefaultEnvironment    = "staging"
	DefaultRepository     = "okctl"
	DefaultClusterName    = "okctl-staging"
	DefaultDomain         = "okctl-staging.oslo.systems"
	DefaultFQDN           = "okctl-staging.oslo.systems."
	DefaultHostedZoneID   = "Z0FAKE41FAKE6I841FAKE"
	DefaultCertificateARN = "arn:aws:acm:eu-west-1:123456789012:certificate/123456789012-1234-1234-1234-12345678"

	StackNameHostedZone  = "okctl-staging-oslo-systems-HostedZone"
	StackNameCertificate = "okctl-staging-oslo-systems-Certificate"
)

// ID returns a fake id
func ID() api.ID {
	return api.ID{
		Region:       DefaultRegion,
		AWSAccountID: DefaultAWSAccountID,
		Environment:  DefaultEnvironment,
		Repository:   DefaultRepository,
		ClusterName:  DefaultClusterName,
	}
}

// CloudFormationTemplate just returns some bytes
func CloudFormationTemplate() []byte {
	return []byte(`AWSTemplateFormatVersion: 2010-09-09
`)
}

// NameServers returns a set of fake name servers
func NameServers() []string {
	return []string{
		"ns-158-fake.awsdns-19.com.",
		"ns-1389-fake.awsdns-45.org.",
		"ns-622-fake.awsdns-13.net.",
		"ns-1614-fake.awsdns-09.co.uk.",
	}
}

// HostedZone returns a fake hosted zone
func HostedZone() *client.HostedZone {
	return &client.HostedZone{
		ID:                     ID(),
		IsDelegated:            true,
		Primary:                true,
		Managed:                true,
		FQDN:                   DefaultFQDN,
		Domain:                 DefaultDomain,
		HostedZoneID:           DefaultHostedZoneID,
		NameServers:            NameServers(),
		StackName:              StackNameHostedZone,
		CloudFormationTemplate: CloudFormationTemplate(),
	}
}

// Certificate returns a fake certificate
func Certificate() *client.Certificate {
	return &client.Certificate{
		ID:                     ID(),
		FQDN:                   DefaultFQDN,
		Domain:                 DefaultDomain,
		HostedZoneID:           DefaultHostedZoneID,
		ARN:                    DefaultCertificateARN,
		StackName:              StackNameCertificate,
		CloudFormationTemplate: CloudFormationTemplate(),
	}
}
