package binary

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"

	"github.com/oslokommune/okctl/pkg/logging"
)

func (c client) applyFile(manifestPath string) error {
	log := logging.GetLogger("kubectl/binary", "applyFile")

	k, err := c.binaryProvider.Kubectl(kubectl.Version)
	if err != nil {
		return fmt.Errorf("acquiring kubectl binary path: %w", err)
	}

	args := []string{
		"apply", "-f", manifestPath,
	}

	cmd := exec.Command(k.BinaryPath, args...) //nolint:gosec

	multiwriter := bytes.Buffer{}

	cmd.Stdout = &multiwriter
	cmd.Stderr = &multiwriter

	cmd.Env, err = c.generateEnv()
	if err != nil {
		return fmt.Errorf("generating environment: %w", err)
	}

	err = cmd.Run()
	if err != nil {
		log.Error(multiwriter.String())

		return fmt.Errorf("running command: %w", err)
	}

	return nil
}
