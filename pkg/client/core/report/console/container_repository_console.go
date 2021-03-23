package console

import (
	"io"

	"github.com/logrusorgru/aurora/v3"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/spinner"
)

type containerRepositoryReport struct {
	console *Console
}

func (c *containerRepositoryReport) ReportCreateContainerRepository(database *client.ContainerRepository, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	return c.console.Report(actions, "container repository", aurora.Green(database.StackName).String())
}

func (c *containerRepositoryReport) ReportDeleteContainerRepository(imageName string, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	return c.console.Report(actions, "container repository", aurora.Green(imageName).String())
}

// NewContainerRepositoryReport returns an initialised container repository reporter
func NewContainerRepositoryReport(out io.Writer, spinner spinner.Spinner) client.ContainerRepositoryReport {
	return &containerRepositoryReport{
		console: New(out, spinner),
	}
}
