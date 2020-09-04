// Package externalsecret implements a manifest builder and applyer
package externalsecret

import (
	"context"
	"fmt"

	v13 "github.com/oslokommune/okctl/pkg/kube/externalsecret/api/types/v1"

	v12 "github.com/oslokommune/okctl/pkg/kube/externalsecret/clientset/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ExternalSecret contains the state for building and applying the manifest
type ExternalSecret struct {
	Namespace string
	Name      string
	Data      map[string]string
	Ctx       context.Context
}

// New returns an initialised runner
func New(name, namespace string, data map[string]string) *ExternalSecret {
	return &ExternalSecret{
		Namespace: namespace,
		Name:      name,
		Data:      data,
		Ctx:       context.Background(),
	}
}

// CreateSecret invokes the client and creates the secret
func (a *ExternalSecret) CreateSecret(_ kubernetes.Interface, config *rest.Config) (interface{}, error) {
	clientSet, err := v12.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create external secrets client set: %w", err)
	}

	externalSecrets, err := clientSet.ExternalSecrets(a.Namespace).List(a.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, es := range externalSecrets.Items {
		if es.Name == a.Name {
			got, err := clientSet.ExternalSecrets(a.Namespace).Get(a.Ctx, es.Name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			return got, nil
		}
	}

	return clientSet.ExternalSecrets(a.Namespace).Create(a.Ctx, a.SecretManifest())
}

// SecretManifest returns the manifest
func (a *ExternalSecret) SecretManifest() *v13.ExternalSecret {
	e := &v13.ExternalSecret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ExternalSecret",
			APIVersion: "kubernetes-client.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      a.Name,
			Namespace: a.Namespace,
		},
		Spec: v13.ExternalSecretSpec{
			BackendType: "systemManager",
			Data:        []v13.ExternalSecretData{},
		},
	}

	for name, key := range a.Data {
		e.Spec.Data = append(e.Spec.Data, v13.ExternalSecretData{
			Key:  key,
			Name: name,
		})
	}

	return e
}
