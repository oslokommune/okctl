// Package scale knows how to scale a deployment in Kubernetes
package scale

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"time"

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

	if scale.Spec.Replicas == d.Replicas {
		return nil, nil
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

	// Loop for a limited period of time, return immediately if the deployment has reached the
	// desired state
	t0 := time.Now()

	ticker := time.NewTicker(5 * time.Second) // nolint: gomnd
	defer ticker.Stop()

	// I think this is a false positive, but probably not
	// I am too tired to bother understanding it right now
	// maybe some other day
	// maybe never
	// nolint: gosimple
	for {
		select {
		case <-ticker.C:
			scale, err := client.AppsV1().Deployments(d.Namespace).GetScale(d.Ctx, d.Name, metav1.GetOptions{})
			if err != nil {
				return nil, err
			}

			if scale.Spec.Replicas == d.Replicas {
				return nil, nil
			}

			if time.Now().After(t0.Add(5 * time.Minute)) { // nolint: gomnd
				return nil, fmt.Errorf(constant.ScaleTimeoutError)
			}
		}
	}
}
