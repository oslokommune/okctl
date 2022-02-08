package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"text/template"

	"github.com/logrusorgru/aurora/v3"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/spf13/afero"

	"sigs.k8s.io/yaml"
)

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

func generateArgoCDApplicationManifestPath(cluster v1alpha1.Cluster, application v1alpha1.Application) string {
	return path.Join(
		cluster.Github.OutputPath,
		constant.DefaultApplicationsOutputDir,
		application.Metadata.Name,
		constant.DefaultApplicationOverlayDir,
		cluster.Metadata.Name,
		"argocd-application.yaml",
	)
}

// ApplyApplicationSuccessMessageOpts contains the values for customizing the apply application success message
type ApplyApplicationSuccessMessageOpts struct {
	ApplicationName           string
	OptionalDockerTagPushStep string
	OptionalDockerImageURI    string
	KubectlApplyArgoCmd       string
}

const applyApplicationSuccessMessage = `
	Successfully scaffolded {{ .ApplicationName }}
	To deploy your application:
		- Commit and push the changes done by okctl{{ .OptionalDockerTagPushStep }}
		- Run {{ .KubectlApplyArgoCmd }}

    If using an ingress, it can take up to five minutes for the routing to configure
`

// WriteApplyApplicationSucessMessageOpts contains necessary information to compile and write a success message
type WriteApplyApplicationSucessMessageOpts struct {
	Out io.Writer

	Application v1alpha1.Application
	Cluster     v1alpha1.Cluster
}

// WriteApplyApplicationSuccessMessage produces a relevant message for successfully reconciling an application
func WriteApplyApplicationSuccessMessage(opts WriteApplyApplicationSucessMessageOpts) error {
	argoCDResourcePath := generateArgoCDApplicationManifestPath(opts.Cluster, opts.Application)

	optionalDockerTagPushStep := ""

	if opts.Application.Image.HasName() {
		optionalDockerTagPushStep = `
        - Tag and push a docker image to your container repository. See instructions on
          https://okctl.io/help/docker-registry/#push-a-docker-image-to-the-amazon-elastic-container-registry-ecr`
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
		KubectlApplyArgoCmd:       aurora.Green(fmt.Sprintf("kubectl apply -f %s", argoCDResourcePath)).String(),
	})
	if err != nil {
		return err
	}

	fmt.Fprint(opts.Out, tmplBuffer.String())

	return nil
}

const deleteApplicationSuccessMessage = `
	Successfully deleted {{ .ApplicationName }}
	To finish the deletion process:
		- Commit and push the changes done by okctl
		- Run {{ .KubectlDeleteArgoCDManifestCmd }}
`

// WriteDeleteApplicationSuccessMessageOpts contains the values for customizing delete application success message
type WriteDeleteApplicationSuccessMessageOpts struct {
	Out io.Writer

	Cluster     v1alpha1.Cluster
	Application v1alpha1.Application
}

// WriteDeleteApplicationSuccessMessage produces a relevant message for successfully deleting an application
func WriteDeleteApplicationSuccessMessage(opts WriteDeleteApplicationSuccessMessageOpts) error {
	argoCDResourcePath := generateArgoCDApplicationManifestPath(opts.Cluster, opts.Application)

	tmpl, err := template.New("t").Parse(deleteApplicationSuccessMessage)
	if err != nil {
		return err
	}

	err = tmpl.Execute(opts.Out, struct {
		ApplicationName                string
		KubectlDeleteArgoCDManifestCmd string
	}{
		ApplicationName:                opts.Application.Metadata.Name,
		KubectlDeleteArgoCDManifestCmd: aurora.Green(fmt.Sprintf("kubectl delete -f %s", argoCDResourcePath)).String(),
	})
	if err != nil {
		return err
	}

	return nil
}
