package scaffold

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/oslokommune/okctl/pkg/storage"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"

	kaex "github.com/oslokommune/kaex/pkg/api"
)

// FetchTemplate downloads an example file and writes it to a buffer
func FetchTemplate(kx kaex.Kaex) ([]byte, error) {
	var (
		err    error
		buffer bytes.Buffer
	)

	err = kaex.FetchTemplate(kx, &buffer, "application")

	if err != nil {
		return nil, fmt.Errorf("unable to fetch example: %w", err)
	}

	return buffer.Bytes(), nil
}

// InterpolateTemplate replaces dummy data in the template with state dependant data
func InterpolateTemplate(o *okctl.Okctl, cmd *cobra.Command, env string, template []byte) ([]byte, error) {
	cluster := GetCluster(o, cmd, env)
	domain := GetHostedZoneDomain(cluster)

	var outputBuffer bytes.Buffer

	output := strings.Replace(
		string(template),
		"my-domain.io",
		fmt.Sprintf("<app-name>.%s", domain),
		1,
	)

	_, err := io.Copy(&outputBuffer, bytes.NewBufferString(output))
	if err != nil {
		return nil, fmt.Errorf("error writing to output buffer: %w", err)
	}

	return outputBuffer.Bytes(), nil
}

// SaveTemplate saves a byte array as an application.yaml file in the current directory
func SaveTemplate(template []byte) error {
	cwd, _ := os.Getwd()
	templateStorage := storage.NewFileSystemStorage(cwd)

	applicationFile, err := templateStorage.Create("", "application.yaml", 0o644)
	if err != nil {
		return fmt.Errorf("error creating application.yaml: %w", err)
	}

	_, err = applicationFile.Write(template)
	if err != nil {
		return fmt.Errorf("error writing to application.yaml: %w", err)
	}

	err = applicationFile.Close()
	if err != nil {
		return fmt.Errorf("unable to close application.yaml after writing: %w", err)
	}

	return err
}
