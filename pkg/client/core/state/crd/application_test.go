package crd

import (
	"io"
	"path"
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clients/kubectl"
	"github.com/oslokommune/okctl/pkg/lib/paths"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInitializeState(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	clusterManifest := v1alpha1.NewCluster()
	absoluteIACRepositoryRootDirectory := "/"

	clusterManifest.Metadata.Name = "mock-cluster"
	clusterManifest.Github.Repository = "mock-repo"

	state := NewApplicationState(fs, &mockKubectl{})

	err := state.Initialize(clusterManifest, absoluteIACRepositoryRootDirectory)
	assert.NoError(t, err)

	okctlConfigDir := path.Join(
		absoluteIACRepositoryRootDirectory,
		paths.GetRelativeClusterOkctlConfigurationDirectory(clusterManifest),
	)

	applicationsDir := path.Join(
		absoluteIACRepositoryRootDirectory,
		paths.GetRelativeClusterApplicationsDirectory(clusterManifest),
	)

	g := goldie.New(t)

	assertExistence(t, g, fs, "application-crd", path.Join(okctlConfigDir, defaultApplicationCustomResourceDefinitionFilename))
	assertExistence(t, g, fs, "okctl-config-argoapp", path.Join(applicationsDir, defaultOkctlConfigurationDirArgoCDApplicationManifestFilename))
}

func assertExistence(t *testing.T, g *goldie.Goldie, fs *afero.Afero, name string, path string) {
	raw, err := fs.ReadFile(path)
	assert.NoError(t, err)

	g.Assert(t, name, raw)
}

type mockKubectl struct{}

func (m mockKubectl) Get(kubectl.Resource) (io.Reader, error) { panic("implement me") }
func (m mockKubectl) Apply(io.Reader) error                   { panic("implement me") }
func (m mockKubectl) DeleteByManifest(io.Reader) error        { panic("implement me") }
func (m mockKubectl) DeleteByResource(kubectl.Resource) error { panic("implement me") }
func (m mockKubectl) Patch(kubectl.PatchOpts) error           { panic("implement me") }
func (m mockKubectl) Exists(kubectl.Resource) (bool, error)   { panic("implement me") }
