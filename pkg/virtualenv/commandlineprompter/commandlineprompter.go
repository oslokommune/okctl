// Package commandlineprompter implements functionality for setting the command prompt (PS1) for various shells
package commandlineprompter

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"path"
	"strings"

	"github.com/oslokommune/okctl/pkg/virtualenv/shelltype"

	"github.com/oslokommune/okctl/pkg/storage"
)

// CommandLinePrompt contains environment variables needed by the command prompt, and possibly a warning about issues
// with the returned command prompt.
type CommandLinePrompt struct {
	Warning string
	Env     map[string]string
}

// CommandLinePrompter defines functionality for creating a command line prompt
type CommandLinePrompter interface {
	CreatePrompt() (CommandLinePrompt, error)
}

// New creates a new command line prompter
func New(opts CommandLinePromptOpts, shellType shelltype.ShellType) (CommandLinePrompter, error) {
	osEnvVars := copyMap(opts.OsEnvVars)

	noPs1, isSet := opts.OsEnvVars["OKCTL_NO_PS1"]
	if isSet && strings.ToLower(strings.TrimSpace(noPs1)) == "true" {
		return &noopPrompter{
			osEnvVars: osEnvVars,
		}, nil
	}

	err := setPs1(opts.UserDirStorage, osEnvVars)
	if err != nil {
		return nil, fmt.Errorf(constant.SetPS1Error, err)
	}

	switch shellType {
	case shelltype.Bash:
		return &bashPrompter{
			clusterName: opts.ClusterName,
			osEnvVars:   osEnvVars,
		}, nil
	case shelltype.Zsh:
		return &zshPrompter{
			userHomeDirStorage: opts.UserHomeDirStorage,
			tmpStorer:          opts.TmpStorage,
			osEnvVars:          osEnvVars,
			clusterName:        opts.ClusterName,
		}, nil
	default:
		return &noopPrompter{
			osEnvVars: osEnvVars,
		}, nil
	}
}

func setPs1(userDirStorage storage.Storer, osEnvVars map[string]string) error {
	ps1Dir, err := createPs1ExecutableIfNotExists(userDirStorage)
	if err != nil {
		return fmt.Errorf(constant.CreateExecutablePS1Error, err)
	}

	if osPath, hasPath := osEnvVars["PATH"]; hasPath {
		osEnvVars["PATH"] = fmt.Sprintf("%s:%s", ps1Dir, osPath)
	} else {
		osEnvVars["PATH"] = ps1Dir
	}

	return nil
}

const (
	// Ps1Dir is the name of the directory where the PS1 executable will be stored
	Ps1Dir = "venv"

	// Ps1Filename is the file name of the PS1 executable that returns data to be put into the command prompt
	Ps1Filename = "venv_ps1"
)

// createPs1ExecutableIfNotExists creates an executable file that returns "myenv:mynamespace", if it doesn't exist.
// The file will be called in the PS1 environment variable.
//
// This function returns the path to the directory containing the file.
func createPs1ExecutableIfNotExists(store storage.Storer) (string, error) {
	ps1FileExists, err := store.Exists(path.Join(Ps1Dir, Ps1Filename))
	if err != nil {
		return "", fmt.Errorf(constant.CheckIfPS1ExecutableExistsError, err)
	}

	if !ps1FileExists {
		ps1File, err := store.Create(Ps1Dir, Ps1Filename, 0o744)
		if err != nil {
			return "", fmt.Errorf(constant.UnableToCreatePS1FileError, err)
		}

		_, err = ps1File.WriteString(`#!/usr/bin/env bash
ENV=$1
ENV=${ENV:-NOENV}

K8S_NAMESPACE="$(kubectl config view --minify --output 'jsonpath={..namespace}' 2>/dev/null)"
K8S_NAMESPACE="${K8S_NAMESPACE:-default}"

echo -e "$ENV:$K8S_NAMESPACE"
`)
		if err != nil {
			return "", fmt.Errorf(constant.WriteContentToPS1FileError, err)
		}

		err = ps1File.Close()
		if err != nil {
			return "", fmt.Errorf(constant.ClosePS1FileError, err)
		}
	}

	return path.Join(store.Path(), Ps1Dir), nil
}

func copyMap(m map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range m {
		res[k] = v
	}

	return res
}
