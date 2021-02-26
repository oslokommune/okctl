package scale

import (
	"context"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Scale struct {
	Name      string
	Namespace string
	Replicas  int32
	Ctx       context.Context
}

func New(name, namespace string, replicas int32) *Scale {
	return &Scale{
		Name:      name,
		Namespace: namespace,
		Replicas:  replicas,
		Ctx:       context.Background(),
	}
}

func (d *Scale) Scale(_ kubernetes.Interface, config *rest.Config) (interface{}, error) {
	client := kubernetes.NewForConfigOrDie(config)

	_, err := client.AppsV1().Deployments(d.Namespace).UpdateScale(
		d.Ctx,
		d.Name,
		&autoscalingv1.Scale{
			Spec: autoscalingv1.ScaleSpec{
				Replicas: d.Replicas,
			},
		},
		metav1.UpdateOptions{},
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
