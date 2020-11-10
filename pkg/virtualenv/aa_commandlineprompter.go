package virtualenv

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/storage"
	"path"
	"strings"
)

type CommandLinePrompt struct {
	Warning string
	Env     map[string]string
}

// commandLinePrompter provides an interface for configuring the command line prompt
type commandLinePrompter interface {
	// setPrompt returns a map of environment variables to be used in the shell that the prompt exists within
	createPrompt() (CommandLinePrompt, error)
}

func newCommandLinePrompter(opts VirtualEnvironmentOpts, shellType shellType) (commandLinePrompter, error) {
	osEnvVars := copyMap(opts.OsEnvVars)

	noPs1, isSet := opts.OsEnvVars["OKCTL_NO_PS1"]
	if isSet && strings.ToLower(strings.TrimSpace(noPs1)) == "true" {
		return &noopPrompter{
			osEnvVars: osEnvVars,
		}, nil
	}

	err := setPs1(opts.UserDirStorage, osEnvVars)
	if err != nil {
		return nil, fmt.Errorf("could not set PS1: %w", err)
	}

	switch shellType {
	case ShellTypeBash:
		return &bashPrompter{
			environment: opts.Environment,
			osEnvVars:   osEnvVars,
		}, nil
	case ShellTypeZsh:
		return &zshPrompter{
			userDirStorage: opts.UserDirStorage,
			tmpStorer:      opts.TmpStorage,
			osEnvVars:      osEnvVars,
			environment:    opts.Environment,
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
		return fmt.Errorf("could not create PS1 executable: %w", err)
	}

	if osPath, hasPath := osEnvVars["PATH"]; hasPath {
		osEnvVars["PATH"] = fmt.Sprintf("%s:%s", ps1Dir, osPath)
	} else {
		osEnvVars["PATH"] = ps1Dir
	}

	return nil
}

// createPs1ExecutableIfNotExists creates an executable file that returns "myenv:mynamespace", if it doesn't exist.
// The file will be called in the PS1 environment variable.
//
// This function returns the path to the directory containing the file.
func createPs1ExecutableIfNotExists(store storage.Storer) (string, error) {
	ps1Filename := "venv_ps1"
	ps1Dir := "venv"

	ps1FileExists, err := store.Exists(path.Join(ps1Dir, ps1Filename))
	if err != nil {
		return "", fmt.Errorf("couldn't create PS1 helper executable: %w", err)
	}

	if !ps1FileExists {
		ps1File, err := store.Create(ps1Dir, ps1Filename, 0o744)
		if err != nil {
			return "", err
		}

		_, err = ps1File.WriteString(`#!/usr/bin/env bash
ENV=$1
ENV=${ENV:-NOENV}

K8S_NAMESPACE="$(kubectl config view --minify --output 'jsonpath={..namespace}' 2>/dev/null)"
K8S_NAMESPACE="${K8S_NAMESPACE:-default}"

echo -e "$ENV:$K8S_NAMESPACE"
`)
		if err != nil {
			return "", err
		}

		err = ps1File.Close()
		if err != nil {
			return "", err
		}
	}

	return path.Join(store.Path(), ps1Dir), nil // TODO check if actually works
}

func copyMap(m map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range m {
		res[k] = v
	}
	return res
}
