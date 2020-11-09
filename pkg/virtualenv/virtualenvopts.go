package virtualenv

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/storage"
)

// CreatePs1ExecutableIfNotExists creates an executable file that returns "myenv:mynamespace", if it doesn't exist.
// The file will be called in the PS1 environment variable.
//
// This function returns the path to the directory containing the file.
func CreatePs1ExecutableIfNotExists(store storage.Storer) (string, error) {
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

	return ps1Dir, nil
}
