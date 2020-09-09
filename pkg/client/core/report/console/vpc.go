package console

import (
	"fmt"
	"io"

	"github.com/theckman/yacspin"

	"github.com/logrusorgru/aurora/v3"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type vpcReport struct {
	console *Console
}

func (r *vpcReport) ReportCreateVPC(vpc *api.Vpc, report *store.Report) error {
	description := fmt.Sprintf("%s (%s)", aurora.Green(vpc.StackName), vpc.VpcID)
	return r.console.Report(report.Actions, "vpc", description)
}

// NewVPCReport returns an initialised VPC reporter
func NewVPCReport(out io.Writer, spinner *yacspin.Spinner, exit chan struct{}) client.VPCReport {
	return &vpcReport{
		console: New(out, exit, spinner),
	}
}
