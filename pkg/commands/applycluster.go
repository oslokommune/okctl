package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
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
			return nil, fmt.Errorf(constant.ReadFileError, err)
		}
	}

	var (
		buffer  bytes.Buffer
		cluster v1alpha1.Cluster
	)

	cluster = v1alpha1.NewCluster()

	_, err = io.Copy(&buffer, inputReader)
	if err != nil {
		return nil, fmt.Errorf(constant.CopyReaderDataError, err)
	}

	err = yaml.Unmarshal(buffer.Bytes(), &cluster)
	if err != nil {
		return nil, fmt.Errorf(constant.UnmarshalBufferError, err)
	}

	return &cluster, nil
}
