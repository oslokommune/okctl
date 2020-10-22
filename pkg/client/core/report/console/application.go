package console

import (
	"io"

	"github.com/logrusorgru/aurora/v3"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/spinner"
)

type applicationReport struct {
	console *Console
}

// ReportDeleteApplication produces relevant output upon application deletion
func (r *applicationReport) ReportDeleteApplication(reports []*store.Report) error {
	var actions []store.Action

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	return r.console.Report(actions, "application", aurora.Green("removing").String())
}

// ReportCreateApplication produces relevant output upon application creation
func (r *applicationReport) ReportCreateApplication(application *client.ScaffoldedApplication, reports []*store.Report) error {
	var actions []store.Action

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	return r.console.Report(actions, "application", application.ApplicationName)
}

// NewApplicationReport returns an initialized application reporter
func NewApplicationReport(out io.Writer, spinner spinner.Spinner) client.ApplicationReport {
	return &applicationReport{
		console: New(out, spinner),
	}
}
