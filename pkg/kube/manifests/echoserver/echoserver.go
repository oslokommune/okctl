// Package echoserver provides a simple echoserver
package echoserver

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// EchoServer contains the state for apply the external-dns
// manifests to kubernetes
type EchoServer struct {
	Namespace     string
	Version       string
	Replicas      int32
	ContainerPort int32
	Ctx           context.Context
}

// New returns an initialised nginx deployment
func New(namespace string) *EchoServer {
	return &EchoServer{
		Namespace:     namespace,
		Version:       "1.10",
		ContainerPort: 8080, // nolint: gomnd
		Replicas:      1,
		Ctx:           context.Background(),
	}
}

// CreateDeployment creates the echoserver deployment manifest
func (e *EchoServer) CreateDeployment(clientSet kubernetes.Interface) (interface{}, error) {
	deployClient := clientSet.AppsV1().Deployments(e.Namespace)
	return deployClient.Create(e.Ctx, e.DeploymentManifest(), metav1.CreateOptions{})
}

// DeploymentManifest returns the deployment manifest
func (e *EchoServer) DeploymentManifest() *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "echoserver",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &e.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "echoserver",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "echoserver",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "echoserver",
							Image:           fmt.Sprintf("gcr.io/google-containers/echoserver:%s", e.Version),
							ImagePullPolicy: "Always",
							Ports: []v1.ContainerPort{
								{
									ContainerPort: e.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}
}

//apiVersion: v1
//kind: Service
//metadata:
//  name: echoserver
//spec:
//  ports:
//  - port: 80
//    targetPort: 8080
//    protocol: TCP
//  selector:
//    app: echoserver

//apiVersion: extensions/v1beta1
//kind: Ingress
//metadata:
//  name: echoserver
//  annotations:
//    kubernetes.io/ingress.class: "nginx"
//spec:
//  rules:
//  - host: echo.example.com
//    http:
//      paths:
//      - path: /
//        backend:
//          serviceName: echoserver
//          servicePort: 80
