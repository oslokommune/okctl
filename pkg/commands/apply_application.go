package commands

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"text/template"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/controller"
	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/spf13/afero"

	"sigs.k8s.io/yaml"
)

// SynchronizeApplicationOpts contains references necessary to synchronize an application
type SynchronizeApplicationOpts struct {
	ReconciliationManager reconciler.Reconciler
	Application           v1alpha1.Application

	Tree *resourcetree.ResourceNode
}

// SynchronizeApplication knows how to discover differences between desired and actual state and rectify them
func SynchronizeApplication(opts SynchronizeApplicationOpts) error {
	desiredResourceOpts := identifyDesiredResources(opts.Application.Image)

	opts.Tree.ApplyFunction(applyDesiredState(desiredResourceOpts), opts.Tree)

	return controller.HandleNode(opts.ReconciliationManager, opts.Tree)
}

type desiredResources struct {
	ContainerRepository bool
}

func identifyDesiredResources(image v1alpha1.ApplicationImage) desiredResources {
	return desiredResources{
		ContainerRepository: image.HasName(),
	}
}

func applyDesiredState(opts desiredResources) resourcetree.ApplyFn {
	return func(receiver *resourcetree.ResourceNode, target *resourcetree.ResourceNode) {
		switch receiver.Type {
		case resourcetree.ResourceNodeTypeContainerRepository:
			receiver.State = controller.BoolToState(opts.ContainerRepository)
		case resourcetree.ResourceNodeTypeApplication:
			receiver.State = resourcetree.ResourceNodeStatePresent
		}
	}
}

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
