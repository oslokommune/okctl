package console

import (
	"io"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/spinner"
)

type serviceAccountReport struct {
	console *Console
}

func (r *serviceAccountReport) ReportCreateServiceAccount(s *api.ServiceAccount, report *store.Report) error {
	return r.console.Report(
		report.Actions,
		"service-account",
		aurora.Green(s.Name).String(),
	)
}

func (r *serviceAccountReport) ReportDeleteServiceAccount(name string, report *store.Report) error {
	return r.console.Report(
		report.Actions,
		"service-account",
		aurora.Green(name).String(),
	)
}

// NewServiceAccountReport returns an initialised reporter
func NewServiceAccountReport(out io.Writer, spinner spinner.Spinner) client.ServiceAccountReport {
	return &serviceAccountReport{
		console: New(out, spinner),
	}
}
