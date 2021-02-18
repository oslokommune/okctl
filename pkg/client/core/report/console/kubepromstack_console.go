package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type kubePromStackReport struct {
	console *Console
}

func (r *kubePromStackReport) ReportKubePromStack(cd *client.KubePromStack, reports []*store.Report) error {
	var actions []store.Action // nolint

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := fmt.Sprintf("%s (url: %s)",
		aurora.Green("kube-prometheus-stack"),
		cd.Hostname,
	)

	return r.console.Report(actions, "KubePromStack", description)
}

// NewKubePromStackReport returns an initialised reporter
func NewKubePromStackReport(out io.Writer, spinner spinner.Spinner) client.KubePromStackReport {
	return &kubePromStackReport{
		console: New(out, spinner),
	}
}
