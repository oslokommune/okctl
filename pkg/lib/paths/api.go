/*
Package paths api.go exposes commonly used paths throughout okctl.

If you require the absolute path to somewhere, combine the relative helper with GetABsoluteIACRepositoryRootDirectory
using path.Join.
*/
package paths

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// GetAbsoluteIACRepositoryRootDirectory returns the absolute path of the repository root no matter what the current
// working directory of the repository the user is in
// I.e.: /home/olly/.../team-iac/
func GetAbsoluteIACRepositoryRootDirectory() (string, error) {
	result, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("getting repository root directory: %w", err)
	}

	pathAsString := string(bytes.Trim(result, "\n"))

	return pathAsString, nil
}

// GetRelativeClusterOutputDirectory returns the relative path of the cluster output directory
// I.e.: /<output directory>/<cluster name>/
func GetRelativeClusterOutputDirectory(clusterManifest v1alpha1.Cluster) string {
	return path.Join(clusterManifest.Github.OutputPath, clusterManifest.Metadata.Name)
}

// GetRelativeClusterApplicationsDirectory returns the relative path of the cluster specific ArgoCD applications
// manifest directory
// I.e.: /<output directory>/<cluster name>/argocd/applications/
func GetRelativeClusterApplicationsDirectory(clusterManifest v1alpha1.Cluster) string {
	return path.Join(
		GetRelativeClusterOutputDirectory(clusterManifest),
		DefaultClusterArgoCDConfigDirectoryName,
		DefaultClusterArgoCDApplicationsDirectoryName,
	)
}

// GetRelativeClusterOkctlConfigurationDirectory returns the relative path of the okctl configuration directory. This is
// the directory where we place configuration applied to the environment that okctl manage.
// I.e.: /<output directory>/<cluster name>/okctl/
func GetRelativeClusterOkctlConfigurationDirectory(clusterManifest v1alpha1.Cluster) string {
	return path.Join(
		GetRelativeClusterOutputDirectory(clusterManifest),
		DefaultClusterOkctlConfigurationDirectoryName,
	)
}
