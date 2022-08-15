package paths

import (
	"path"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

const (
	// DefaultArgoCDClusterConfigDir is where we do cluster specific ArgoCD configuration
	DefaultArgoCDClusterConfigDir = "argocd"
	// DefaultArgoCDClusterConfigApplicationsDir is where we put ArgoCD application manifests for applications
	DefaultArgoCDClusterConfigApplicationsDir = "applications"
	// DefaultArgoCDClusterConfigNamespacesDir is where we put namespace manifests for applications
	DefaultArgoCDClusterConfigNamespacesDir = "namespaces"
)

// GetRelativeArgoCDApplicationsDir knows how to construct the relative path to the directory synched by ArgoCD where
// we place ArgoCD applications representing active applications in the cluster
func GetRelativeArgoCDApplicationsDir(cluster v1alpha1.Cluster) string {
	return path.Join(
		cluster.Github.OutputPath,
		cluster.Metadata.Name,
		DefaultArgoCDClusterConfigDir,
		DefaultArgoCDClusterConfigApplicationsDir,
	)
}
