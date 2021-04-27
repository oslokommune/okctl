package client

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/config/constant"
)

const (
	numExpectedNameServers = 4
)

// InitiateDomainDelegationOpts contains required inputs for creating a
// DNS zone delegation
type InitiateDomainDelegationOpts struct {
	ClusterID api.ID

	PrimaryHostedZoneFQDN string
	Nameservers           []string
	Labels                []string
}

// Validate the provided inputs
func (o InitiateDomainDelegationOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterID, validation.Required),
		validation.Field(&o.PrimaryHostedZoneFQDN, validation.Required, is.DNSName),
		validation.Field(
			&o.Nameservers,
			validation.Required,
			validation.Each(is.DNSName),
			validation.Length(numExpectedNameServers, numExpectedNameServers),
		),
		validation.Field(&o.Labels, validation.Each(validation.In(constant.DefaultAutomaticPullRequestMergeLabel))),
	)
}

// RevokeDomainDelegationOpts contains required inputs for removing
// a DNS Zone delegation
type RevokeDomainDelegationOpts struct {
	ClusterID             api.ID
	PrimaryHostedZoneFQDN string
	Labels                []string
}

// Validate the inputs
func (o RevokeDomainDelegationOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ClusterID, validation.Required),
		validation.Field(&o.PrimaryHostedZoneFQDN, validation.Required, is.DNSName),
		validation.Field(&o.Labels, validation.Each(validation.In(constant.DefaultAutomaticPullRequestMergeLabel))),
	)
}

// NSRecordDelegationService defines required functionality for requesting a nameserver delegation record in
// the top level domain.
//
// If a team wants 'team.oslo.systems', okctl will create that domain which will get its own nameservers assigned.
// The top level domain 'oslo.systems' then needs to delegate DNS inquiries for 'team.oslo.systems' to the assigned
// nameservers. This is the delegation this service should handle.
type NSRecordDelegationService interface {
	InitiateDomainDelegation(opts InitiateDomainDelegationOpts) error
	RevokeDomainDelegation(opts RevokeDomainDelegationOpts) error
}
