package console

import (
	"io"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/spinner"
)

type managedPolicyReport struct {
	console *Console
}

func (m *managedPolicyReport) ReportCreatePolicy(policy *api.ManagedPolicy, report *store.Report) error {
	return m.console.Report(
		report.Actions,
		"managed-policy",
		aurora.Green(policy.StackName).String(),
	)
}

func (m *managedPolicyReport) ReportDeletePolicy(stackName string, report *store.Report) error {
	return m.console.Report(
		report.Actions,
		"managed-policy",
		aurora.Green(stackName).String(),
	)
}

// NewManagedPolicyReport returns an initialised reporter
func NewManagedPolicyReport(out io.Writer, spinner spinner.Spinner) client.ManagedPolicyReport {
	return &managedPolicyReport{
		console: New(out, spinner),
	}
}
