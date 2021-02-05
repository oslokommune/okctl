package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type awsLoadBalancerControllerReport struct {
	console *Console
}

func (r *awsLoadBalancerControllerReport) ReportDeleteAWSLoadBalancerController(report *store.Report) error {
	return r.console.Report(report.Actions, "aws-load-balancer-controller", aurora.Green("deleting").String())
}

func (r *awsLoadBalancerControllerReport) ReportCreateAWSLoadBalancerController(controller *client.AWSLoadBalancerController, report *store.Report) error {
	description := fmt.Sprintf("%s (policy), %s (service account), %s (chart)",
		aurora.Green(controller.Policy.StackName),
		aurora.Green(controller.ServiceAccount.Config.Metadata.Name),
		aurora.Green(controller.Chart.Release.Name),
	)

	return r.console.Report(report.Actions, "alb-ingress-controller", description)
}

// NewAWSLoadBalancerControllerReport returns an initialised reporter
func NewAWSLoadBalancerControllerReport(out io.Writer, spinner spinner.Spinner) client.AWSLoadBalancerControllerReport {
	return &awsLoadBalancerControllerReport{
		console: New(out, spinner),
	}
}
