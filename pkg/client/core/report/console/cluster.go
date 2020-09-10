package console

import (
	"io"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/theckman/yacspin"
)

type clusterReport struct {
	console *Console
}

func (r *clusterReport) ReportCreateCluster(cluster *api.Cluster, report *store.Report) error {
	description := aurora.Green(cluster.Config.Metadata).String()
	return r.console.Report(report.Actions, "cluster", description)
}

// NewClusterReport returns an initialised cluster reporter
func NewClusterReport(out io.Writer, exit chan struct{}, spinner *yacspin.Spinner) client.ClusterReport {
	return &clusterReport{
		console: New(out, exit, spinner),
	}
}
