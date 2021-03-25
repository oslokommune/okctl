package console

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/logrusorgru/aurora"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
)

type githubReport struct {
	console *Console
}

func (r *githubReport) ReportRepositoryDeployKey(repository *client.GithubRepository, report *store.Report) error {
	description := fmt.Sprintf("%s (deploykey: %s)", aurora.Green(repository.FullName), repository.DeployKey.Title)
	return r.console.Report(report.Actions, "github-repository-deploy-key", description)
}

// NewGithubReport returns an initialised reporter
func NewGithubReport(out io.Writer, spinner spinner.Spinner) client.GithubReport {
	return &githubReport{
		console: New(out, spinner),
	}
}
