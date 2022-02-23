package client

// KubernetesState defines functionality for handling state in Kubernetes
type KubernetesState interface {
	HasResource(kind, namespace, name string) (bool, error)
}
