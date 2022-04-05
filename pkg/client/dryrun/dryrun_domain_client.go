package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type domainClient struct {
	out io.Writer
}

func (d domainClient) CreatePrimaryHostedZone(_ context.Context, opts client.CreatePrimaryHostedZoneOpts) (*client.HostedZone, error) {
	fmt.Fprintf(d.out, formatCreate(fmt.Sprintf("%s %s", "primary hosted zone", opts.Domain)))

	return &client.HostedZone{
		ID:           opts.ID,
		Primary:      true,
		FQDN:         opts.Domain,
		Domain:       opts.Domain,
		HostedZoneID: toBeGenerated,
		NameServers: []string{
			toBeGenerated,
			toBeGenerated,
			toBeGenerated,
			toBeGenerated,
		},
		StackName:              toBeGenerated,
		CloudFormationTemplate: []byte(toBeGenerated),
	}, nil
}

func (d domainClient) GetPrimaryHostedZone(_ context.Context) (*client.HostedZone, error) {
	return &client.HostedZone{
		ID:                     api.ID{},
		IsDelegated:            false,
		Primary:                true,
		FQDN:                   "N/A",
		Domain:                 "N/A",
		HostedZoneID:           toBeGenerated,
		NameServers:            []string{toBeGenerated},
		StackName:              toBeGenerated,
		CloudFormationTemplate: []byte(toBeGenerated),
	}, nil
}

func (d domainClient) DeletePrimaryHostedZone(_ context.Context, _ client.DeletePrimaryHostedZoneOpts) error {
	fmt.Fprintf(d.out, formatDelete("primary hosted zone"))

	return nil
}

func (d domainClient) SetHostedZoneDelegation(_ context.Context, _ string, _ bool) error {
	return nil
}
