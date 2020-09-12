package console

import (
	"fmt"
	"io"

	"github.com/logrusorgru/aurora/v3"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/theckman/yacspin"
)

type domainReport struct {
	console *Console
}

func (r *domainReport) ReportCreatePrimaryHostedZone(zone *client.HostedZone, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := fmt.Sprintf("%s (%s)", aurora.Green(zone.HostedZone.StackName), zone.HostedZone.HostedZoneID)

	return r.console.Report(actions, "primary-hosted-zone", description)
}

// NewDomainReport returns an initialised domain reporter
func NewDomainReport(out io.Writer, exit chan struct{}, spinner *yacspin.Spinner) client.DomainReport {
	return &domainReport{
		console: New(out, exit, spinner),
	}
}
