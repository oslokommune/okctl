package scaffold

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/afero"

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

// InterpolationOpts defines possible data to inject into the templates
type InterpolationOpts struct {
	Domain string
}

/*
InterpolateTemplate replaces dummy data in the template with state dependant data

Parameters:
template []byte: the template in which we should do the interpolation
opts *InterpolationOpts: What values to interpolate with what
*/
func InterpolateTemplate(template []byte, opts *InterpolationOpts) (interpolatedTemplate []byte, err error) {
	var outputBuffer bytes.Buffer

	output := strings.Replace(
		string(template),
		"my-domain.io",
		fmt.Sprintf("<app-name>.%s", opts.Domain),
		1,
	)

	_, err = io.Copy(&outputBuffer, bytes.NewBufferString(output))
	if err != nil {
		return nil, fmt.Errorf("error writing to output buffer: %w", err)
	}

	return outputBuffer.Bytes(), nil
}

// SaveTemplate saves a byte array as an application.yaml file in the current directory
func SaveTemplate(fs *afero.Afero, path string, template []byte) error {
	applicationFile, err := fs.Create(path)
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
