package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type monitoringReport struct {
	console *Console
}

func (r *monitoringReport) ReportSaveLoki(loki *client.Loki, report *store.Report) error {
	return r.console.Report(
		report.Actions,
		"loki",
		aurora.Green("loki").String(),
	)
}

func (r *monitoringReport) ReportRemoveLoki(report *store.Report) error {
	return r.console.Report(
		report.Actions,
		"loki",
		aurora.Green("loki").String(),
	)
}

func (r *monitoringReport) ReportRemoveKubePromStack(reports []*store.Report) error {
	var actions []store.Action // nolint

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	return r.console.Report(
		actions,
		"kube-prometheus-stack",
		aurora.Green("kube-prometheus-stack").String(),
	)
}

func (r *monitoringReport) ReportSaveKubePromStack(cd *client.KubePromStack, reports []*store.Report) error {
	var actions []store.Action // nolint

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := fmt.Sprintf("%s (url: %s)",
		aurora.Green("kube-prometheus-stack"),
		cd.Hostname,
	)

	return r.console.Report(actions, "kube-prometheus-stack", description)
}

// NewMonitoringReport returns an initialised reporter
func NewMonitoringReport(out io.Writer, spinner spinner.Spinner) client.MonitoringReport {
	return &monitoringReport{
		console: New(out, spinner),
	}
}
