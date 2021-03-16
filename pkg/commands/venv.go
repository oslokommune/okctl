package commands

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/oslokommune/okctl/pkg/config/constant"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/helm"
	"github.com/oslokommune/okctl/pkg/okctl"
)

// OkctlEnvironment contains data about the okctl environment
type OkctlEnvironment struct {
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
func (o *OkctlEnvironment) Validate() error {
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

// GetOkctlEnvironment returns data needed to connect to an okctl cluster
func GetOkctlEnvironment(o *okctl.Okctl) (OkctlEnvironment, error) {
	meta := o.RepoStateWithEnv.GetMetadata()
	cluster := o.RepoStateWithEnv.GetCluster()

	userDataDir, err := o.GetUserDataDir()
	if err != nil {
		return OkctlEnvironment{}, err
	}

	k, err := o.BinariesProvider.Kubectl(kubectl.Version)
	if err != nil {
		return OkctlEnvironment{}, err
	}

	a, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return OkctlEnvironment{}, err
	}

	opts := OkctlEnvironment{
		Region:                 meta.Region,
		AWSAccountID:           cluster.AWSAccountID,
		Environment:            cluster.Environment,
		Repository:             meta.Name,
		ClusterName:            o.RepoStateWithEnv.GetClusterName(),
		UserDataDir:            userDataDir,
		Debug:                  o.Debug,
		KubectlBinaryDir:       path.Dir(k.BinaryPath),
		AwsIamAuthenticatorDir: path.Dir(a.BinaryPath),
	}

	err = opts.Validate()
	if err != nil {
		return OkctlEnvironment{}, fmt.Errorf("failed to validate okctl environment: %w", err)
	}

	return opts, nil
}

// GetOkctlEnvVars converts an okctl environment to a map with environmental variables
func GetOkctlEnvVars(opts OkctlEnvironment) map[string]string {
	appDir := opts.UserDataDir

	kubeConfig := path.Join(appDir, constant.DefaultCredentialsDirName, opts.ClusterName, constant.DefaultClusterKubeConfig)
	awsConfig := path.Join(appDir, constant.DefaultCredentialsDirName, opts.ClusterName, constant.DefaultClusterAwsConfig)
	awsCredentials := path.Join(appDir, constant.DefaultCredentialsDirName, opts.ClusterName, constant.DefaultClusterAwsCredentials)

	h := &helm.Config{
		HelmPluginsDirectory: path.Join(appDir, constant.DefaultHelmBaseDir, constant.DefaultHelmPluginsDirectory),
		HelmRegistryConfig:   path.Join(appDir, constant.DefaultHelmBaseDir, constant.DefaultHelmRegistryConfig),
		HelmRepositoryCache:  path.Join(appDir, constant.DefaultHelmBaseDir, constant.DefaultHelmRepositoryCache),
		HelmRepositoryConfig: path.Join(appDir, constant.DefaultHelmBaseDir, constant.DefaultHelmRepositoryConfig),
		HelmBaseDir:          path.Join(appDir, constant.DefaultHelmBaseDir),
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

func getPathWithOkctlBinaries(opts OkctlEnvironment) string {
	okctlPath := fmt.Sprintf("%s:%s", opts.KubectlBinaryDir, opts.AwsIamAuthenticatorDir)
	osPath, osPathExists := os.LookupEnv("PATH")

	if osPathExists {
		return fmt.Sprintf("%s:%s", okctlPath, osPath)
	}

	return okctlPath
}

// MergeEnvVars first converts the given slice to a map. The provided slice must contain strings on the form "KEY=VALUE.
// It then merges this map with the other provided map.
// If both map contains a PATH key, they will be merged.
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
