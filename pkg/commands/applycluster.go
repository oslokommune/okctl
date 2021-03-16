package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"sigs.k8s.io/yaml"
)

// InferClusterFromStdinOrFile initializes a v1alpha1.Cluster based on a path. If the path is "-", the Cluster is
// initialized based on stdin
func InferClusterFromStdinOrFile(stdin io.Reader, path string) (*v1alpha1.Cluster, error) {
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
			return nil, fmt.Errorf("unable to read file: %w", err)
		}
	}

	var (
		buffer  bytes.Buffer
		cluster v1alpha1.Cluster
	)

	cluster = v1alpha1.NewDefaultCluster(
		"",
		"",
		constant.DefaultGithubOrganization,
		"",
		"",
	)

	_, err = io.Copy(&buffer, inputReader)
	if err != nil {
		return nil, fmt.Errorf("copying reader data: %w", err)
	}

	err = yaml.Unmarshal(buffer.Bytes(), &cluster)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling buffer: %w", err)
	}

	return &cluster, nil
}
