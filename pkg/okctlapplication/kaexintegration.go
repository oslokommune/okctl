package okctlapplication

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1beta1"

	"github.com/oslokommune/kaex/pkg/api"
)

/*
AcquireApplication returns an okctl Application based on stdin or a file
*/
func AcquireApplication(path string) (api.Application, error) {
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

	if app.Ingress.Annotations == nil {
		app.Ingress.Annotations = map[string]string{}
	}

	app.Ingress.Annotations["kubernetes.io/ingress.class"] = "alb"
	app.Ingress.Annotations["alb.ingress.kubernetes.io/scheme"] = "internet-facing"

	return app, err
}

func createOkctlVolume(app api.Application, volume map[string]string) (v1.PersistentVolumeClaim, error) {
	var (
		path string
		size string
	)

	for path, size = range volume {
		break
	}

	pvc, err := api.CreatePersistentVolume(app, path, size)
	if err != nil {
		return v1.PersistentVolumeClaim{}, fmt.Errorf("error creating pvc: %w", err)
	}

	return pvc, nil
}

func createOkctlService(app api.Application) (v1.Service, error) {
	service, err := api.CreateService(app)
	if err != nil {
		return v1.Service{}, fmt.Errorf("error creating kaex service: %w", err)
	}

	service.Spec.Type = "NodePort"

	return service, nil
}

func createOkctlIngress(app api.Application) (networkingv1.Ingress, error) {
	ingress, err := api.CreateIngress(app)
	if err != nil {
		return networkingv1.Ingress{}, err
	}

	return ingress, nil
}

func createOkctlDeployment(app api.Application) (appsv1.Deployment, error) {
	deployment, err := api.CreateDeployment(app)
	if err != nil {
		return appsv1.Deployment{}, err
	}

	return deployment, nil
}

/*
ArgoCDDeploymentResources contains the necessary resources for a ArgoCD deployment
*/
type ArgoCDDeploymentResources struct {
	KubernetesResourcesBuffer bytes.Buffer
	ArgoAppBuffer             bytes.Buffer
}

/*
ConvertApplicationToResources turns an api.Application into a ArgoCDDeploymentResources containing buffers with relevant
*/
func ConvertApplicationToResources(app api.Application, iacRepoURL string) (*ArgoCDDeploymentResources, error) {
	expandedApp := ArgoCDDeploymentResources{}

	for path := range app.Volumes {
		pvc, err := createOkctlVolume(app, app.Volumes[path])
		if err != nil {
			return nil, err
		}

		err = api.WriteResource(&expandedApp.KubernetesResourcesBuffer, pvc)
		if err != nil {
			return nil, err
		}
	}

	if app.Port != 0 {
		service, err := createOkctlService(app)
		if err != nil {
			return nil, fmt.Errorf("unable to create service resource: %w", err)
		}

		err = api.WriteResource(&expandedApp.KubernetesResourcesBuffer, service)
		if err != nil {
			return nil, fmt.Errorf("error writing service to buffer: %w", err)
		}
	}

	if app.Url != "" && app.Port != 0 {
		ingress, err := createOkctlIngress(app)
		if err != nil {
			return nil, err
		}

		err = api.WriteResource(&expandedApp.KubernetesResourcesBuffer, ingress)
		if err != nil {
			return nil, err
		}
	}

	deployment, err := createOkctlDeployment(app)
	if err != nil {
		return nil, err
	}

	err = api.WriteResource(&expandedApp.KubernetesResourcesBuffer, deployment)
	if err != nil {
		return nil, err
	}

	argoApp, err := CreateArgoApp(app, iacRepoURL)
	if err != nil {
		return &expandedApp, fmt.Errorf("error creating ArgoApp from application.yaml: %w", err)
	}

	err = api.WriteResource(&expandedApp.ArgoAppBuffer, argoApp)
	if err != nil {
		return &expandedApp, fmt.Errorf("error writing ArgoApp to buffer: %w", err)
	}

	return &expandedApp, nil
}
