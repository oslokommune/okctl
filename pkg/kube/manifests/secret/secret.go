// Package secret provides a secret creator and applier
package secret

import (
	"context"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Secret contains the state for creating a secret
type Secret struct {
	Name      string
	Namespace string
	Secret    *v1.Secret
	Ctx       context.Context
}

// New returns an initialised secret creator
func New(name, namespace string, secret *v1.Secret) *Secret {
	return &Secret{
		Name:      name,
		Namespace: namespace,
		Secret:    secret,
		Ctx:       context.Background(),
	}
}

// DeleteSecret deletes the secret
func (s *Secret) DeleteSecret(_ kubernetes.Interface, config *rest.Config) (interface{}, error) {
	client := kubernetes.NewForConfigOrDie(config)

	ns, err := client.CoreV1().Secrets(s.Namespace).List(s.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	found := false

	for _, item := range ns.Items {
		if item.Name == s.Name {
			found = true
		}
	}

	if !found {
		return nil, nil
	}

	policy := metav1.DeletePropagationForeground

	return nil, client.CoreV1().Secrets(s.Namespace).Delete(s.Ctx, s.Name, metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
}

// CreateSecret creates the secret
func (s *Secret) CreateSecret(clientset kubernetes.Interface, _ *rest.Config) (interface{}, error) {
	client := clientset.CoreV1().Secrets(s.Namespace)

	secrets, err := client.List(s.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ns := range secrets.Items {
		if ns.Name == s.Name {
			r, err := client.Get(s.Ctx, ns.Name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			return r, nil
		}
	}

	return client.Create(s.Ctx, s.Secret, metav1.CreateOptions{})
}

// NewManifest returns the secret manifest
func NewManifest(name, namespace string, data, labels map[string]string) *v1.Secret {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		StringData: data,
		Type:       v1.SecretTypeOpaque,
	}
}
