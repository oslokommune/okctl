package commands

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/oslokommune/okctl/pkg/context"

	"github.com/oslokommune/okctl/pkg/config/constant"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubens"
	"github.com/oslokommune/okctl/pkg/helm"
	"github.com/oslokommune/okctl/pkg/okctl"
)

// OkctlEnvironment contains data about the okctl environment
type OkctlEnvironment struct {
	AWSAccountID           string
	Region                 string
	AwsProfile             string
	ClusterName            string
	UserDataDir            string
	UserHomeDir            string
	Debug                  bool
	KubectlBinaryDir       string
	KubensBinaryDir        string
	AwsIamAuthenticatorDir string
	ClusterDeclarationPath string
	AWSCredentialsType     string
	GithubCredentialsType  string
}

// Validate the inputs
func (o *OkctlEnvironment) Validate() error {
	validators := []*validation.FieldRules{
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
		validation.Field(&o.UserDataDir, validation.Required),
		validation.Field(&o.UserHomeDir, validation.Required),
		validation.Field(&o.KubectlBinaryDir, validation.Required),
		validation.Field(&o.KubensBinaryDir, validation.Required),
		validation.Field(&o.AwsIamAuthenticatorDir, validation.Required),
		validation.Field(&o.ClusterDeclarationPath, validation.Required),
		validation.Field(&o.AWSCredentialsType, validation.Required),
		validation.Field(&o.GithubCredentialsType, validation.Required),
	}

	if o.AWSCredentialsType == context.AWSCredentialsTypeAwsProfile {
		validators = append(validators, validation.Field(&o.AwsProfile, validation.Required))
	}

	return validation.ValidateStruct(o, validators...)
}

// GetOkctlEnvironment returns data needed to connect to an okctl cluster
func GetOkctlEnvironment(o *okctl.Okctl, clusterDeclarationPath string) (OkctlEnvironment, error) {
	userDataDir, err := o.GetUserDataDir()
	if err != nil {
		return OkctlEnvironment{}, err
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return OkctlEnvironment{}, err
	}

	awsProfile := os.Getenv("AWS_PROFILE")

	k, err := o.BinariesProvider.Kubectl(kubectl.Version)
	if err != nil {
		return OkctlEnvironment{}, err
	}

	kn, err := o.BinariesProvider.Kubens(kubens.Version)
	if err != nil {
		return OkctlEnvironment{}, err
	}

	a, err := o.BinariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return OkctlEnvironment{}, err
	}

	absoluteClusterDeclarationPath, err := ensureAbsolutePath(clusterDeclarationPath)
	if err != nil {
		return OkctlEnvironment{}, fmt.Errorf("ensuring absolute declaration path: %w", err)
	}

	opts := OkctlEnvironment{
		AWSAccountID:           o.Declaration.Metadata.AccountID,
		Region:                 o.Declaration.Metadata.Region,
		AwsProfile:             awsProfile,
		ClusterName:            o.Declaration.Metadata.Name,
		UserDataDir:            userDataDir,
		UserHomeDir:            userHomeDir,
		Debug:                  o.Debug,
		KubectlBinaryDir:       path.Dir(k.BinaryPath),
		KubensBinaryDir:        path.Dir(kn.BinaryPath),
		AwsIamAuthenticatorDir: path.Dir(a.BinaryPath),
		ClusterDeclarationPath: absoluteClusterDeclarationPath,
		AWSCredentialsType:     o.Context.AWSCredentialsType,
		GithubCredentialsType:  o.Context.GithubCredentialsType,
	}

	err = opts.Validate()
	if err != nil {
		return OkctlEnvironment{}, fmt.Errorf("failed to validate okctl environment: %w", err)
	}

	return opts, nil
}

func ensureAbsolutePath(declarationPath string) (string, error) {
	if path.IsAbs(declarationPath) {
		return declarationPath, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting current directory: %w", err)
	}

	return path.Join(cwd, declarationPath), nil
}

// GetVenvEnvVars returns the environmental variables needed by a virtual environment. This contains environment variables from
// the user's shell merged with those from the okctl.
func GetVenvEnvVars(okctlEnvironment OkctlEnvironment) (map[string]string, error) {
	okctlEnvVars, err := GetOkctlEnvVars(okctlEnvironment)
	if err != nil {
		return nil, err
	}

	return MergeEnvVars(CleanOsEnvVars(os.Environ()), okctlEnvVars), nil
}

// GetOkctlEnvVars converts an okctl environment to a map with environmental variables
func GetOkctlEnvVars(opts OkctlEnvironment) (map[string]string, error) {
	appDir := opts.UserDataDir

	var awsProfile, awsConfig, awsCredentials string

	if opts.AWSCredentialsType == context.AWSCredentialsTypeAwsProfile {
		awsProfile = opts.AwsProfile
		if awsProfile == "" {
			return nil, fmt.Errorf("environment variable AWS_PROFILE not set")
		}

		awsConfig = path.Join(opts.UserHomeDir, ".aws", "config")
		awsCredentials = path.Join(opts.UserHomeDir, ".aws", "credentials")
	} else {
		awsConfig = path.Join(appDir, constant.DefaultCredentialsDirName, opts.ClusterName, constant.DefaultClusterAwsConfig)
		awsCredentials = path.Join(appDir, constant.DefaultCredentialsDirName, opts.ClusterName, constant.DefaultClusterAwsCredentials)
	}

	kubeConfig := path.Join(appDir, constant.DefaultCredentialsDirName, opts.ClusterName, constant.DefaultClusterKubeConfig)

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
	envMap["AWS_PROFILE"] = awsProfile
	envMap["KUBECONFIG"] = kubeConfig
	envMap["PATH"] = getPathWithOkctlBinaries(opts)

	envMap[constant.EnvClusterDeclaration] = opts.ClusterDeclarationPath
	envMap[context.DefaultAWSCredentialsType] = opts.AWSCredentialsType
	envMap[context.DefaultGithubCredentialsType] = opts.GithubCredentialsType

	return envMap, nil
}

func getPathWithOkctlBinaries(opts OkctlEnvironment) string {
	okctlPath := fmt.Sprintf("%s:%s:%s", opts.KubectlBinaryDir, opts.KubensBinaryDir, opts.AwsIamAuthenticatorDir)
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
	numberOfResultingSubstrings := 2

	for _, env := range slice {
		split := strings.SplitN(env, "=", numberOfResultingSubstrings)
		key := split[0]
		val := split[1]
		m[key] = val
	}

	return m
}

func toSlice(m map[string]string) []string {
	result := make([]string, len(m))
	index := 0

	for key, val := range m {
		result[index] = fmt.Sprintf("%s=%s", key, val)

		index++
	}

	return result
}

// CleanOsEnvVars ensures blacklisted variables are removed from the list
func CleanOsEnvVars(environ []string) []string {
	// We want to control AWS_PROFILE for e.g. SAML login
	keyBlacklist := []string{constant.EnvClusterDeclaration, "AWS_PROFILE"}
	cleanedVars := toMap(environ)

	for _, blacklistedKey := range keyBlacklist {
		delete(cleanedVars, blacklistedKey)
	}

	return toSlice(cleanedVars)
}
