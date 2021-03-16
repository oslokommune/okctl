// Package ask knows how to ask about stuff in the terminal
package ask

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/route53"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/domain"

	"github.com/miekg/dns"

	"github.com/pkg/browser"

	"github.com/oslokommune/okctl/pkg/github"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// Ask contains stating for asking stuff
type Ask struct {
	In      terminal.FileReader
	Out     terminal.FileWriter
	Err     io.Writer
	spinner spinner.Spinner
}

// New returns an initialised ask
func New() *Ask {
	return &Ask{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
}

// WithSpinner will allow the package to pause and unpause
// the spinner
func (a *Ask) WithSpinner(spinner spinner.Spinner) *Ask {
	a.spinner = spinner
	return a
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
func (a *Ask) ConfirmPostingNameServers(to io.Writer, domain string, nameServers []string) (bool, error) {
	_, err := fmt.Fprintf(to, nameServersMsg, domain, domain, strings.Join(nameServers, ","))
	if err != nil {
		return false, err
	}

	prompt := &survey.Confirm{
		Message: "Have you sent us the information outlined above?",
		Default: true,
		Help:    "We have printed the name of your domain and name servers above, these need to be sent to use so we can create a delegated zone",
	}

	haveSent := false

	if a.spinner != nil {
		err = a.spinner.Pause()
		if err != nil {
			return false, fmt.Errorf("stopping spinner: %w", err)
		}

		defer func() {
			_ = a.spinner.Unpause()
		}()
	}

	err = survey.AskOne(prompt, &haveSent, survey.WithStdio(a.In, a.Out, a.Err))
	if err != nil {
		return false, err
	}

	if !haveSent {
		_, err = fmt.Fprintf(to, "! You have not sent the domain and name server, your DNS records will not resolve until you do so")
		if err != nil {
			return false, err
		}
	}

	return true, nil
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

		private := mappedRepos[s.Value].GetPrivate()
		if !private {
			return fmt.Errorf("repository must be private")
		}

		return nil
	}

	if a.spinner != nil {
		err := a.spinner.Pause()
		if err != nil {
			return nil, fmt.Errorf("stopping spinner: %w", err)
		}

		defer func() {
			_ = a.spinner.Unpause()
		}()
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

	if a.spinner != nil {
		err := a.spinner.Pause()
		if err != nil {
			return nil, fmt.Errorf("stopping spinner: %w", err)
		}

		defer func() {
			_ = a.spinner.Unpause()
		}()
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
// nolint: funlen
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

	if a.spinner != nil {
		err := a.spinner.Pause()
		if err != nil {
			return nil, fmt.Errorf("stopping spinner: %w", err)
		}

		defer func() {
			_ = a.spinner.Unpause()
		}()
	}

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

// Domain contains the response
type Domain struct {
	FQDN   string
	Domain string
}

// Domain asks the user if they accept the given domain
func (a *Ask) Domain(defaultDomain string) (*Domain, error) {
	for {
		d, err := a.askDomain(defaultDomain)
		if err != nil {
			if err == terminal.InterruptErr {
				return nil, err
			}

			_, err = fmt.Fprintln(a.Err, err.Error())
			if err != nil {
				return nil, err
			}

			continue
		}

		return &Domain{
			FQDN:   dns.Fqdn(d),
			Domain: d,
		}, nil
	}
}

func (a *Ask) askDomain(defaultDomain string) (string, error) {
	d := ""

	q := &survey.Input{
		Message: "Provide the name of the domain you want to delegate to this cluster",
		Default: defaultDomain,
		Help:    "This is the domain name we will delegate to your AWS account and that the cluster will create hostnames from",
	}

	validatorFn := func(val interface{}) error {
		f, ok := val.(string)
		if !ok {
			return fmt.Errorf("could not convert input to a string")
		}

		err := domain.Validate(f)
		if err != nil {
			return err
		}

		return domain.NotTaken(f)
	}

	if a.spinner != nil {
		err := a.spinner.Pause()
		if err != nil {
			return "", fmt.Errorf("stopping spinner: %w", err)
		}

		defer func() {
			_ = a.spinner.Unpause()
		}()
	}

	err := survey.AskOne(q, &d,
		survey.WithStdio(a.In, a.Out, a.Err),
		survey.WithValidator(survey.Required),
		survey.WithValidator(validatorFn),
	)
	if err != nil {
		return "", err
	}

	return d, nil
}

// Username asks the user for their username
func (a *Ask) Username() (string, error) {
	username := ""

	prompt := &survey.Input{
		Message: "Your username:",
		Help:    "This is your AD user, e.g., yyyXXXXXX (y = letter, x = digit). We store it in the application configuration, so you don't have to enter it each time.",
	}

	if a.spinner != nil {
		err := a.spinner.Pause()
		if err != nil {
			return "", fmt.Errorf("stopping spinner: %w", err)
		}

		defer func() {
			_ = a.spinner.Unpause()
		}()
	}

	err := survey.AskOne(prompt, &username, survey.WithStdio(a.In, a.Out, a.Err))
	if err != nil {
		return "", err
	}

	return username, nil
}

// RepositoryConfig contains the user provided inputs
type RepositoryConfig struct {
	Name    string
	Region  string
	BaseDir string
}

// RepositoryConfig asks the user for repo configuration
func (a *Ask) RepositoryConfig() (*RepositoryConfig, error) {
	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Project or product name:",
				Help:    "A descriptive name, e.g., team or project, used among other things to prefix AWS resources",
			},
		},
		{
			Name: "region",
			Prompt: &survey.Select{
				Message: "Choose AWS region:",
				Options: v1alpha1.SupportedRegions(),
				Help:    "The AWS region to create resources in",
			},
		},
		{
			Name: "basedir",
			Prompt: &survey.Input{
				Message: "Output directory:",
				Default: constant.DefaultOutputDirectory,
				Help:    "Directory in the repository to store the data in",
			},
		},
	}

	answers := struct {
		Name    string
		Region  string
		Basedir string
	}{}

	if a.spinner != nil {
		err := a.spinner.Pause()
		if err != nil {
			return nil, fmt.Errorf("stopping spinner: %w", err)
		}

		defer func() {
			_ = a.spinner.Unpause()
		}()
	}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return nil, fmt.Errorf("getting repository config: %w", err)
	}

	return &RepositoryConfig{
		Name:    answers.Name,
		Region:  answers.Region,
		BaseDir: answers.Basedir,
	}, nil
}

const zoneHelp = `
The hosted zone you select will be associated
with the cluster, but its lifecycle will not
be managed by it, e.g., we will not remove it.
`

// SelectHostedZone queries the user to select a team for authorisation towards argo cd
func (a *Ask) SelectHostedZone(zones []*route53.HostedZone) (*route53.HostedZone, error) {
	keys := make([]string, len(zones))
	mappedZones := make(map[string]*route53.HostedZone, len(zones))

	for i, z := range zones {
		mappedZones[z.Domain] = z
		keys[i] = z.Domain
	}

	repo := ""

	prompt := &survey.Select{
		Message: "Select hosted zone that you want to use with your cluster: ",
		Options: keys,
		Help:    zoneHelp,
	}

	if a.spinner != nil {
		err := a.spinner.Pause()
		if err != nil {
			return nil, fmt.Errorf("stopping spinner: %w", err)
		}

		defer func() {
			_ = a.spinner.Unpause()
		}()
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

	return mappedZones[repo], nil
}
