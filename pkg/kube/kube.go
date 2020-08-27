// Package kube implements a kubernetes client
package kube

import (
	"context"
	"encoding/json"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Kuber provides the methods that are available
// by a concrete implementation
type Kuber interface {
	Apply(fn ApplyFn, fns ...ApplyFn)
}

// Kube contains state for communicating with
// a kubernetes cluster
type Kube struct {
	KubeConfigPath string
	ClientSet      *kubernetes.Clientset
	Ctx            context.Context
}

// ApplyFn defines the signature of a function that applies
// some operation to the kubernetes cluster
type ApplyFn func(clientSet kubernetes.Interface) error

// New returns an initialised kubernetes client
func New(kubeConfigPath string) (*Kube, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Kube{
		KubeConfigPath: kubeConfigPath,
		ClientSet:      clientSet,
		Ctx:            context.Background(),
	}, nil
}

// Apply all the functions to the cluster
func (k *Kube) Apply(first ApplyFn, rest ...ApplyFn) error {
	fns := append([]ApplyFn{first}, rest...)

	for _, fn := range fns {
		err := fn(k.ClientSet)
		if err != nil {
			return err
		}
	}

	return nil
}

// Debug a namespace
func (k *Kube) Debug(namespace string) (map[string][]string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", k.KubeConfigPath)
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	events, err := clientSet.CoreV1().Events(namespace).List(k.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil
	}

	eventStrings := make([]string, len(events.Items))

	for i, event := range events.Items {
		j, err := json.MarshalIndent(event, "", "  ")
		if err != nil {
			return nil, err
		}

		eventStrings[i] = string(j)
	}

	pods, err := clientSet.CoreV1().Pods(namespace).List(k.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	podLogs := make([]string, len(pods.Items))
	podSpec := make([]string, len(pods.Items))

	for i, pod := range pods.Items {
		request := clientSet.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{})

		raw, err := request.DoRaw(k.Ctx)
		if err != nil {
			return nil, err
		}

		podLogs[i] = string(raw)

		j, err := json.MarshalIndent(pod, "", "  ")
		if err != nil {
			return nil, err
		}

		podSpec[i] = string(j)
	}

	return map[string][]string{
		"podLogs":  podLogs,
		"podSpecs": podSpec,
		"events":   eventStrings,
	}, nil
}
