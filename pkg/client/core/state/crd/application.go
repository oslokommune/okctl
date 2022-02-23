package crd

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/clients/kubectl"
	"github.com/oslokommune/okctl/pkg/logging"

	"sigs.k8s.io/yaml"
)

// Put adds an application manifest to etcd
func (a applicationState) Put(application v1alpha1.Application) error {
	rawApplication, err := yaml.Marshal(application)
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}

	err = a.kubectl.Apply(bytes.NewReader(rawApplication))
	if err != nil {
		return fmt.Errorf("applying: %w", err)
	}

	return nil
}

// Get retrieves an application manifest from etcd
func (a applicationState) Get(name string) (v1alpha1.Application, error) {
	resource, err := a.kubectl.Get(kubectl.Resource{
		Name:      name,
		Namespace: "okctl",
	})
	if err != nil {
		return v1alpha1.Application{}, fmt.Errorf("retrieving resource: %w", err)
	}

	rawResource, err := io.ReadAll(resource)
	if err != nil {
		return v1alpha1.Application{}, fmt.Errorf("buffering: %w", err)
	}

	var appManifest v1alpha1.Application

	err = yaml.Unmarshal(rawResource, &appManifest)
	if err != nil {
		return v1alpha1.Application{}, fmt.Errorf("unmarshalling app manifest: %w", err)
	}

	return appManifest, nil
}

// Delete removes an application manifest from etcd
func (a applicationState) Delete(name string) error {
	err := a.kubectl.DeleteByResource(kubectl.Resource{
		Name:      name,
		Kind:      "application.okctl.io",
		Namespace: "okctl",
	})
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}

	return nil
}

// List returns all application manifests in etcd
func (a applicationState) List() ([]v1alpha1.Application, error) {
	resources, err := a.kubectl.Get(kubectl.Resource{
		Namespace: "okctl",
		Kind:      "application.okctl.io",
	})
	if err != nil {
		return nil, fmt.Errorf("retrieving: %w", err)
	}

	rawResources, err := io.ReadAll(resources)
	if err != nil {
		return nil, fmt.Errorf("buffering: %w", err)
	}

	var manifests []v1alpha1.Application

	err = yaml.Unmarshal(rawResources, &manifests)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	return manifests, nil
}

// Initialize knows how to apply the CustomResourceDefinition, so we can store application manifests in etcd
func (a applicationState) Initialize() error {
	log := logging.GetLogger("applicationState", "Initialize")

	_, err := a.kubectl.Get(kubectl.Resource{
		Name:      "applications.okctl.io",
		Namespace: "okctl",
		Kind:      "CustomResourceDefinition",
	})
	if err == nil {
		log.Debug("Found existing custom resource definition")

		return nil
	}

	if err != nil && !errors.Is(err, kubectl.ErrNotFound) {
		return fmt.Errorf("retrieving application CRD: %w", err)
	}

	log.Debug("No existing custom resource definition found, creating.")

	err = a.kubectl.Apply(bytes.NewReader(applicationCRDTemplate))
	if err != nil {
		return fmt.Errorf("applying application custom resource definition: %w", err)
	}

	return nil
}

// NewApplicationState returns an initialized application state instance
func NewApplicationState(kubectlClient kubectl.Client) client.ApplicationState {
	return &applicationState{
		kubectl: kubectlClient,
	}
}

type applicationState struct {
	kubectl kubectl.Client
}

//go:embed application-crd-template.yaml
var applicationCRDTemplate []byte
