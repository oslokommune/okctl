package kube

import (
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
