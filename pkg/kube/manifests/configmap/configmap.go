// Package configmap provides a configmap creator and applier
package configmap

import (
	"context"

	v1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ConfigMap contains the state for creating a configmap
type ConfigMap struct {
	Name      string
	Namespace string
	Manifest  *v1.ConfigMap
	Ctx       context.Context
}

// New returns an initialised ConfigMap creator
func New(name, namespace string, manifest *v1.ConfigMap) *ConfigMap {
	return &ConfigMap{
		Name:      name,
		Namespace: namespace,
		Manifest:  manifest,
		Ctx:       context.Background(),
	}
}

// DeleteConfigMap deletes the ConfigMap
func (s *ConfigMap) DeleteConfigMap(_ kubernetes.Interface, config *rest.Config) (interface{}, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	ns, err := client.CoreV1().ConfigMaps(s.Namespace).List(s.Ctx, metav1.ListOptions{})
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

	return nil, client.CoreV1().ConfigMaps(s.Namespace).Delete(s.Ctx, s.Name, metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
}

// CreateConfigMap creates the ConfigMap
func (s *ConfigMap) CreateConfigMap(clientset kubernetes.Interface, _ *rest.Config) (interface{}, error) {
	client := clientset.CoreV1().ConfigMaps(s.Namespace)

	ConfigMaps, err := client.List(s.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, ns := range ConfigMaps.Items {
		if ns.Name == s.Name {
			r, err := client.Get(s.Ctx, ns.Name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			return r, nil
		}
	}

	return client.Create(s.Ctx, s.Manifest, metav1.CreateOptions{})
}

// NewManifest returns the ConfigMap manifest
func NewManifest(name, namespace string, data, labels map[string]string) *v1.ConfigMap {
	return &v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: data,
	}
}
