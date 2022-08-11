package paths

import (
	"path"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

const (
	// DefaultApplicationsOutputDir is where the application declarations reside
	DefaultApplicationsOutputDir = "applications"
	// DefaultApplicationBaseDir is where the directory where application base files reside
	DefaultApplicationBaseDir = "base"
	// DefaultApplicationOverlayDir is where the directory where application overlay files reside
	DefaultApplicationOverlayDir = "overlays"
)

// GetRelativeApplicationDir knows how to construct the relative path to a specific applications root directory
func GetRelativeApplicationDir(cluster v1alpha1.Cluster, app v1alpha1.Application) string {
	return path.Join(cluster.Github.OutputPath, DefaultApplicationsOutputDir, app.Metadata.Name)
}
