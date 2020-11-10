package commandlineprompter

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/storage"
	"strings"
)

type zshPrompter struct {
	userDirStorage storage.Storer
	tmpStorer      storage.Storer
	osEnvVars      map[string]string
	environment    string
}

type CreateZshPromptWarning struct {
	Warning string
	Err     error
}

func (e *CreateZshPromptWarning) Unwrap() error {
	return e.Err
}

// createPrompt makes zsh show a custom command prompt. It does so by creating a temporary .zshrc file.
// The return value is this temporary file, which should be cleaned up after use.
func (p *zshPrompter) CreatePrompt() (CommandLinePrompt, error) {
	if _, ok := p.osEnvVars["ZDOTDIR"]; ok {
		// We're dependent on being able to set ZDOTDIR ourself to launch zsh to a temporary path with a custom .zshrc
		// file. If the user has already set ZDOTDIR, we cannot do this. However, instead of returning an error, show
		// this warning to the user.
		msg := "WARNING: Could not set command prompt (PS1) because ZDOTDIR is already set. "
		msg += "Either start okctl venv with no ZDOTDIR set, or set environment variable OKCTL_NO_PS1=true to get "
		msg += "rid of this message"

		return CommandLinePrompt{
			Warning: msg,
			Env:     p.osEnvVars,
		}, nil
	}

	zshrcDir, err := p.writeZshrcFile()
	if err != nil {
		return CommandLinePrompt{}, err
	}

	// Make zsh use the new .zshrc file
	p.osEnvVars["ZDOTDIR"] = zshrcDir

	return CommandLinePrompt{
		Env: p.osEnvVars,
	}, nil
}

func (p *zshPrompter) writeZshrcFile() (string, error) {
	zshrcTmpFile, err := p.tmpStorer.Create(".", ".zshrc", 0o644)
	if err != nil {
		return "", err
	}

	zshrcContents, err := p.createZshrcContents()
	if err != nil {
		return "", err
	}

	_, err = zshrcTmpFile.WriteString(zshrcContents)
	if err != nil {
		return "", err
	}

	return p.tmpStorer.Path(), nil
}

func (p *zshPrompter) createZshrcContents() (string, error) {
	zshrcBuilder := strings.Builder{}

	zshrcExists, err := p.userDirStorage.Exists(".zshrc")
	if err != nil {
		return "", err
	}

	if zshrcExists {
		zshrcBuilder.WriteString("source ~/.zshrc\n\n")
	}

	zshrcBuilder.WriteString(`setopt PROMPT_SUBST
autoload -U colors && colors # Enable colors
prompt() {
`)

	ps1, overridePs1 := p.osEnvVars["OKCTL_PS1"]
	if overridePs1 {
		zshrcBuilder.WriteString(fmt.Sprintf("PS1=%s", ps1))
	} else {
		zshrcBuilder.WriteString(fmt.Sprint(`PS1="%F{red}%~ %f%F{blue}($(venv_ps1 ` + p.environment + `)%f) $ "`))
	}

	zshrcBuilder.WriteString(`
}
precmd_functions+=(prompt)
`)

	return zshrcBuilder.String(), nil
}
