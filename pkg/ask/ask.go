// Package ask knows how to ask about stuff in the terminal
package ask

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/oslokommune/okctl/pkg/github"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// Ask contains stating for asking stuff
type Ask struct {
	In  terminal.FileReader
	Out terminal.FileWriter
	Err io.Writer
}

// New returns an initialised ask
func New() *Ask {
	return &Ask{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

const nameServersMsg = `	We have created a domain for you at: %s

	For this domain to resolve we need to delegate an
	AWS HostedZone to your account, to do this we need
	that you send us the domain and NS records that we
	have printed below to one of these places:

		- Slack: #kjøremiljø-support on slack
		- Email: okctl@oslo.kommune.no

	----
	Hi!

	Please create a delegated zone for our domain:

	Domain: %s
	NameServers: %s

	Thank you.
    ----
`

// ConfirmPostingNameServers asks the user to confirm that they have posted the nameservers
// to a channel where we can receive them
func (a *Ask) ConfirmPostingNameServers(to io.Writer, domain string, nameServers []string) error {
	_, err := fmt.Fprintf(to, nameServersMsg, domain, domain, strings.Join(nameServers, ","))
	if err != nil {
		return err
	}

	prompt := &survey.Confirm{
		Message: "Have you sent us the information outlined above?",
		Default: true,
		Help:    "We have printed the name of your domain and name servers above, these need to be sent to use so we can create a delegated zone",
	}

	haveSent := false

	err = survey.AskOne(prompt, &haveSent, survey.WithStdio(a.In, a.Out, a.Err))
	if err != nil {
		return err
	}

	if !haveSent {
		_, err = fmt.Fprintf(to, "! You have not sent the domain and name server, your DNS records will not resolve until you do so")
		if err != nil {
			return err
		}
	}

	return nil
}

const iacRepoMsg = `	The github repository that you select will
	be used as the infrastructure as code (IAC) repository. 

	Requirements:

		- Must be a private repository
		- Must be in oslokommune organisation
		- Should be the same repository as this
	
	It should be the same repository as you are in now, but you can 
	override this if you require so. Argo CD will be setup so it can 
	read the contents of this IAC repository and looks for new Kubernetes
	manifests for the deployments it is watching.

	If you select a different repository, please ensure that you
	move the outputs and .okctl.yml to that location.
`

// SelectInfrastructureRepository queries the user to select a repository as their IAC repository
func (a *Ask) SelectInfrastructureRepository(defaultRepo string, repos []*github.Repository) (*github.Repository, error) {
	keys := make([]string, len(repos))
	mappedRepos := make(map[string]*github.Repository, len(repos))

	for i, r := range repos {
		mappedRepos[*r.FullName] = r
		keys[i] = *r.FullName
	}

	repo := ""

	prompt := &survey.Select{
		Message: "Select repository that Argo CD will use for infrastructure as code:",
		Options: keys,
		Help:    iacRepoMsg,
		Default: defaultRepo,
	}

	isPrivateFn := func(val interface{}) error {
		s, ok := val.(survey.OptionAnswer)
		if !ok {
			return fmt.Errorf("not an option answer")
		}

		private := *mappedRepos[s.Value].Private
		if !private {
			return fmt.Errorf("repository must be private")
		}

		return nil
	}

	err := survey.AskOne(
		prompt,
		&repo,
		survey.WithValidator(survey.Required),
		survey.WithValidator(isPrivateFn),
		survey.WithStdio(a.In, a.Out, a.Err),
	)
	if err != nil {
		return nil, fmt.Errorf("failed while asking user to select repo: %w", err)
	}

	return mappedRepos[repo], nil
}
