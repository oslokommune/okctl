// Package scaffold knows how to generate necessary resources for deploying an okctl application
package scaffold

import (
	"bytes"
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/jsonpatch"

	"github.com/oslokommune/okctl/pkg/scaffold/resources"
	v1 "k8s.io/api/core/v1"
)

const (
	defaultMainServicePortName = "main"
	deploymentImagePath        = "/spec/template/spec/containers/0/image"
)

// GenerateApplicationBase converts an Application to Kustomize base files
//nolint:funlen,gocyclo,gocognit
func GenerateApplicationBase(opts GenerateApplicationBaseOpts) error {
	if len(opts.Application.Volumes) != 0 {
		volumes := make([]*v1.PersistentVolumeClaim, len(opts.Application.Volumes))

		for index := range opts.Application.Volumes {
			pvc, err := resources.CreateOkctlVolume(opts.Application, opts.Application.Volumes[index])
			if err != nil {
				return fmt.Errorf("creating PersistentVolumeClaim resource: %w", err)
			}

			volumes[index] = &pvc
		}

		raw, err := volumesAsBytes(volumes)
		if err != nil {
			return fmt.Errorf("converting pvc's to bytes: %w", err)
		}

		err = opts.SaveManifest("volumes.yaml", raw)
		if err != nil {
			return fmt.Errorf("saving volume manifest: %w", err)
		}
	}

	if opts.Application.HasService() {
		service := resources.CreateOkctlService(opts.Application, defaultMainServicePortName)

		raw, err := ResourceAsBytes(service)
		if err != nil {
			return err
		}

		err = opts.SaveManifest("service.yaml", raw)
		if err != nil {
			return fmt.Errorf("saving service manifest: %w", err)
		}
	}

	if opts.Application.HasIngress() && opts.Application.HasService() {
		ingress, err := resources.CreateOkctlIngress(opts.Application)
		if err != nil {
			return err
		}

		raw, err := ResourceAsBytes(ingress)
		if err != nil {
			return err
		}

		err = opts.SaveManifest("ingress.yaml", raw)
		if err != nil {
			return fmt.Errorf("saving ingress manifest: %w", err)
		}
	}

	if opts.Application.HasPrometheus() {
		monitor := resources.CreateServiceMonitor(opts.Application, defaultMainServicePortName)

		raw, err := ResourceAsBytes(monitor)
		if err != nil {
			return err
		}

		err = opts.SaveManifest("service-monitor.yaml", raw)
		if err != nil {
			return fmt.Errorf("saving service monitor manifest: %w", err)
		}
	}

	rawNamespace, err := ResourceAsBytes(resources.CreateNamespace(opts.Application))
	if err != nil {
		return err
	}

	err = opts.SaveManifest("namespace.yaml", rawNamespace)
	if err != nil {
		return fmt.Errorf("saving namespace manifest: %w", err)
	}

	rawDeployment, err := ResourceAsBytes(resources.CreateOkctlDeployment(opts.Application))
	if err != nil {
		return err
	}

	err = opts.SaveManifest("deployment.yaml", rawDeployment)
	if err != nil {
		return fmt.Errorf("saving deployment manifest: %w", err)
	}

	return nil
}

// GenerateApplicationOverlay generates patches for environment specific parts of the kubernetes resources
func GenerateApplicationOverlay(opts GenerateApplicationOverlayOpts) error {
	patch := *jsonpatch.New()

	patch.Add(jsonpatch.Operation{
		Type:  jsonpatch.OperationTypeAdd,
		Path:  deploymentImagePath,
		Value: opts.Application.Image.URI,
	})

	err := opts.SavePatch("Deployment", patch)
	if err != nil {
		return fmt.Errorf("saving deployment patch: %w", err)
	}

	if opts.Application.HasIngress() {
		patch = createIngressPatch(opts.Application.SubDomain, opts.Domain, opts.CertificateARN)

		err = opts.SavePatch("Ingress", patch)
		if err != nil {
			return fmt.Errorf("saving ingress patch: %w", err)
		}
	}

	return nil
}

// GenerateArgoCDApplicationManifest generates an ArgoCD Application manifest
func GenerateArgoCDApplicationManifest(opts GenerateArgoCDApplicationManifestOpts) (io.Reader, error) {
	argoApp := resources.CreateArgoApp(opts.Name, opts.Namespace, opts.IACRepoURL, opts.SourceSyncDir)

	raw, err := ResourceAsBytes(argoApp)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(raw), nil
}

func createIngressPatch(subDomain, domain, certArn string) jsonpatch.Patch {
	patch := *jsonpatch.New()

	host := fmt.Sprintf("%s.%s", subDomain, domain)

	patch.Add(
		jsonpatch.Operation{
			Type:  jsonpatch.OperationTypeAdd,
			Path:  "/metadata/annotations/alb.ingress.kubernetes.io~1certificate-arn",
			Value: certArn,
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

	return patch
}
