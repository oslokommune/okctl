package console

import (
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type clusterReport struct {
	console *Console
}

func (r *clusterReport) ReportCreateCluster(cluster *api.Cluster, reports []*store.Report) error {
	var actions []store.Action // nolint: prealloc

	for _, report := range reports {
		actions = append(actions, report.Actions...)
	}

	description := aurora.Green(cluster.Config.Metadata.Name).String()

	return r.console.Report(actions, "cluster", description)
}

// NewClusterReport returns an initialised cluster reporter
func NewClusterReport(out io.Writer, spinner spinner.Spinner) client.ClusterReport {
	return &clusterReport{
		console: New(out, spinner),
	}
}
