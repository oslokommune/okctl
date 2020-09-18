package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type albIngressControllerReport struct {
	console *Console
}

func (r *albIngressControllerReport) ReportDeleteALBIngressController(report *store.Report) error {
	return r.console.Report(report.Actions, "alb-ingress-controller", aurora.Green("deleting").String())
}

func (r *albIngressControllerReport) ReportCreateALBIngressController(controller *client.ALBIngressController, report *store.Report) error {
	description := fmt.Sprintf("%s (policy), %s (service account), %s (chart)",
		aurora.Green(controller.Policy.StackName),
		aurora.Green(controller.ServiceAccount.Config.Metadata.Name),
		aurora.Green(controller.Chart.Release.Name),
	)

	return r.console.Report(report.Actions, "alb-ingress-controller", description)
}

// NewAlbIngressControllerReport returns an initialised reporter
func NewAlbIngressControllerReport(out io.Writer, spinner spinner.Spinner) client.ALBIngressControllerReport {
	return &albIngressControllerReport{
		console: New(out, spinner),
	}
}
