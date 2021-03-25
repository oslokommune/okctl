package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/oslokommune/okctl/pkg/client"

	"sigs.k8s.io/yaml"
)

// InferApplicationFromStdinOrFile initializes a okctl application based on a path. If the path is "-", the application
// is initialized based on stdin
func InferApplicationFromStdinOrFile(stdin io.Reader, path string) (client.OkctlApplication, error) {
	var (
		inputReader io.Reader
		err         error
		buf         []byte
		app         client.OkctlApplication
	)

	switch path {
	case "-":
		inputReader = stdin
	default:
		inputReader, err = os.Open(filepath.Clean(path))
		if err != nil {
			return app, fmt.Errorf("unable to read file: %w", err)
		}
	}

	buf, err = io.ReadAll(inputReader)
	if err != nil {
		return app, fmt.Errorf("reading application data: %w", err)
	}

	err = yaml.Unmarshal(buf, &app)
	if err != nil {
		return app, fmt.Errorf("unmarshalling buffer: %w", err)
	}

	return app, nil
}
