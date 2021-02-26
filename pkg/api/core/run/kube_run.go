package run

import (
	"fmt"
	"time"

	"github.com/oslokommune/okctl/pkg/kube/manifests/secret"

	"github.com/oslokommune/okctl/pkg/kube/manifests/storageclass"

	"github.com/oslokommune/okctl/pkg/credentials/aws"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/kube/manifests/namespace"

	"github.com/oslokommune/okctl/pkg/kube/manifests/externalsecret"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/kube"
	"github.com/oslokommune/okctl/pkg/kube/manifests/externaldns"
	"sigs.k8s.io/yaml"
)

type kubeRun struct {
	provider v1alpha1.CloudProvider
	auth     aws.Authenticator
}

func (k *kubeRun) CreateNativeSecret(opts api.CreateNativeSecretOpts) (*api.NativeSecretKube, error) {
	sec := secret.New(opts.Name, opts.Namespace, secret.NewManifest(
		opts.Name,
		opts.Namespace,
		opts.Data,
		opts.Labels,
	))

	client, err := kube.New(kube.NewFromEKSCluster(opts.ID.ClusterName, opts.ID.Region, k.provider, k.auth))
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	_, err = client.Apply(kube.Applier{
		Fn:          sec.CreateSecret,
		Description: fmt.Sprintf("creating secret: %s, at: %s", opts.Name, opts.Namespace),
	})
	if err != nil {
		return nil, fmt.Errorf("creating secret: %w", err)
	}

	data, err := yaml.Marshal(sec.Manifest)
	if err != nil {
		return nil, fmt.Errorf("marshalling manifest: %w", err)
	}

	return &api.NativeSecretKube{
		ID:        opts.ID,
		Name:      opts.Name,
		Namespace: opts.Namespace,
		Manifest:  data,
	}, nil
}

func (k *kubeRun) DeleteNativeSecret(opts api.DeleteNativeSecretOpts) error {
	client, err := kube.New(kube.NewFromEKSCluster(opts.ID.ClusterName, opts.ID.Region, k.provider, k.auth))
	if err != nil {
		return fmt.Errorf("creating kubernetes client: %w", err)
	}

	_, err = client.Apply(kube.Applier{
		Fn:          secret.New(opts.Name, opts.Namespace, nil).DeleteSecret,
		Description: fmt.Sprintf("deleting secret: %s, from: %s", opts.Name, opts.Namespace),
	})
	if err != nil {
		return fmt.Errorf("deleting secret: %w", err)
	}

	return nil
}

func (k *kubeRun) CreateStorageClass(opts api.CreateStorageClassOpts) (*api.StorageClassKube, error) {
	sc, err := storageclass.New(opts.Name, opts.Parameters, opts.Annotations)
	if err != nil {
		return nil, fmt.Errorf("creating manifest: %w", err)
	}

	client, err := kube.New(kube.NewFromEKSCluster(opts.ID.ClusterName, opts.ID.Region, k.provider, k.auth))
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	_, err = client.Apply(kube.Applier{
		Fn:          sc.CreateStorageClass,
		Description: fmt.Sprintf("storageclass: %s", opts.Name),
	})
	if err != nil {
		return nil, fmt.Errorf("creating storageclass: %w", err)
	}

	manifest, err := sc.StorageClassManifest().Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshalling manifest: %w", err)
	}

	return &api.StorageClassKube{
		ID:       opts.ID,
		Name:     opts.Name,
		Manifest: manifest,
	}, nil
}

func (k *kubeRun) DeleteNamespace(opts api.DeleteNamespaceOpts) error {
	ns := namespace.New(opts.Namespace)

	client, err := kube.New(kube.NewFromEKSCluster(opts.ID.ClusterName, opts.ID.Region, k.provider, k.auth))
	if err != nil {
		return fmt.Errorf("creating kubernetes client: %w", err)
	}

	_, err = client.Apply(kube.Applier{
		Fn:          ns.DeleteNamespace,
		Description: fmt.Sprintf("deleting namespace: %s", opts.Namespace),
	})
	if err != nil {
		return fmt.Errorf("deleting namespace: %w", err)
	}

	return nil
}

func (k *kubeRun) DeleteExternalSecrets(opts api.DeleteExternalSecretsOpts) error {
	fns := make([]kube.Applier, len(opts.Manifests))

	i := 0

	for name, ns := range opts.Manifests {
		fns[i] = kube.Applier{
			Fn:          externalsecret.New(name, ns, nil).DeleteSecret,
			Description: fmt.Sprintf("delete secret: %s, from: %s", name, ns),
		}

		i++
	}

	client, err := kube.New(kube.NewFromEKSCluster(opts.ID.ClusterName, opts.ID.Region, k.provider, k.auth))
	if err != nil {
		return fmt.Errorf("creating kubernetes client: %w", err)
	}

	_, err = client.Apply(fns...)
	if err != nil {
		return fmt.Errorf("applying kubernetes manifests: %w", err)
	}

	return nil
}

// In all fairness, we should refactor this, probably by extending the functionality
// on the kube side. First we collect all apply things, then we apply, or something like
// this.
// nolint: funlen
func (k *kubeRun) CreateExternalSecrets(opts api.CreateExternalSecretsOpts) (*api.ExternalSecretsKube, error) {
	fns := make([]kube.Applier, len(opts.Manifests))
	manifests := map[string][]byte{}
	namespaces := map[string]struct{}{}

	for i, manifest := range opts.Manifests {
		data := map[string]string{}

		for _, d := range manifest.Data {
			data[d.Name] = d.Key
		}

		secretManifest := externalsecret.SecretManifest(manifest.Name, manifest.Namespace, manifest.Annotations, manifest.Labels, data)
		fn := externalsecret.New(manifest.Name, manifest.Namespace, secretManifest)

		raw, err := yaml.Marshal(secretManifest)
		if err != nil {
			return nil, fmt.Errorf("marshalling manifest: %w", err)
		}

		fns[i] = kube.Applier{
			Fn:          fn.CreateSecret,
			Description: fmt.Sprintf("external secret %s in %s", manifest.Name, manifest.Namespace),
		}

		manifests[fmt.Sprintf("external-secret-%s.yml", manifest.Name)] = raw
		namespaces[manifest.Namespace] = struct{}{}
	}

	for ns := range namespaces {
		newNS := namespace.New(ns)

		fns = append([]kube.Applier{
			{
				Fn:          newNS.CreateNamespace,
				Description: fmt.Sprintf("namespace %s", ns),
			},
		}, fns...)

		data, err := yaml.Marshal(newNS.NamespaceManifest())
		if err != nil {
			return nil, fmt.Errorf("marshalling manifest: %w", err)
		}

		manifests[fmt.Sprintf("namespace-%s.yml", ns)] = data
	}

	client, err := kube.New(kube.NewFromEKSCluster(opts.ID.ClusterName, opts.ID.Region, k.provider, k.auth))
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	// Should probably watch these..
	_, err = client.Apply(fns...)
	if err != nil {
		return nil, fmt.Errorf("applying kubernetes manifests: %w", err)
	}

	return &api.ExternalSecretsKube{
		ID:        opts.ID,
		Manifests: manifests,
	}, nil
}

func (k *kubeRun) CreateExternalDNSKubeDeployment(opts api.CreateExternalDNSKubeDeploymentOpts) (*api.ExternalDNSKube, error) {
	ext := externaldns.New(opts.HostedZoneID, opts.DomainFilter)

	client, err := kube.New(kube.NewFromEKSCluster(opts.ID.ClusterName, opts.ID.Region, k.provider, k.auth))
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	resources, err := client.Apply(
		kube.Applier{Fn: ext.CreateDeployment, Description: "external dns deployment"},
		kube.Applier{Fn: ext.CreateClusterRole, Description: "external dns cluster role"},
		kube.Applier{Fn: ext.CreateClusterRoleBinding, Description: "external dns cluster role binding"},
	)
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

	return &api.ExternalDNSKube{
		ID:           opts.ID,
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
func NewKubeRun(provider v1alpha1.CloudProvider, auth aws.Authenticator) api.KubeRun {
	return &kubeRun{
		auth:     auth,
		provider: provider,
	}
}
