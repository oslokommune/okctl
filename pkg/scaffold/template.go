package scaffold

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"

	"github.com/oslokommune/kaex/pkg/api"
)

// FetchExample downloads an example file and writes it to a buffer
func FetchExample(fullExample bool) ([]byte, error) {
	var (
		err    error
		buffer bytes.Buffer
	)

	if fullExample {
		err = api.FetchFullExample(&buffer)
	} else {
		err = api.FetchMinimalExample(&buffer)
	}

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
