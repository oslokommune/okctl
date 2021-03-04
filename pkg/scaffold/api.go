package scaffold

import (
	"encoding/json"
	"fmt"
	kaex "github.com/oslokommune/kaex/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/scaffold/resources"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1beta1"
	"sigs.k8s.io/yaml"
)

type ApplicationBase struct {
	Kustomization   []byte
	ArgoApplication []byte

	Deployment []byte
	Ingress    []byte
	Service    []byte
	Volumes    []byte
}

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
func GenerateApplicationBase(app kaex.Application, iacRepoURL string, applicationOutputDir string) (ApplicationBase, error) {
	var err error
	applicationBase := NewApplicationBase()
	kustomization := NewKustomization()

	volumes := make([]*v1.PersistentVolumeClaim, len(app.Volumes))
	for index := range app.Volumes {

		pvc, err := resources.CreateOkctlVolume(app, app.Volumes[index])
		if err != nil {
			return applicationBase, fmt.Errorf("creating PersistentVolumeClaim resource: %w", err)
		}

		volumes[index] = &pvc
	}

	applicationBase.Volumes, err = resources.VolumesAsBytes(volumes)
	if err != nil {
		return applicationBase, err
	}

	if len(applicationBase.Volumes) > 0 {
		kustomization.AddResource("volumes.yaml")
	}

	var service v1.Service
	if app.Port != 0 {
		service, err = resources.CreateOkctlService(app)
		if err != nil {
			return applicationBase, fmt.Errorf("creating service resource: %w", err)
		}

		kustomization.AddResource("service.yaml")
	}

	applicationBase.Service, err = resources.ResourceAsBytes(service)
	if err != nil {
		return applicationBase, err
	}

	var ingress networkingv1.Ingress
	if app.Url != "" && app.Port != 0 {
		ingress, err = resources.CreateOkctlIngress(app)
		if err != nil {
			return applicationBase, err
		}

		kustomization.AddResource("ingress.yaml")
	}

	// TODO this needs to go into if above
	applicationBase.Ingress, err = resources.ResourceAsBytes(ingress)
	if err != nil {
		return applicationBase, err
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

	argoApp := resources.CreateArgoApp(app, iacRepoURL, applicationOutputDir)

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

func GenerateApplicationOverlay(application client.OkctlApplication, hostedZoneDomain, certARN string) (applicationOverlay, error) {
	var err error
	overlay := newApplicationOverlay()

	if application.HasIngress() {
		ingressPatch := NewPatch()

		host := fmt.Sprintf("%s.%s", application.SubDomain, hostedZoneDomain)

		ingressPatch.AddOperations(
			Operation{
				Type: OperationTypeAdd,
				Path: "/metadata/annotations",
				Value: map[string]string{
					"alb.ingress.kubernetes.io~1certificate-arn": certARN,
				},
			},
			Operation{
				Type:  OperationTypeAdd,
				Path:  "/spec/rules/0/host",
				Value: host,
			},
			Operation{
				Type:  OperationTypeAdd,
				Path:  "/spec/tls/0/hosts",
				Value: []string{host},
			},
		)

		overlay.IngressPatch, err = json.Marshal(ingressPatch)
		if err != nil {
			return overlay, fmt.Errorf("marshalling ingress patch: %w", err)
		}
	}

	return overlay, nil
}
