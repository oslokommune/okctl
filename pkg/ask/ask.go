// Package ask knows how to ask about stuff in the terminal
package ask

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/browser"

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

const iacHelp = `	The github repository that you select will
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
		Message: "Select infrastructure as code repository:",
		Options: keys,
		Help:    iacHelp,
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

const teamHelp = `	The team you select will be the only one
	that is allowed to login to the Argo CD UI using Github
	as the authorisation provider.`

// SelectTeam queries the user to select a team for authorisation towards argo cd
func (a *Ask) SelectTeam(teams []*github.Team) (*github.Team, error) {
	keys := make([]string, len(teams))
	mappedTeams := make(map[string]*github.Team, len(teams))

	for i, r := range teams {
		mappedTeams[r.GetName()] = r
		keys[i] = r.GetName()
	}

	repo := ""

	prompt := &survey.Select{
		Message: "Select team that is authorised to access the Argo CD UI:",
		Options: keys,
		Help:    teamHelp,
	}

	err := survey.AskOne(
		prompt,
		&repo,
		survey.WithValidator(survey.Required),
		survey.WithStdio(a.In, a.Out, a.Err),
	)
	if err != nil {
		return nil, fmt.Errorf("failed while asking user to select team: %w", err)
	}

	return mappedTeams[repo], nil
}

const oauthAppMsg = `	We will now be creating an oauth application
	for authorising access to a UI. Below you will find the
	the steps to follow and information you will need to fill
	in:

	1. Open the following URL: https://github.com/settings/developers
	   or answer Yes to the terminal query.

	2. Click on "Register a new application"

	3. Fill in the fields with the information below:

		Application name: '%s'
		Homepage URL: '%s'
		Application description: '%s'
		Authorization callback URL: '%s'

	4. Click on "Register application"

	5. Respond to the terminal queries for the "Client ID" and
	   "Client Secret"

	6. Click on the "Transfer Ownership" button

	7. Enter the following information:

		Type the name of the application to confirm: '%s'
		New owner’s GitHub username or organization name: '%s'

	8. Click on "Transfer this application"

`

const oauthAppHelp = `	We will be creating an oauth app that authorises
	those who have access to a given UI. It is not possible to
	automate this task, and therefore it requires some manual
	work.
`

// OauthApp contains the oauth app state
type OauthApp struct {
	Name         string
	Organisation string
	URL          string
	CallbackURL  string
	ClientID     string
	ClientSecret string
}

// OauthAppOpts contains the inputs
type OauthAppOpts struct {
	Organisation string
	Name         string
	URL          string
	CallbackURL  string
}

// CreateOauthApp helps the user create and transfer and oauth app
func (a *Ask) CreateOauthApp(to io.Writer, opts OauthAppOpts) (*OauthApp, error) {
	_, err := fmt.Fprintf(to, oauthAppMsg, opts.Name, opts.URL, opts.Name, opts.CallbackURL, opts.Name, opts.Organisation)
	if err != nil {
		return nil, err
	}

	prompt := &survey.Confirm{
		Message: "Attempt to open browser window to github oauth apps?",
		Default: true,
		Help:    oauthAppHelp,
	}

	doOpen := false

	err = survey.AskOne(prompt, &doOpen, survey.WithStdio(a.In, a.Out, a.Err))
	if err != nil {
		return nil, err
	}

	if doOpen {
		_ = browser.OpenURL("https://github.com/settings/developers")
	}

	qs := []*survey.Question{
		{
			Name:     "clientid",
			Prompt:   &survey.Input{Message: "Enter Client ID of the oauth app:"},
			Validate: survey.Required,
		},
		{
			Name:     "clientsecret",
			Prompt:   &survey.Password{Message: "Enter Client Secret of the oauth app:"},
			Validate: survey.Required,
		},
	}

	answers := struct {
		ClientID     string
		ClientSecret string
	}{}

	err = survey.Ask(qs, &answers, survey.WithStdio(a.In, a.Out, a.Err))
	if err != nil {
		return nil, err
	}

	return &OauthApp{
		Name:         opts.Name,
		Organisation: opts.Organisation,
		URL:          opts.URL,
		CallbackURL:  opts.CallbackURL,
		ClientID:     answers.ClientID,
		ClientSecret: answers.ClientSecret,
	}, nil
}
