package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type parameterReport struct {
	console *Console
}

func (p *parameterReport) SaveSecret(parameter *api.SecretParameter, report *store.Report) error {
	description := fmt.Sprintf("%s (path: %s, version: %d)",
		aurora.Green(parameter.Name),
		parameter.Path,
		parameter.Version,
	)

	return p.console.Report(report.Actions, "parameter", description)
}

// NewParameterReport returns an initialised reporter
func NewParameterReport(out io.Writer, spinner spinner.Spinner) client.ParameterReport {
	return &parameterReport{
		console: New(out, spinner),
	}
}
