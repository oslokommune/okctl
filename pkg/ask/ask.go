// Package ask knows how to ask about stuff in the terminal
package ask

import (
	"fmt"
	"io"
	"os"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

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
