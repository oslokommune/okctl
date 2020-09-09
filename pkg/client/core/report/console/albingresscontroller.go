package console

import (
	"fmt"
	"io"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/theckman/yacspin"
)

type albIngressControllerReport struct {
	console *Console
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
func NewAlbIngressControllerReport(out io.Writer, exit chan struct{}, spinner *yacspin.Spinner) client.ALBIngressControllerReport {
	return &albIngressControllerReport{
		console: New(out, exit, spinner),
	}
}
