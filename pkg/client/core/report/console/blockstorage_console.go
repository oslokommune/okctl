package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type blockstorageReport struct {
	console *Console
}

func (r *blockstorageReport) ReportDeleteBlockstorage(report *store.Report) error {
	return r.console.Report(report.Actions, "blockstorage", aurora.Green("deleting").String())
}

func (r *blockstorageReport) ReportCreateBlockstorage(secret *client.Blockstorage, report *store.Report) error {
	description := fmt.Sprintf("%s (policy), %s (service account), %s (helm)",
		aurora.Green(secret.Policy.StackName),
		aurora.Green(secret.ServiceAccount.Config.Metadata.Name),
		aurora.Green("blockstorage"),
	)

	return r.console.Report(report.Actions, "blockstorage", description)
}

// NewBlockstorageReport returns an initialised reporter
func NewBlockstorageReport(out io.Writer, spinner spinner.Spinner) client.BlockstorageReport {
	return &blockstorageReport{
		console: New(out, spinner),
	}
}
