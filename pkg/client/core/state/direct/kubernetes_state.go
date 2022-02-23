package direct

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/clients/kubectl"
)

type kubernetesState struct {
	kubectl kubectl.Client
}

// HasResource returns a boolean indicating if a resource exists or not
func (k kubernetesState) HasResource(kind, namespace, name string) (bool, error) {
	exists, err := k.kubectl.Exists(kubectl.Resource{
		Namespace: namespace,
		Kind:      kind,
		Name:      name,
	})
	if err != nil {
		return false, fmt.Errorf("checking existence: %w", err)
	}

	return exists, nil
}

// NewKubernetesState returns an initialized Kubernetes state client
func NewKubernetesState(kubectlClient kubectl.Client) client.KubernetesState {
	return &kubernetesState{
		kubectl: kubectlClient,
	}
}
