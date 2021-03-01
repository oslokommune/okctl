package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type manifestReport struct {
	console *Console
}

func (r *manifestReport) SaveConfigMap(secret *client.ConfigMap, report *store.Report) error {
	return r.console.Report(
		report.Actions,
		"manifest",
		fmt.Sprintf("%s (%s)", aurora.Green("config-map").String(), secret.Name),
	)
}

func (r *manifestReport) RemoveConfigMap(report *store.Report) error {
	return r.console.Report(
		report.Actions,
		"manifest",
		aurora.Green("config-map").String(),
	)
}

func (r *manifestReport) RemoveExternalSecret(report *store.Report) error {
	return r.console.Report(
		report.Actions,
		"manifest",
		aurora.Green("external-secret").String(),
	)
}

func (r *manifestReport) SaveStorageClass(_ *client.StorageClass, report *store.Report) error {
	description := aurora.Green("storage-class").String()
	return r.console.Report(report.Actions, "manifest", description)
}

func (r *manifestReport) SaveExternalSecret(_ *client.ExternalSecret, report *store.Report) error {
	description := aurora.Green("external-secret").String()
	return r.console.Report(report.Actions, "manifest", description)
}

// NewManifestReport returns an initialised manifest reporter
func NewManifestReport(out io.Writer, spinner spinner.Spinner) client.ManifestReport {
	return &manifestReport{
		console: New(out, spinner),
	}
}
