package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type autoscalerReport struct {
	console *Console
}

func (r *autoscalerReport) ReportDeleteAutoscaler(report *store.Report) error {
	return r.console.Report(report.Actions, "autoscaler", aurora.Green("deleting").String())
}

func (r *autoscalerReport) ReportCreateAutoscaler(secret *client.Autoscaler, report *store.Report) error {
	description := fmt.Sprintf("%s (policy), %s (service account), %s (kubernetes)",
		aurora.Green(secret.Policy.StackName),
		aurora.Green(secret.ServiceAccount.Config.Metadata.Name),
		aurora.Green("autoscaler"),
	)

	return r.console.Report(report.Actions, "autoscaler", description)
}

// NewAutoscalerReport returns an initialised reporter
func NewAutoscalerReport(out io.Writer, spinner spinner.Spinner) client.AutoscalerReport {
	return &autoscalerReport{
		console: New(out, spinner),
	}
}
