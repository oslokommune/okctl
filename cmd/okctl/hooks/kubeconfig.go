package hooks

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// WriteKubeConfig writes the Kubernetes configuration to disk
func WriteKubeConfig(o *okctl.Okctl) RunEer {
	return func(_ *cobra.Command, _ []string) error {
		kubeconfigStore, err := o.KubeConfigStore()
		if err != nil {
			return fmt.Errorf("getting kubeconfig store: %w", err)
		}

		// GetKubeConfig actually writes the config to disk
		_, err = kubeconfigStore.GetKubeConfig(o.Declaration.Metadata.Name)
		if err != nil {
			return fmt.Errorf("getting kubeconfig: %w", err)
		}

		return nil
	}
}
