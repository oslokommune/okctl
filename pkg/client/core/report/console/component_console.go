package console

import (
	"io"

	"github.com/logrusorgru/aurora/v3"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/spinner"
)

type componentReport struct {
	console *Console
}

func (c *componentReport) ReportCreatePostgresDatabase(database *client.PostgresDatabase, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	return c.console.Report(actions, "postgres", aurora.Green(database.StackName).String())
}

func (c *componentReport) ReportDeletePostgresDatabase(applicationName string, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	return c.console.Report(actions, "postgres", aurora.Green(applicationName).String())
}

// NewComponentReport returns an initialised component reporter
func NewComponentReport(out io.Writer, spinner spinner.Spinner) client.ComponentReport {
	return &componentReport{
		console: New(out, spinner),
	}
}
