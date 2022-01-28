package binary

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/google/uuid"
)

// envAsArray converts a map to a string array of KEY=VALUE pairs
func envAsArray(m map[string]string) []string {
	result := make([]string, len(m))
	index := 0

	for key, value := range m {
		result[index] = fmt.Sprintf("%s=%s", key, value)

		index++
	}

	return result
}

// cacheReaderOnFs writes f to a temporary file, then returns the path and a clean up function
func (c *client) cacheReaderOnFs(r io.Reader) (string, teardownFn, error) {
	dir, err := c.fs.TempDir("", "okctl")
	if err != nil {
		return "", nil, fmt.Errorf("opening temporary file: %w", err)
	}

	targetPath := path.Join(dir, uuid.New().String())

	err = c.fs.WriteReader(targetPath, r)
	if err != nil {
		return "", nil, fmt.Errorf("writing file: %w", err)
	}

	teardowner := func() error {
		return c.fs.RemoveAll(dir)
	}

	return targetPath, teardowner, nil
}

func (c client) generateEnv() ([]string, error) {
	awsEnvCredentials, err := c.credentialsProvider.Aws().AsEnv()
	if err != nil {
		return []string{}, fmt.Errorf("acquiring AWS credentials: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []string{}, fmt.Errorf("acquiring user home directory: %w", err)
	}

	generatedPaths, err := c.generatePath()
	if err != nil {
		return []string{}, fmt.Errorf("generating PATH: %w", err)
	}

	env := envAsArray(map[string]string{
		"KUBECONFIG": path.Join(homeDir,
			constant.DefaultDir,
			constant.DefaultCredentialsDirName,
			c.cluster.Metadata.Name,
			constant.DefaultClusterKubeConfig,
		),
		"PATH": strings.Join(generatedPaths, ":"),
	})

	env = append(env, awsEnvCredentials...)

	return env, nil
}

func (c client) generatePath() ([]string, error) {
	k, err := c.binaryProvider.Kubectl(kubectl.Version)
	if err != nil {
		return []string{}, fmt.Errorf("acquiring kubectl binary path: %w", err)
	}

	a, err := c.binaryProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return []string{}, fmt.Errorf("acquiring AWS IAM authenticator binary path: %w", err)
	}

	return []string{
		path.Dir(k.BinaryPath),
		path.Dir(a.BinaryPath),
	}, nil
}
