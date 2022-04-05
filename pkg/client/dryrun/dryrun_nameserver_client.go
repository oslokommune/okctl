package dryrun

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type nameserverService struct {
	out io.Writer
}

func (n nameserverService) InitiateDomainDelegation(opts client.InitiateDomainDelegationOpts) error {
	fmt.Fprintf(n.out, formatCreate(fmt.Sprintf("nameserver delegation for %s", opts.PrimaryHostedZoneFQDN)))

	return nil
}

func (n nameserverService) RevokeDomainDelegation(opts client.RevokeDomainDelegationOpts) error {
	fmt.Fprintf(n.out, formatDelete(fmt.Sprintf("nameserver delegation for %s", opts.PrimaryHostedZoneFQDN)))

	return nil
}
