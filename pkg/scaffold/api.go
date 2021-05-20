// Package scaffold knows how to generate necessary resources for deploying an okctl application
package scaffold

import (
	"encoding/json"
	"fmt"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/jsonpatch"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/scaffold/resources"
	v1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"sigs.k8s.io/yaml"
)

const defaultMainServicePortName = "main"

// ApplicationBase contains the content of the Kubernetes resource files
type ApplicationBase struct {
	Kustomization   []byte
	ArgoApplication []byte

	Deployment     []byte
	Ingress        []byte
	Service        []byte
	Volumes        []byte
	ServiceMonitor []byte
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

// GenerateApplicationBase converts an Application to Kustomize base files
//nolint:funlen,gocyclo
func GenerateApplicationBase(app v1alpha1.Application, iacRepoURL, relativeApplicationOverlayDir string) (ApplicationBase, error) {
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

		applicationBase.Volumes, err = volumesAsBytes(volumes)
		if err != nil {
			return applicationBase, err
		}

		if len(applicationBase.Volumes) > 0 {
			kustomization.AddResource("volumes.yaml")
		}
	}

	if app.HasService() {
		service := resources.CreateOkctlService(app, defaultMainServicePortName)

		kustomization.AddResource("service.yaml")

		applicationBase.Service, err = resourceAsBytes(service)
		if err != nil {
			return applicationBase, err
		}
	}

	if app.HasIngress() && app.HasService() {
		var ingress networkingv1beta1.Ingress

		ingress, err = resources.CreateOkctlIngress(app)
		if err != nil {
			return applicationBase, err
		}

		kustomization.AddResource("ingress.yaml")

		applicationBase.Ingress, err = resourceAsBytes(ingress)
		if err != nil {
			return applicationBase, err
		}
	}

	if app.HasPrometheus() {
		monitor := resources.CreateServiceMonitor(app, defaultMainServicePortName)

		kustomization.AddResource("service-monitor.yaml")

		applicationBase.ServiceMonitor, err = resourceAsBytes(monitor)
		if err != nil {
			return applicationBase, err
		}
	}

	deployment := resources.CreateOkctlDeployment(app)

	kustomization.AddResource("deployment.yaml")

	applicationBase.Deployment, err = resourceAsBytes(deployment)
	if err != nil {
		return applicationBase, err
	}

	argoApp := resources.CreateArgoApp(app, iacRepoURL, relativeApplicationOverlayDir)

	applicationBase.ArgoApplication, err = resourceAsBytes(argoApp)
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
func GenerateApplicationOverlay(application v1alpha1.Application, hostedZoneDomain, certARN string) (ApplicationOverlay, error) {
	var (
		err     error
		overlay = newApplicationOverlay()
	)

	kustomization := resources.NewKustomization()
	kustomization.AddResource("../../base")

	if application.HasIngress() {
		ingressPatch := jsonpatch.New()

		host := fmt.Sprintf("%s.%s", application.SubDomain, hostedZoneDomain)

		ingressPatch.Add(
			jsonpatch.Operation{
				Type:  jsonpatch.OperationTypeAdd,
				Path:  "/metadata/annotations/alb.ingress.kubernetes.io~1certificate-arn",
				Value: certARN,
			},
			jsonpatch.Operation{
				Type:  jsonpatch.OperationTypeAdd,
				Path:  "/spec/rules/0/host",
				Value: host,
			},
			jsonpatch.Operation{
				Type:  jsonpatch.OperationTypeAdd,
				Path:  "/spec/tls",
				Value: []string{},
			},
			jsonpatch.Operation{
				Type:  jsonpatch.OperationTypeAdd,
				Path:  "/spec/tls/0",
				Value: map[string]string{},
			},
			jsonpatch.Operation{
				Type:  jsonpatch.OperationTypeAdd,
				Path:  "/spec/tls/0/hosts",
				Value: []string{host},
			},
		)

		overlay.IngressPatch, err = json.Marshal(ingressPatch)
		if err != nil {
			return overlay, fmt.Errorf("marshalling ingress patch: %w", err)
		}

		kustomization.AddPatch(resources.PatchReference{
			Path:   constant.DefaultIngressPatchFilename,
			Target: resources.PatchTarget{Kind: "Ingress"},
		})
	}

	overlay.Kustomization, err = yaml.Marshal(kustomization)
	if err != nil {
		return overlay, fmt.Errorf("marshalling kustomization: %w", err)
	}

	return overlay, nil
}
