package scaffold

import (
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/jsonpatch"
)

// ManifestSaver defines a function which store manifests
type ManifestSaver func(filename string, content []byte) error

// PatchSaver defines a function which store patches
type PatchSaver func(kind string, patch jsonpatch.Patch) error

// GenerateApplicationBaseOpts contains required data to generate application base manifests
type GenerateApplicationBaseOpts struct {
	SaveManifest ManifestSaver
	Application  v1alpha1.Application
}

// GenerateApplicationOverlayOpts contains required data to generate application patches
type GenerateApplicationOverlayOpts struct {
	SavePatch   PatchSaver
	Application v1alpha1.Application

	Domain         string
	CertificateARN string
}

// GenerateArgoCDApplicationManifestOpts contains required information to generate an ArgoCD Application Manifest
type GenerateArgoCDApplicationManifestOpts struct {
	Name          string
	Namespace     string
	IACRepoURL    string
	SourceSyncDir string
}
