package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"sigs.k8s.io/yaml"
)

// InferClusterFromStdinOrFile initializes a v1alpha1.Cluster based on a path. If the path is "-", the Cluster is
// initialized based on stdin
func InferClusterFromStdinOrFile(stdin io.Reader, path string, cluster *v1alpha1.Cluster) error {
	var (
		inputReader io.Reader
		err         error
	)

	switch path {
	case "-":
		inputReader = stdin
	default:
		inputReader, err = os.Open(filepath.Clean(path))
		if err != nil {
			return fmt.Errorf("unable to read file: %w", err)
		}
	}

	var (
		buffer bytes.Buffer
	)

	_, err = io.Copy(&buffer, inputReader)
	if err != nil {
		return fmt.Errorf("copying reader data: %w", err)
	}

	err = yaml.Unmarshal(buffer.Bytes(), &cluster)
	if err != nil {
		return fmt.Errorf("unmarshalling buffer: %w", err)
	}

	return nil
}

func ValidateClusterInput(cluster *v1alpha1.Cluster) error {
	// TODO: Can we do this? VPC will be empty for instance
	return cluster.Validate()
}
