package commandlineprompter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/oslokommune/okctl/pkg/storage"
)

const zshrcFilename = ".zshrc"

var reClusterDeclarationExport = regexp.MustCompile("export OKCTL_CLUSTER_DECLARATION.*")

type zshPrompter struct {
	userHomeDirStorage storage.Storer
	tmpStorer          storage.Storer
	osEnvVars          map[string]string
	clusterName        string
}

// CreatePrompt returns environment variables that when set in zsh will show a command prompt.
// The warning is set in case something prevented the prompt to be set the expected way.
func (p *zshPrompter) CreatePrompt() (CommandLinePrompt, error) {
	if _, ok := p.osEnvVars["ZDOTDIR"]; ok {
		// We're dependent on being able to set ZDOTDIR ourself to launch zsh to a temporary path with a custom .zshrc
		// file. If the user has already set ZDOTDIR, we cannot do this. However, instead of returning an error, show
		// this warning to the user.
		msg := "WARNING: Could not set command prompt (PS1) because ZDOTDIR is already set. "
		msg += "Either start okctl venv with no ZDOTDIR set, or set environment variable OKCTL_NO_PS1=true to get "
		msg += "rid of this message."

		return CommandLinePrompt{
			Warning: msg,
			Env:     p.osEnvVars,
		}, nil
	}

	zshrcDir, err := p.writeZshrcFile()
	if err != nil {
		return CommandLinePrompt{}, fmt.Errorf("could not write .zshrc file: %w", err)
	}

	// Make zsh use the new .zshrc file
	p.osEnvVars["ZDOTDIR"] = zshrcDir

	return CommandLinePrompt{
		Env: p.osEnvVars,
	}, nil
}

func (p *zshPrompter) writeZshrcFile() (string, error) {
	zshrcTmpFile, err := p.tmpStorer.Create(".", zshrcFilename, 0o644)
	if err != nil {
		return "", fmt.Errorf("could not create temporary .zshrc file: %w", err)
	}

	zshrcContents, err := p.createZshrcContents()
	if err != nil {
		return "", fmt.Errorf("could not create temporary .zshrc contents: %w", err)
	}

	_, err = zshrcTmpFile.WriteString(zshrcContents)
	if err != nil {
		return "", fmt.Errorf("could not write contents to temporary .zshrc file: %w", err)
	}

	return p.tmpStorer.Path(), nil
}

func (p *zshPrompter) createZshrcContents() (string, error) {
	zshrcBuilder := strings.Builder{}

	zshrcExists, err := p.userHomeDirStorage.Exists(zshrcFilename)
	if err != nil {
		return "", fmt.Errorf("could not check if temporary .zshrc file exists: %w", err)
	}

	if zshrcExists {
		zshrcFileContents, err := p.userHomeDirStorage.ReadAll(zshrcFilename)
		if err != nil {
			return "", fmt.Errorf("reading existing .zshrc content: %w", err)
		}

		cleanedContent := reClusterDeclarationExport.ReplaceAll(zshrcFileContents, []byte(""))

		zshrcBuilder.Write(cleanedContent)
	}

	zshrcBuilder.WriteString(`setopt PROMPT_SUBST
autoload -U colors && colors # Enable colors
prompt() {
`)

	ps1, overridePs1 := p.osEnvVars["OKCTL_PS1"]
	if overridePs1 {
		withEnv := strings.ReplaceAll(ps1, "%env", p.clusterName)
		zshrcBuilder.WriteString(fmt.Sprintf(`PS1="%s"`, withEnv))
	} else {
		zshrcBuilder.WriteString(fmt.Sprint(`PS1="%F{red}%~ %f%F{blue}($(venv_ps1 ` + p.clusterName + `)%f) $ "`))
	}

	zshrcBuilder.WriteString(`
}
precmd_functions+=(prompt)
`)

	return zshrcBuilder.String(), nil
}
