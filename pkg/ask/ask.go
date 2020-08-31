// Package ask knows how to ask about stuff in the terminal
package ask

import (
	"fmt"
	"io"
	"os"
	"strings"

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

const message = `	We have created a domain for you at: %s

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
	_, err := fmt.Fprintf(to, message, domain, domain, strings.Join(nameServers, ","))
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
