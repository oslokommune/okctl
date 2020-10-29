// Package virtualenv helps finding the environment variables needed to use a okctl cluster.
package virtualenv

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/helm"
	"github.com/oslokommune/okctl/pkg/okctl"
)

// GetVirtualEnvironmentOpts returns data needed to set up a virtual environment.
func GetVirtualEnvironmentOpts(o *okctl.Okctl) (VirtualEnvironmentOpts, error) {
	meta := o.RepoStateWithEnv.GetMetadata()
	cluster := o.RepoStateWithEnv.GetCluster()

	userDataDir, err := o.GetUserDataDir()
	if err != nil {
		return VirtualEnvironmentOpts{}, err
	}

	k, err := o.BinariesProvider.Kubectl(kubectl.Version)
	if err != nil {
		return VirtualEnvironmentOpts{}, err
	}

	a, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return VirtualEnvironmentOpts{}, err
	}

	opts := VirtualEnvironmentOpts{
		Region:                 meta.Region,
		AWSAccountID:           cluster.AWSAccountID,
		Environment:            cluster.Environment,
		Repository:             meta.Name,
		ClusterName:            cluster.Name,
		UserDataDir:            userDataDir,
		Debug:                  o.Debug,
		KubectlBinaryDir:       path.Dir(k.BinaryPath),
		AwsIamAuthenticatorDir: path.Dir(a.BinaryPath),
	}

	err = opts.validate()
	if err != nil {
		return VirtualEnvironmentOpts{}, errors.E(err, "failed to validate show credentials options")
	}

	return opts, nil
}

// GetVirtualEnvironment merges the passed virtual environment with the OS environment variables, and returns them as
// strings on the form "key=value".
func GetVirtualEnvironment(opts *VirtualEnvironmentOpts, osEnvVars []string) ([]string, error) {
	venv := getOkctlEnvVars(opts)
	osEnv := toMap(osEnvVars)

	osEnv["PATH"] = getPathWithOkctlBinaries(opts, osEnv)

	// TODO detect shell
	ps1, overridePs1 := os.LookupEnv("OKCTL_PS1")
	if overridePs1 {
		osEnv["PROMPT_COMMAND"] = ps1
	} else {
		osEnv["PROMPT_COMMAND"] = fmt.Sprintf(`PS1="\[\033[0;31m\]\w\[\033[0;34m\]\$(__git_ps1)\[\e[0m\] \[\033[0;32m\](\$(venv_ps1 %s)) \[\e[0m\]\$ "`, opts.Environment)
	}

	// Merge maps
	for osEnv, value := range osEnv {
		venv[osEnv] = value
	}

	return toEnvVarsSlice(&venv), nil
}

// Returns a map with environmental variables, where the map's key is the environment variable's name and map's value is
// the environment variable's value.
func getOkctlEnvVars(opts *VirtualEnvironmentOpts) map[string]string {
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

	for _, env := range slice {
		split := strings.SplitN(env, "=", 2)
		key := split[0]
		val := split[1]
		m[key] = val
	}

	return m
}

func getPathWithOkctlBinaries(opts *VirtualEnvironmentOpts, osEnv map[string]string) string {
	venvPath := createPs1ExecutableIfNotExists()
	okctlPath := fmt.Sprintf("%s:%s:%s", opts.KubectlBinaryDir, opts.AwsIamAuthenticatorDir, venvPath)
	osPath, hasKey := osEnv["PATH"]
	if hasKey {
		return fmt.Sprintf("%s:%s", okctlPath, osPath)
	} else {
		return okctlPath
	}
}

func createPs1ExecutableIfNotExists() string {
	return "/home/yngvar/.okctl/binaries/venv"
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
