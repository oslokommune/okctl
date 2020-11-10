// Package virtualenv helps finding the environment variables needed to use a okctl cluster.
package virtualenv

import (
	"fmt"
	"sort"
)

// VirtualEnvironment contains environment variables in a virtual environment.
type VirtualEnvironment struct {
	env          map[string]string
	Warning      string
	ShellCommand string
}

// Environ returns all environment variables in the virtual environment, on the form
// []string { "KEY1=VALUE1", "KEY2=VALUE2", ... }
// This is the same form as os.Environ.
func (v *VirtualEnvironment) Environ() []string {
	return toEnvVarsSlice(&v.env)
}

// Getenv returns the environment variable with the given key, and a bool indicating if the key was found or not.
func (v *VirtualEnvironment) Getenv(key string) (string, bool) {
	val, hasKey := v.env[key]
	return val, hasKey
}

func toEnvVarsSlice(venv *map[string]string) []string {
	venvs := make([]string, 0, len(*venv))

	for k, v := range *venv {
		keyEqualsValue := fmt.Sprintf("%s=%s", k, v)
		venvs = append(venvs, keyEqualsValue)
	}

	sort.Strings(venvs)

	return venvs
}
