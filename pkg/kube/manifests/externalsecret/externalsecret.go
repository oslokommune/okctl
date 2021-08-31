// Package externalsecret implements a manifest builder and applyer
package externalsecret

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

	typesv1 "github.com/oslokommune/okctl/pkg/kube/externalsecret/api/types/v1"

	clientv1 "github.com/oslokommune/okctl/pkg/kube/externalsecret/clientset/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ExternalSecret contains the state for building and applying the manifest
type ExternalSecret struct {
	Namespace string
	Name      string
	Manifest  *typesv1.ExternalSecret
	Ctx       context.Context
}

// New returns an initialised runner
func New(name, namespace string, manifest *typesv1.ExternalSecret) *ExternalSecret {
	return &ExternalSecret{
		Namespace: namespace,
		Name:      name,
		Manifest:  manifest,
		Ctx:       context.Background(),
	}
}

// CreateSecret invokes the client and creates the secret
func (a *ExternalSecret) CreateSecret(_ kubernetes.Interface, config *rest.Config) (interface{}, error) {
	clientSet, err := clientv1.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf(constant.CreateExternalSecretSetError, err)
	}

	externalSecrets, err := clientSet.ExternalSecrets(a.Namespace).List(a.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf(constant.ListExternalSecretsError, a.Namespace, err)
	}

	for _, es := range externalSecrets.Items {
		if es.Name == a.Name {
			got, err := clientSet.ExternalSecrets(a.Namespace).Get(a.Ctx, es.Name, metav1.GetOptions{})
			if err != nil {
				return nil, fmt.Errorf(constant.GetExternalSecretsError, es.Name, es.Namespace, err)
			}

			return got, nil
		}
	}

	return clientSet.ExternalSecrets(a.Namespace).Create(a.Ctx, a.Manifest)
}

// DeleteSecret invokes the client and deletes the secret
func (a *ExternalSecret) DeleteSecret(_ kubernetes.Interface, config *rest.Config) (interface{}, error) {
	clientSet, err := clientv1.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf(constant.CreateExternalSecretSetError, err)
	}

	externalSecrets, err := clientSet.ExternalSecrets(a.Namespace).List(a.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf(constant.ListExternalSecretsError, a.Namespace, err)
	}

	deletePolicy := metav1.DeletePropagationForeground

	for _, es := range externalSecrets.Items {
		if es.Name == a.Name {
			return nil, clientSet.ExternalSecrets(a.Namespace).Delete(a.Ctx, a.Name, metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			})
		}
	}

	return nil, nil
}

// SecretManifest returns the manifest
func SecretManifest(name, namespace, backendType string, annotations, labels map[string]string, data []typesv1.ExternalSecretData) *typesv1.ExternalSecret {
	e := &typesv1.ExternalSecret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ExternalSecret",
			APIVersion: "kubernetes-client.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: typesv1.ExternalSecretSpec{
			BackendType: backendType,
			Data:        data,
		},
	}

	if labels != nil || annotations != nil {
		e.Spec.Template = &typesv1.ExternalSecretTemplate{
			Metadata: typesv1.ExternalSecretTemplateMetadata{
				Annotations: annotations,
				Labels:      labels,
			},
		}
	}

	return e
}
