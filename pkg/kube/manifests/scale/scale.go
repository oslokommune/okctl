// Package scale knows how to scale a deployment in Kubernetes
package scale

import (
	"context"

	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Scale contains the required state for scaling a deployment
type Scale struct {
	Name      string
	Namespace string
	Replicas  int32
	Ctx       context.Context
}

// New returns an initialised deployment scaler
func New(name, namespace string, replicas int32) *Scale {
	return &Scale{
		Name:      name,
		Namespace: namespace,
		Replicas:  replicas,
		Ctx:       context.Background(),
	}
}

// Scale the deployment with the provided number of replicas
func (d *Scale) Scale(_ kubernetes.Interface, config *rest.Config) (interface{}, error) {
	client := kubernetes.NewForConfigOrDie(config)

	scale, err := client.AppsV1().Deployments(d.Namespace).GetScale(d.Ctx, d.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	scale.Spec.Replicas = d.Replicas

	_, err = client.AppsV1().Deployments(d.Namespace).UpdateScale(
		d.Ctx,
		d.Name,
		scale,
		metav1.UpdateOptions{},
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
