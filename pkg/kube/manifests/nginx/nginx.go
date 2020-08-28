// Package nginx provides kubernetes manifests for deploying nginx
// this is primarily used for testing kube
package nginx

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Nginx contains the state for apply the external-dns
// manifests to kubernetes
type Nginx struct {
	Namespace    string
	DomainFilter string
	Version      string
	OwnerID      string
	FsGroup      int64
	RunAsNonRoot bool
	Replicas     int32
	Ctx          context.Context
}

// New returns an initialised nginx deployment
func New(namespace string) *Nginx {
	return &Nginx{
		Namespace: namespace,
		Version:   "1.14.2",
		Replicas:  1,
		Ctx:       context.Background(),
	}
}

// CreateDeployment creates the external-dns Deployment manifest
func (e *Nginx) CreateDeployment(clientSet kubernetes.Interface) (interface{}, error) {
	deployClient := clientSet.AppsV1().Deployments(e.Namespace)
	return deployClient.Create(e.Ctx, e.DeploymentManifest(), metav1.CreateOptions{})
}

// DeploymentManifest returns the deployment manifest
func (e *Nginx) DeploymentManifest() *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &e.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "nginx",
							Image: fmt.Sprintf("nginx:%s", e.Version),
						},
					},
				},
			},
		},
	}
}
