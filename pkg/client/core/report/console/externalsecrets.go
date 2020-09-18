package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type externalSecretsReport struct {
	console *Console
}

func (r *externalSecretsReport) ReportDeleteExternalSecrets(report *store.Report) error {
	return r.console.Report(report.Actions, "external-secrets", aurora.Green("deleting").String())
}

func (r *externalSecretsReport) ReportCreateExternalSecrets(secret *client.ExternalSecrets, report *store.Report) error {
	description := fmt.Sprintf("%s (policy), %s (service account), %s (chart)",
		aurora.Green(secret.Policy.StackName),
		aurora.Green(secret.ServiceAccount.Config.Metadata.Name),
		aurora.Green(secret.Chart.Release.Name),
	)

	return r.console.Report(report.Actions, "external-secrets", description)
}

// NewExternalSecretsReport returns an initialised reporter
func NewExternalSecretsReport(out io.Writer, spinner spinner.Spinner) client.ExternalSecretsReport {
	return &externalSecretsReport{
		console: New(out, spinner),
	}
}
