package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/spf13/afero"

	"sigs.k8s.io/yaml"
)

const pushECRImageInstructionURL = "https://www.okctl.io/running-a-docker-image-in-your-cluster/#push-a-docker-image-to-the-amazon-elastic-container-registry-ecr"

// InferApplicationFromStdinOrFile returns an okctl application based on input. The function will parse input either
// from the reader or from the fs based on if path is a path or if it is "-". "-" represents stdin
func InferApplicationFromStdinOrFile(declaration v1alpha1.Cluster, stdin io.Reader, fs *afero.Afero, path string) (v1alpha1.Application, error) {
	var (
		err         error
		inputReader io.Reader
		app         = v1alpha1.NewApplication(declaration)
	)

	switch path {
	case "-":
		inputReader = stdin
	default:
		inputReader, err = fs.Open(filepath.Clean(path))
		if err != nil {
			return app, fmt.Errorf("opening application file: %w", err)
		}
	}

	var buf []byte

	buf, err = ioutil.ReadAll(inputReader)
	if err != nil {
		return app, fmt.Errorf("reading application file: %w", err)
	}

	err = yaml.Unmarshal(buf, &app)
	if err != nil {
		return app, fmt.Errorf("parsing application yaml: %w", err)
	}

	return app, nil
}

// ApplyApplicationSuccessMessageOpts contains the values for customizing the apply application success message
type ApplyApplicationSuccessMessageOpts struct {
	ApplicationName           string
	OptionalDockerTagPushStep string
	OptionalDockerImageURI    string
	OptionalIngressInfo       string
}

const applyApplicationSuccessMessage = `
Successfully applied {{ .ApplicationName }}

To finalize the changes:
- Commit and push the changes done by okctl{{ .OptionalDockerTagPushStep }}
{{ .OptionalIngressInfo }}
`

// WriteApplyApplicationSucessMessageOpts contains necessary information to compile and write a success message
type WriteApplyApplicationSucessMessageOpts struct {
	Out io.Writer

	Application v1alpha1.Application
	Cluster     v1alpha1.Cluster
}

// WriteApplyApplicationSuccessMessage produces a relevant message for successfully reconciling an application
func WriteApplyApplicationSuccessMessage(opts WriteApplyApplicationSucessMessageOpts) error {
	optionalIngressInfo := ""
	optionalDockerTagPushStep := ""

	if opts.Application.Image.HasName() {
		optionalDockerTagPushStep = fmt.Sprintf("\n%s\n  %s",
			"- Tag and push a docker image to your container repository. See instructions on",
			pushECRImageInstructionURL,
		)
	}

	if opts.Application.HasIngress() {
		optionalIngressInfo = fmt.Sprintf("\n%s",
			"N.B.: it can take up to five minutes for the routing to configure",
		)
	}

	tmpl, err := template.New("t").Parse(applyApplicationSuccessMessage)
	if err != nil {
		return err
	}

	var tmplBuffer bytes.Buffer

	err = tmpl.Execute(&tmplBuffer, ApplyApplicationSuccessMessageOpts{
		ApplicationName:           opts.Application.Metadata.Name,
		OptionalDockerTagPushStep: optionalDockerTagPushStep,
		OptionalDockerImageURI:    opts.Application.Image.URI,
		OptionalIngressInfo:       optionalIngressInfo,
	})
	if err != nil {
		return err
	}

	fmt.Fprint(opts.Out, tmplBuffer.String())

	return nil
}
