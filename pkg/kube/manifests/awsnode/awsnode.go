// Package awsnode performs operations on the aws-node daemonset
package awsnode

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

// AWSNode contains the required state for
// interacting with the aws-node
type AWSNode struct {
	Name      string
	Namespace string
	client    kubernetes.Interface
	ctx       context.Context
}

// New returns an initialised aws-node client
func New(client kubernetes.Interface) *AWSNode {
	return &AWSNode{
		Name:      "aws-node",
		Namespace: "kube-system",
		client:    client,
		ctx:       context.Background(),
	}
}

// EnablePodENI enables the pod ENI's so we can attach
// security groups to pods
// Based on the output from this command, this appears to be what kubectl does:
// `kubectl -n kube-system set env daemonset aws-node ENABLE_POD_ENI=true --v=9`
func (a *AWSNode) EnablePodENI() error {
	ds, err := a.client.AppsV1().DaemonSets(a.Namespace).Get(a.ctx, a.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	for i, c := range ds.Spec.Template.Spec.Containers {
		if c.Name == a.Name {
			for j, e := range c.Env {
				if e.Name == "ENABLE_POD_ENI" {
					e.Value = "true"
					ds.Spec.Template.Spec.Containers[i].Env[j] = e
				}
			}
		}
	}

	_, err = a.client.AppsV1().DaemonSets(a.Namespace).Update(a.ctx, ds, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
