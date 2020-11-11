package commands

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/helm"
	"github.com/oslokommune/okctl/pkg/okctl"
	"os"
	"path"
	"strings"
)

type CredentialsOpts struct {
	Region                 string
	AWSAccountID           string
	Environment            string
	Repository             string
	ClusterName            string
	UserDataDir            string
	Debug                  bool
	KubectlBinaryDir       string
	AwsIamAuthenticatorDir string
}

// Validate the inputs
func (o *CredentialsOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Environment, validation.Required),
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
		validation.Field(&o.UserDataDir, validation.Required),
		validation.Field(&o.KubectlBinaryDir, validation.Required),
		validation.Field(&o.AwsIamAuthenticatorDir, validation.Required),
	)
}

// GetCredentialsOpts returns data needed to connect to an okctl cluster
func GetCredentialsOpts(o *okctl.Okctl) (CredentialsOpts, error) {
	meta := o.RepoStateWithEnv.GetMetadata()
	cluster := o.RepoStateWithEnv.GetCluster()

	userDataDir, err := o.GetUserDataDir()
	if err != nil {
		return CredentialsOpts{}, err
	}

	k, err := o.BinariesProvider.Kubectl(kubectl.Version)
	if err != nil {
		return CredentialsOpts{}, err
	}

	a, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return CredentialsOpts{}, err
	}

	opts := CredentialsOpts{
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

	err = opts.Validate()
	if err != nil {
		return CredentialsOpts{}, errors.E(err, "failed to validate show credentials options")
	}

	return opts, nil
}

// Returns a map with environmental variables, where the map's key is the environment variable's name and map's value is
// the environment variable's value.
func GetOkctlEnvVars(opts CredentialsOpts) map[string]string {
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
	envMap["PATH"] = getPathWithOkctlBinaries(opts)

	return envMap
}

func getPathWithOkctlBinaries(opts CredentialsOpts) string {
	okctlPath := fmt.Sprintf("%s:%s", opts.KubectlBinaryDir, opts.AwsIamAuthenticatorDir)
	osPath, osPathExists := os.LookupEnv("PATH")

	if osPathExists {
		return fmt.Sprintf("%s:%s", okctlPath, osPath)
	} else {
		return okctlPath
	}
}

func MergeEnvVars(osEnvs []string, venvMap map[string]string) map[string]string {
	merged := make(map[string]string)

	for key, val := range venvMap {
		merged[key] = val
	}

	osEnvMap := toMap(osEnvs)
	for key, val := range osEnvMap {
		merged[key] = val
	}

	// Merge PATHs
	venvPath, venvHasPath := venvMap["PATH"]
	osPath, osHasPath := osEnvMap["PATH"]
	if osHasPath && venvHasPath {
		merged["PATH"] = fmt.Sprintf("%s:%s", venvPath, osPath)
	}

	return merged
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
