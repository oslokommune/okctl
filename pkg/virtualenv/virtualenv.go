// Package virtualenv helps finding the environment variables needed to use a okctl cluster.
package virtualenv

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/helm"
)

// VirtualEnvironment contains environment variables in a virtual environment.
type VirtualEnvironment struct {
	env map[string]string
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

// GetVirtualEnvironment merges the passed virtual environment with the OS environment variables, and returns them as
// strings on the form "key=value".
func GetVirtualEnvironment(opts VirtualEnvironmentOpts) (*VirtualEnvironment, error) {
	osEnv := toMap(opts.OsEnvVars)
	venv := getOkctlEnvVars(opts)

	addOkctlBinariesToPath(opts, osEnv)

	// Merge maps
	for osEnv, value := range osEnv {
		venv[osEnv] = value
	}

	return &VirtualEnvironment{
		env: venv,
	}, nil
}

// Returns a map with environmental variables, where the map's key is the environment variable's name and map's value is
// the environment variable's value.
func getOkctlEnvVars(opts VirtualEnvironmentOpts) map[string]string {
	appDir := opts.UserDataDir

	kubeConfig := path.Join(appDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterKubeConfig)
	awsConfig := path.Join(appDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterAwsConfig)
	awsCredentials := path.Join(appDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterAwsCredentials)

	h := &helm.Config{
		HelmPluginsDirectory: path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmPluginsDirectory),
		HelmRegistryConfig:   path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmRegistryConfig),
		HelmRepositoryCache:  path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryCache),
		HelmRepositoryConfig: path.Join(appDir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryConfig),
		HelmBaseDir:          path.Join(appDir, config.DefaultHelmBaseDir),
		Debug:                opts.Debug,
	}

	envMap := make(map[string]string)

	for k, v := range h.Envs() {
		if k == "HOME" || k == "PATH" {
			continue
		}

		envMap[k] = v
	}

	envMap["AWS_CONFIG_FILE"] = awsConfig
	envMap["AWS_SHARED_CREDENTIALS_FILE"] = awsCredentials
	envMap["AWS_PROFILE"] = "default"
	envMap["KUBECONFIG"] = kubeConfig

	return envMap
}

func toMap(slice []string) map[string]string {
	m := make(map[string]string)
77
	for _, env := range slice {
		split := strings.SplitN(env, "=", 2)
		key := split[0]
		val := split[1]
		m[key] = val
	}

	return m
}

func addOkctlBinariesToPath(opts VirtualEnvironmentOpts, osEnv map[string]string) {
	okctlPath := calcOkctlPath(opts)

	osPath, osPathExists := osEnv["PATH"]

	if osPathExists {
		osEnv["PATH"] = fmt.Sprintf("%s:%s", okctlPath, osPath)
	} else {
		osEnv["PATH"] = okctlPath
	}
}

func calcOkctlPath(opts VirtualEnvironmentOpts) string {
	// TODO don't add ps1 here, should be for later!
	var okctlPath string
	if opts.Ps1Dir == "" {
		okctlPath = fmt.Sprintf("%s:%s", opts.KubectlBinaryDir, opts.AwsIamAuthenticatorDir)
	} else {
		okctlPath = fmt.Sprintf("%s:%s:%s", opts.KubectlBinaryDir, opts.AwsIamAuthenticatorDir, opts.Ps1Dir)
	}

	return okctlPath
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
