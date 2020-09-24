package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type argoCDReport struct {
	console *Console
}

func (r *argoCDReport) CreateArgoCD(cd *client.ArgoCD, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := fmt.Sprintf("%s (url: %s)",
		aurora.Green("argocd"),
		cd.ArgoURL,
	)

	return r.console.Report(actions, "argocd", description)
}

// NewArgoCDReport returns an initialised reporter
func NewArgoCDReport(out io.Writer, spinner spinner.Spinner) client.ArgoCDReport {
	return &argoCDReport{
		console: New(out, spinner),
	}
}
