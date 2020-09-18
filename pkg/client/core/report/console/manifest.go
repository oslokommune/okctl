package console

import (
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type manifestReport struct {
	console *Console
}

func (m *manifestReport) SaveExternalSecret(_ *client.ExternalSecret, report *store.Report) error {
	description := aurora.Blue("external-secret").String()
	return m.console.Report(report.Actions, "manifest", description)
}

// NewManifestReport returns an initialised manifest reporter
func NewManifestReport(out io.Writer, spinner spinner.Spinner) client.ManifestReport {
	return &manifestReport{
		console: New(out, spinner),
	}
}
