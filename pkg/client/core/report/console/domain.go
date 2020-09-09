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

func (r *domainReport) ReportCreatePrimaryHostedZone(zone *client.HostedZone, report *store.Report) error {
	description := fmt.Sprintf("%s (%s)", aurora.Green(zone.HostedZone.StackName), zone.HostedZone.HostedZoneID)
	return r.console.Report(report.Actions, "primary-hosted-zone", description)
}

// NewDomainReport returns an initialised domain reporter
func NewDomainReport(out io.Writer, exit chan struct{}, spinner *yacspin.Spinner) client.DomainReport {
	return &domainReport{
		console: New(out, exit, spinner),
	}
}
