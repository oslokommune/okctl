package binary

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/oslokommune/okctl/pkg/logging"

	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/google/uuid"
)

// runKubectlCommand knows how to run a command with the kubectl binary
func (c client) runKubectlCommand(logActivity string, args []string) error {
	log := logging.GetLogger("kubectl/binary", logActivity)

	k, err := c.binaryProvider.Kubectl(kubectl.Version)
	if err != nil {
		return fmt.Errorf("acquiring kubectl binary path: %w", err)
	}

	cmd := exec.Command(k.BinaryPath, args...) //nolint:gosec
	log.Debug(fmt.Sprintf("passing arguments: %+v", args))

	multiwriter := bytes.Buffer{}

	cmd.Stdout = &multiwriter
	cmd.Stderr = &multiwriter

	cmd.Env, err = c.generateEnv()
	if err != nil {
		return fmt.Errorf("generating environment: %w", err)
	}

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("%s: %w", multiwriter.String(), err)

		return errorHandler(err, fmt.Errorf("running command: %w", err))
	}

	return nil
}

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

// generateEnv knows how to produce an array of environment variables required to run kubectl commands
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
		"HOME": homeDir,
	})

	env = append(env, awsEnvCredentials...)

	return env, nil
}

// generatePath knows how to produce the value of a PATH environment variable containing all necessary binaries to run
// kubectl commands
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
