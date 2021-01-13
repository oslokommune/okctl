package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora/v3"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type domainReport struct {
	console *Console
}

func (r *domainReport) ReportDeletePrimaryHostedZone(reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	return r.console.Report(actions, "primary-hosted-zone", aurora.Green("removing").String())
}

func (r *domainReport) ReportCreatePrimaryHostedZone(zone *client.HostedZone, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := fmt.Sprintf("%s (%s)", aurora.Green(zone.HostedZone.StackName), zone.HostedZone.HostedZoneID)

	return r.console.Report(actions, "primary-hosted-zone", description)
}

func (r *domainReport) ReportHostedZoneDelegation(zone *client.HostedZone, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := fmt.Sprintf("%s delegated status: %t", zone.HostedZone.Domain, zone.IsDelegated)

	return r.console.Report(actions, "hosted-zone", description)
}

// NewDomainReport returns an initialised domain reporter
func NewDomainReport(out io.Writer, spinner spinner.Spinner) client.DomainReport {
	return &domainReport{
		console: New(out, spinner),
	}
}
