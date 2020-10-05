// Package scaffold knows how to scaffold okctl applications
package scaffold

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/oslokommune/okctl/internal/third_party/argoproj/argo-cd/pkg/apis/application/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1beta1"

	"github.com/oslokommune/kaex/pkg/api"
)

// ApplicationDeployment contains necessary data for a deployment
type ApplicationDeployment struct {
	ArgoApplication v1alpha1.Application
	Deployment      appsv1.Deployment
	Ingress         networkingv1.Ingress
	Service         v1.Service
	Volumes         []v1.PersistentVolumeClaim
}

func (deployment *ApplicationDeployment) Write(writer io.Writer) error {
	err := deployment.WriteKubernetesResources(writer)
	if err != nil {
		return err
	}

	err = deployment.WriteArgoResources(writer)
	if err != nil {
		return err
	}

	return nil
}

// WriteKubernetesResources writes kubernetes resources to stream as yaml
func (deployment *ApplicationDeployment) WriteKubernetesResources(writer io.Writer) error {
	for index := range deployment.Volumes {
		err := api.WriteResource(writer, deployment.Volumes[index])
		if err != nil {
			return fmt.Errorf("error writing volume to buffer: %w", err)
		}
	}

	if deployment.Service.Name != "" {
		err := api.WriteResource(writer, deployment.Service)
		if err != nil {
			return fmt.Errorf("error writing service to buffer: %w", err)
		}
	}

	if deployment.Ingress.Name != "" {
		err := api.WriteResource(writer, deployment.Ingress)
		if err != nil {
			return fmt.Errorf("error writing ingress to buffer: %w", err)
		}
	}

	err := api.WriteResource(writer, deployment.Deployment)
	if err != nil {
		return fmt.Errorf("error writing deployment to buffer: %w", err)
	}

	return nil
}

// WriteArgoResources writes Argo-cd resources to stream as yaml
func (deployment *ApplicationDeployment) WriteArgoResources(writer io.Writer) error {
	err := api.WriteResource(writer, deployment.ArgoApplication)
	if err != nil {
		return fmt.Errorf("error writing ArgoApp to buffer: %w", err)
	}

	return nil
}

// NewApplicationDeployment converts a Kaex Application to an okctl deployment
func NewApplicationDeployment(app api.Application, iacRepoURL string) (*ApplicationDeployment, error) {
	applicationDeployment := ApplicationDeployment{}

	for index := range app.Volumes {
		applicationDeployment.Volumes = make([]v1.PersistentVolumeClaim, len(app.Volumes))

		pvc, err := createOkctlVolume(app, app.Volumes[index])
		if err != nil {
			return nil, fmt.Errorf("unable to create PersistentVolumeClaim resource: %w", err)
		}

		applicationDeployment.Volumes[index] = pvc
	}

	if app.Port != 0 {
		service, err := createOkctlService(app)
		if err != nil {
			return nil, fmt.Errorf("unable to create service resource: %w", err)
		}

		applicationDeployment.Service = service
	}

	if app.Url != "" && app.Port != 0 {
		ingress, err := createOkctlIngress(app)
		if err != nil {
			return nil, err
		}

		applicationDeployment.Ingress = ingress
	}

	deployment, err := createOkctlDeployment(app)
	if err != nil {
		return nil, err
	}

	argoApp := createArgoApp(app, iacRepoURL)

	applicationDeployment.Deployment = deployment
	applicationDeployment.ArgoApplication = *argoApp

	return &applicationDeployment, nil
}

/*
ReadApplication returns an okctl Application based on stdin or a file
*/
func ReadApplication(path string) (api.Application, error) {
	var (
		rawApplication []byte
		err            error
	)

	if path == "-" {
		rawApplication, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return api.Application{}, fmt.Errorf("failed to read stdin: %w", err)
		}
	} else {
		rawApplication, err = ioutil.ReadFile(filepath.Clean(path))
		if err != nil {
			return api.Application{}, fmt.Errorf("failed to read file: %w", err)
		}
	}

	app, err := api.ParseApplication(string(rawApplication))
	if err != nil {
		return api.Application{}, fmt.Errorf("unable to parse application: %w", err)
	}

	return app, err
}
