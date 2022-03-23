package core

import (
	"context"
	"fmt"
	"io"
	"path"
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clients/kubectl"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSetupNamespaceSync(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "Should produce an ArgoCD application manifest for namespaces in the correct location",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			fs := &afero.Afero{Fs: afero.NewMemMapFs()}
			absoluteRepositoryDir := "/"

			cluster := v1alpha1.NewCluster()
			cluster.Metadata.Name = "mockCluster"
			cluster.Github.Repository = "mock-iac"

			service := NewArgoCDService(NewArgoCDServiceOpts{
				Fs:              fs,
				AbsoluteRepoDir: absoluteRepositoryDir,
			})

			err := service.SetupNamespacesSync(ctx, &mockKubectlClient{}, cluster)
			assert.NoError(t, err)

			argocdManifestPath := path.Join(
				absoluteRepositoryDir,
				cluster.Github.OutputPath,
				cluster.Metadata.Name,
				constant.DefaultArgoCDClusterConfigDir,
				fmt.Sprintf("%s.yaml", defaultArgoCDNamespacesManifestName),
			)

			result, err := fs.ReadFile(argocdManifestPath)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, result)
		})
	}
}

func TestSetupApplicationSync(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "Should produce an ArgoCD application manifest for applications in the correct location",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			fs := &afero.Afero{Fs: afero.NewMemMapFs()}
			absoluteRepositoryDir := "/"

			cluster := v1alpha1.NewCluster()
			cluster.Metadata.Name = "mockCluster"
			cluster.Github.Repository = "mock-iac"

			service := NewArgoCDService(NewArgoCDServiceOpts{
				Fs:              fs,
				AbsoluteRepoDir: absoluteRepositoryDir,
			})

			err := service.SetupApplicationsSync(ctx, &mockKubectlClient{}, cluster)
			assert.NoError(t, err)

			argocdManifestPath := path.Join(
				absoluteRepositoryDir,
				cluster.Github.OutputPath,
				cluster.Metadata.Name,
				constant.DefaultArgoCDClusterConfigDir,
				fmt.Sprintf("%s.yaml", defaultArgoCDApplicationManifestName),
			)

			result, err := fs.ReadFile(argocdManifestPath)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, result)
		})
	}
}

type mockKubectlClient struct{}

func (m mockKubectlClient) Apply(_ io.Reader) error                 { return nil }
func (m mockKubectlClient) Delete(_ io.Reader) error                { panic("implement me") }
func (m mockKubectlClient) Patch(_ kubectl.PatchOpts) error         { panic("implement me") }
func (m mockKubectlClient) Exists(_ kubectl.Resource) (bool, error) { panic("implement me") }
