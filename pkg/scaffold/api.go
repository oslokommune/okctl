// Package scaffold knows how to generate necessary resources for deploying an okctl application
package scaffold

import (
	"encoding/json"
	"fmt"

	kaex "github.com/oslokommune/kaex/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/scaffold/resources"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1beta1"
	"sigs.k8s.io/yaml"
)

// ApplicationBase contains the content of the Kubernetes resource files
type ApplicationBase struct {
	Kustomization   []byte
	ArgoApplication []byte

	Deployment []byte
	Ingress    []byte
	Service    []byte
	Volumes    []byte
}

// NewApplicationBase returns an initialized ApplicationBase struct
func NewApplicationBase() ApplicationBase {
	return ApplicationBase{
		Kustomization:   []byte(""),
		ArgoApplication: []byte(""),
		Deployment:      []byte(""),
		Ingress:         []byte(""),
		Service:         []byte(""),
		Volumes:         []byte(""),
	}
}

// GenerateApplicationBase converts a Kaex Application to Kustomize base files
// nolint: funlen gocyclo
func GenerateApplicationBase(app kaex.Application, iacRepoURL, relativeApplicationOverlayDir string) (ApplicationBase, error) {
	var (
		err             error
		applicationBase = NewApplicationBase()
		kustomization   = resources.NewKustomization()
	)

	volumes := make([]*v1.PersistentVolumeClaim, len(app.Volumes))

	for index := range app.Volumes {
		pvc, err := resources.CreateOkctlVolume(app, app.Volumes[index])
		if err != nil {
			return applicationBase, fmt.Errorf("creating PersistentVolumeClaim resource: %w", err)
		}

		volumes[index] = &pvc

		applicationBase.Volumes, err = resources.VolumesAsBytes(volumes)
		if err != nil {
			return applicationBase, err
		}

		if len(applicationBase.Volumes) > 0 {
			kustomization.AddResource("volumes.yaml")
		}
	}

	if app.Port != 0 {
		var service v1.Service

		service, err = resources.CreateOkctlService(app)
		if err != nil {
			return applicationBase, fmt.Errorf("creating service resource: %w", err)
		}

		kustomization.AddResource("service.yaml")

		applicationBase.Service, err = resources.ResourceAsBytes(service)
		if err != nil {
			return applicationBase, err
		}
	}

	if app.Url != "" && app.Port != 0 {
		var ingress networkingv1.Ingress

		ingress, err = resources.CreateOkctlIngress(app)
		if err != nil {
			return applicationBase, err
		}

		kustomization.AddResource("ingress.yaml")

		applicationBase.Ingress, err = resources.ResourceAsBytes(ingress)
		if err != nil {
			return applicationBase, err
		}
	}

	deployment, err := resources.CreateOkctlDeployment(app)
	if err != nil {
		return applicationBase, err
	}

	kustomization.AddResource("deployment.yaml")

	applicationBase.Deployment, err = resources.ResourceAsBytes(deployment)
	if err != nil {
		return applicationBase, err
	}

	argoApp := resources.CreateArgoApp(app, iacRepoURL, relativeApplicationOverlayDir)

	applicationBase.ArgoApplication, err = resources.ResourceAsBytes(argoApp)
	if err != nil {
		return applicationBase, err
	}

	applicationBase.Kustomization, err = yaml.Marshal(kustomization)
	if err != nil {
		return applicationBase, err
	}

	return applicationBase, nil
}

// GenerateApplicationOverlay generates patches for environment specific parts of the kubernetes resources
func GenerateApplicationOverlay(application client.OkctlApplication, hostedZoneDomain, certARN string) (ApplicationOverlay, error) {
	var (
		err     error
		overlay = newApplicationOverlay()
	)

	kustomization := resources.NewKustomization()
	kustomization.AddResource("../../base")

	if application.HasIngress() {
		ingressPatch := resources.NewPatch()

		host := fmt.Sprintf("%s.%s", application.SubDomain, hostedZoneDomain)

		ingressPatch.AddOperations(
			resources.Operation{
				Type:  resources.OperationTypeAdd,
				Path:  "/metadata/annotations/alb.ingress.kubernetes.io~1certificate-arn",
				Value: certARN,
			},
			resources.Operation{
				Type:  resources.OperationTypeAdd,
				Path:  "/spec/rules/0/host",
				Value: host,
			},
			resources.Operation{
				Type:  resources.OperationTypeAdd,
				Path:  "/spec/tls",
				Value: []string{},
			},
			resources.Operation{
				Type:  resources.OperationTypeAdd,
				Path:  "/spec/tls/0",
				Value: map[string]string{},
			},
			resources.Operation{
				Type:  resources.OperationTypeAdd,
				Path:  "/spec/tls/0/hosts",
				Value: []string{host},
			},
		)

		overlay.IngressPatch, err = json.Marshal(ingressPatch)
		if err != nil {
			return overlay, fmt.Errorf("marshalling ingress patch: %w", err)
		}

		kustomization.AddPatch(resources.PatchReference{
			Path:   config.DefaultIngressPatchFilename,
			Target: resources.PatchTarget{Kind: "Ingress"},
		})
	}

	overlay.Kustomization, err = yaml.Marshal(kustomization)
	if err != nil {
		return overlay, fmt.Errorf("marshalling kustomization: %w", err)
	}

	return overlay, nil
}
