package run

import (
	"fmt"
	"time"

	v1 "github.com/oslokommune/okctl/pkg/kube/externalsecret/api/types/v1"

	"github.com/oslokommune/okctl/pkg/kube/manifests/scale"

	"github.com/oslokommune/okctl/pkg/kube/manifests/configmap"

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

func (k *kubeRun) ScaleDeployment(opts api.ScaleDeploymentOpts) error {
	s := scale.New(opts.Name, opts.Namespace, opts.Replicas)

	client, err := kube.New(kube.NewFromEKSCluster(opts.ID.ClusterName, opts.ID.Region, k.provider, k.auth))
	if err != nil {
		return fmt.Errorf("creating kubernetes client: %w", err)
	}

	_, err = client.Apply(kube.Applier{
		Fn:          s.Scale,
		Description: fmt.Sprintf("scaling deployment: %s, at: %s, to: %d", opts.Name, opts.Namespace, opts.Replicas),
	})
	if err != nil {
		return fmt.Errorf("scaling deployment: %w", err)
	}

	return nil
}

func (k *kubeRun) CreateConfigMap(opts api.CreateConfigMapOpts) (*api.ConfigMap, error) {
	sec := configmap.New(opts.Name, opts.Namespace, configmap.NewManifest(
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
		Fn:          sec.CreateConfigMap,
		Description: fmt.Sprintf("creating configmap: %s, at: %s", opts.Name, opts.Namespace),
	})
	if err != nil {
		return nil, fmt.Errorf("creating configmap: %w", err)
	}

	data, err := yaml.Marshal(sec.Manifest)
	if err != nil {
		return nil, fmt.Errorf("marshalling manifest: %w", err)
	}

	return &api.ConfigMap{
		ID:        opts.ID,
		Name:      opts.Name,
		Namespace: opts.Namespace,
		Manifest:  data,
	}, nil
}

func (k *kubeRun) DeleteConfigMap(opts api.DeleteConfigMapOpts) error {
	client, err := kube.New(kube.NewFromEKSCluster(opts.ID.ClusterName, opts.ID.Region, k.provider, k.auth))
	if err != nil {
		return fmt.Errorf("creating kubernetes client: %w", err)
	}

	_, err = client.Apply(kube.Applier{
		Fn:          configmap.New(opts.Name, opts.Namespace, nil).DeleteConfigMap,
		Description: fmt.Sprintf("deleting configmap: %s, from: %s", opts.Name, opts.Namespace),
	})
	if err != nil {
		return fmt.Errorf("deleting configmap: %w", err)
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

func (k *kubeRun) CreateNamespace(opts api.CreateNamespaceOpts) (*api.Namespace, error) {
	ns := namespace.New(opts.Namespace, opts.Labels)

	client, err := kube.New(kube.NewFromEKSCluster(opts.ID.ClusterName, opts.ID.Region, k.provider, k.auth))
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	_, err = client.Apply(kube.Applier{
		Fn:          ns.CreateNamespace,
		Description: fmt.Sprintf("creating namespace: %s", opts.Namespace),
	})
	if err != nil {
		return nil, fmt.Errorf("creating namespace: %w", err)
	}

	manifest, err := ns.NamespaceManifest().Marshal()
	if err != nil {
		return nil, fmt.Errorf("marshalling manifest: %w", err)
	}

	return &api.Namespace{
		ID:        opts.ID,
		Namespace: opts.Namespace,
		Labels:    opts.Labels,
		Manifest:  manifest,
	}, nil
}

func (k *kubeRun) DeleteNamespace(opts api.DeleteNamespaceOpts) error {
	ns := namespace.New(opts.Namespace, nil)

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

func (k *kubeRun) CreateExternalSecrets(opts api.CreateExternalSecretsOpts) (*api.ExternalSecretsKube, error) {
	data := make([]v1.ExternalSecretData, len(opts.Manifest.Data))

	for i, d := range opts.Manifest.Data {
		data[i] = v1.ExternalSecretData{
			Key:      d.Key,
			Name:     d.Name,
			Property: d.Property,
		}
	}

	secretManifest := externalsecret.SecretManifest(
		opts.Manifest.Name,
		opts.Manifest.Namespace,
		opts.Manifest.Backend,
		opts.Manifest.Annotations,
		opts.Manifest.Labels,
		data,
	)

	fn := externalsecret.New(opts.Manifest.Name, opts.Manifest.Namespace, secretManifest)

	raw, err := yaml.Marshal(secretManifest)
	if err != nil {
		return nil, fmt.Errorf("marshalling manifest: %w", err)
	}

	client, err := kube.New(kube.NewFromEKSCluster(opts.ID.ClusterName, opts.ID.Region, k.provider, k.auth))
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}

	_, err = client.Apply(kube.Applier{
		Fn:          fn.CreateSecret,
		Description: fmt.Sprintf("external secret %s in %s", opts.Manifest.Name, opts.Manifest.Namespace),
	})
	if err != nil {
		return nil, fmt.Errorf("applying kubernetes manifests: %w", err)
	}

	return &api.ExternalSecretsKube{
		ID:        opts.ID,
		Name:      opts.Manifest.Name,
		Namespace: opts.Manifest.Namespace,
		Content:   raw,
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

	err = client.Watch(resources, 4*time.Minute) // nolint: gomnd
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
