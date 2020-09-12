package console

import (
	"fmt"
	"io"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/theckman/yacspin"
)

type parameterReport struct {
	console *Console
}

func (p *parameterReport) SaveSecret(parameter *api.SecretParameter, report *store.Report) error {
	description := fmt.Sprintf("%s (path: %s, version: %d)",
		aurora.Blue(parameter.Name),
		parameter.Path,
		parameter.Version,
	)

	return p.console.Report(report.Actions, "parameter", description)
}

// NewParameterReport returns an initialised reporter
func NewParameterReport(out io.Writer, exit chan struct{}, spinner *yacspin.Spinner) client.ParameterReport {
	return &parameterReport{
		console: New(out, exit, spinner),
	}
}
