package console

import (
	"io"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/spinner"
)

type identityManagerReport struct {
	console *Console
}

func (r *identityManagerReport) ReportDeleteIdentityPool(report *store.Report) error {
	return r.console.Report(report.Actions, "identitypool", aurora.Green("deleting").String())
}

func (r *identityManagerReport) ReportIdentityPoolUser(client *api.IdentityPoolUser, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := aurora.Green("identity-pool-user").String()

	return r.console.Report(actions, "identity-manager", description)
}

func (r *identityManagerReport) ReportIdentityPoolClient(client *api.IdentityPoolClient, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := aurora.Green("identity-pool-client").String()

	return r.console.Report(actions, "identity-manager", description)
}

func (r *identityManagerReport) ReportIdentityPool(pool *api.IdentityPool, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := aurora.Green("identity-pool").String()

	return r.console.Report(actions, "identity-manager", description)
}

// NewIdentityManagerReport returns an initialised reporter
func NewIdentityManagerReport(out io.Writer, spinner spinner.Spinner) client.IdentityManagerReport {
	return &identityManagerReport{
		console: New(out, spinner),
	}
}
