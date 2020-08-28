package run

import (
	"fmt"
	"time"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/kube"
	"github.com/oslokommune/okctl/pkg/kube/manifests/externaldns"
	"sigs.k8s.io/yaml"
)

type kubeRun struct {
	kubeConfStore api.KubeConfigStore
}

func (k *kubeRun) CreateExternalDNSKubeDeployment(opts api.CreateExternalDNSKubeDeploymentOpts) (*api.Kube, error) {
	kubeConfig, err := k.kubeConfStore.GetKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve kubeconfig: %w", err)
	}

	ext := externaldns.New(opts.HostedZoneID, opts.DomainFilter)

	client, err := kube.New(kubeConfig.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	resources, err := client.Apply(ext.CreateDeployment, ext.CreateClusterRole, ext.CreateClusterRoleBinding)
	if err != nil {
		return nil, fmt.Errorf("failed to apply kubernets manifests: %w", err)
	}

	err = client.Watch(resources, 2*time.Minute) // nolint: gomnd
	if err != nil {
		return nil, fmt.Errorf("failed while waiting for resources to be created: %w", err)
	}

	deployment, err := yaml.Marshal(ext.DeploymentManifest())
	if err != nil {
		return nil, fmt.Errorf("failed to serialise Deployment manifest: %w", err)
	}

	clusterRole, err := yaml.Marshal(ext.ClusterRoleManifest())
	if err != nil {
		return nil, fmt.Errorf("failed to serialise ClusterRole manifest: %w", err)
	}

	clusterRoleBinding, err := yaml.Marshal(ext.ClusterRoleBindingManifest())
	if err != nil {
		return nil, fmt.Errorf("failed to serialise ClusterRoleBinding manifest: %w", err)
	}

	return &api.Kube{
		HostedZoneID: opts.HostedZoneID,
		DomainFilter: opts.DomainFilter,
		Manifests: map[string][]byte{
			"deployment.yaml":         deployment,
			"clusterrole.yaml":        clusterRole,
			"clusterrolebinding.yaml": clusterRoleBinding,
		},
	}, nil
}

// NewKubeRun returns an initialised kube runner
func NewKubeRun(kubeConfStore api.KubeConfigStore) api.KubeRun {
	return &kubeRun{
		kubeConfStore: kubeConfStore,
	}
}
