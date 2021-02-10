// Package namespace provides a namespace creator and applier
package namespace

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Namespace contains the state for creating a namespace
type Namespace struct {
	Namespace string
	Ctx       context.Context
}

// New returns an initialised namespace creator
func New(namespace string) *Namespace {
	return &Namespace{
		Namespace: namespace,
		Ctx:       context.Background(),
	}
}

// DeleteNamespace deletes the namespace
func (n *Namespace) DeleteNamespace(clientset kubernetes.Interface, _ *rest.Config) (interface{}, error) {
	return nil, clientset.CoreV1().Namespaces().Delete(n.Ctx, n.Namespace, metav1.DeleteOptions{})
}

// CreateNamespace creates the namespace
func (n *Namespace) CreateNamespace(clientset kubernetes.Interface, _ *rest.Config) (interface{}, error) {
	client := clientset.CoreV1().Namespaces()

	namespaces, err := client.List(n.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ns := range namespaces.Items {
		if ns.Name == n.Namespace {
			r, err := client.Get(n.Ctx, ns.Name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			return r, nil
		}
	}

	return client.Create(n.Ctx, n.NamespaceManifest(), metav1.CreateOptions{})
}

// NamespaceManifest returns the namespace manifest
func (n *Namespace) NamespaceManifest() *v1.Namespace {
	return &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: n.Namespace,
		},
	}
}
