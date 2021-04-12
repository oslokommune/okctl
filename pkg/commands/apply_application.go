package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"text/template"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/spf13/afero"

	"sigs.k8s.io/yaml"
)

// InferApplicationFromStdinOrFile returns an okctl application based on input. The function will parse input either
// from the reader or from the fs based on if path is a path or if it is "-". "-" represents stdin
func InferApplicationFromStdinOrFile(stdin io.Reader, fs *afero.Afero, path string) (client.OkctlApplication, error) {
	var (
		err         error
		app         client.OkctlApplication
		inputReader io.Reader
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

// WriteApplyApplicationSuccessMessage produces a relevant message for successfully reconciling an application
func WriteApplyApplicationSuccessMessage(writer io.Writer, applicationName, outputDir string) error {
	argoCDResourcePath := path.Join(
		outputDir,
		constant.DefaultApplicationsOutputDir,
		applicationName,
		"argocd-application.yaml",
	)

	templateString := `
	Successfully scaffolded {{ .ApplicationName }}
	To deploy your application:
		1. Commit and push the changes done by okctl
		2. Run kubectl apply -f {{ .ArgoCDResourcePath }}
	If using an ingress, it can take up to five minutes for the routing to configure
`

	tmpl, err := template.New("t").Parse(templateString)
	if err != nil {
		return err
	}

	var tmplBuffer bytes.Buffer

	err = tmpl.Execute(&tmplBuffer, struct {
		ApplicationName    string
		ArgoCDResourcePath string
	}{
		ApplicationName:    applicationName,
		ArgoCDResourcePath: argoCDResourcePath,
	})
	if err != nil {
		return err
	}

	fmt.Fprint(writer, tmplBuffer.String())

	return nil
}
