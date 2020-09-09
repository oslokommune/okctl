package console

import (
	"fmt"
	"io"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/theckman/yacspin"
)

type externalSecretsReport struct {
	console *Console
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
func NewExternalSecretsReport(out io.Writer, exit chan struct{}, spinner *yacspin.Spinner) client.ExternalSecretsReport {
	return &externalSecretsReport{
		console: New(out, exit, spinner),
	}
}
