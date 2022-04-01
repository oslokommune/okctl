package binary

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/oslokommune/okctl/pkg/logging"
)

// runCommand knows how to run a command using the binary
func (c client) runCommand(logActivity string, baseDir string, args []string) error {
	log := logging.GetLogger("terraform/binary", logActivity)

	k, err := c.binaryProvider.Terraform(c.version)
	if err != nil {
		return fmt.Errorf("acquiring kubectl binary path: %w", err)
	}

	cmd := exec.Command(k.BinaryPath, args...) //nolint:gosec
	log.Debug(fmt.Sprintf("passing arguments: %+v", args))

	multiwriter := bytes.Buffer{}

	cmd.Stdout = os.Stdout // Should probably not assume os.Stdout?
	cmd.Stderr = &multiwriter

	cmd.Env, err = c.generateEnv()
	if err != nil {
		return fmt.Errorf("generating environment: %w", err)
	}

	cmd.Dir = baseDir

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %w", multiwriter.String(), err)
	}

	return nil
}

// generateEnv knows how to produce an array of environment variables required to run terraform commands
func (c client) generateEnv() ([]string, error) {
	awsEnvCredentials, err := c.credentialsProvider.Aws().AsEnv()
	if err != nil {
		return []string{}, fmt.Errorf("acquiring AWS credentials: %w", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []string{}, fmt.Errorf("acquiring user home directory: %w", err)
	}

	env := envAsArray(map[string]string{
		"HOME": homeDir,
	})

	env = append(env, awsEnvCredentials...)

	return env, nil
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
