package virtualenv

import (
	"fmt"
	"path"
	"strings"

	"github.com/oslokommune/okctl/pkg/storage"
)

// ShellIsZsh returns true if provided command will run zsh
func ShellIsZsh(shellCmd string) bool {
	return strings.HasSuffix(shellCmd, "zsh")
}

// SetCmdPromptZsh makes zsh show a custom command prompt. It does so by creating a temporary .zshrc file.
// The return value is this temporary file, which should be cleaned up after use.
func SetCmdPromptZsh(opts *VirtualEnvironmentOpts, venv *VirtualEnvironment, fileSystemReader storage.Storer, tmpStorer storage.Storer) error {
	err := writeZshrcFile(opts, venv, fileSystemReader, tmpStorer)
	if err != nil {
		return err
	}

	// Make zsh use the new .zshrc file
	venv.env["ZDOTDIR"] = tmpStorer.Path()

	return nil
}

func writeZshrcFile(opts *VirtualEnvironmentOpts, venv *VirtualEnvironment, fileSystemReader storage.Storer, tmpStorer storage.Storer) error {
	zshrcTmpFile, err := tmpStorer.Create(".", ".zshrc", 0o644)
	if err != nil {
		return err
	}

	zshrcContents, err := createZshrcContents(opts, venv, fileSystemReader)
	if err != nil {
		return err
	}

	_, err = zshrcTmpFile.WriteString(zshrcContents)
	if err != nil {
		return err
	}

	return nil
}

func createZshrcContents(opts *VirtualEnvironmentOpts, venv *VirtualEnvironment, fileSystemReader storage.Storer) (string, error) {
	zshrcBuilder := strings.Builder{}

	zDotDir, zDotDirExists := venv.Getenv("ZDOTDIR")

	zshrcExists, err := fileSystemReader.Exists(".zshrc")
	if err != nil {
		return "", err
	}

	// zsh uses ZDOTDIR to read startup files from, if set. If not, it uses $HOME.
	// http://zsh.sourceforge.net/Intro/intro_3.html
	// This code ensures we follow this contract, in case the user has set ZDOTDIR.
	if zDotDirExists {
		zshrcPath := path.Join(zDotDir, ".zshrc")
		zshrcBuilder.WriteString(fmt.Sprintf("source %s", zshrcPath))
	} else if zshrcExists {
		zshrcBuilder.WriteString("source ~/.zshrc\n\n")
	}

	zshrcBuilder.WriteString(`setopt PROMPT_SUBST
autoload -U colors && colors # Enable colors
prompt() {
`)

	ps1, overridePs1 := venv.env["OKCTL_PS1"]
	if overridePs1 {
		zshrcBuilder.WriteString(fmt.Sprintf("PS1=%s", ps1))
	} else {
		zshrcBuilder.WriteString(fmt.Sprint(`PS1="%F{red}%~ %f%F{blue}($(venv_ps1 ` + opts.Environment + `)%f) $ "`))
	}

	zshrcBuilder.WriteString(`
}
precmd_functions+=(prompt)
`)

	return zshrcBuilder.String(), nil
}
