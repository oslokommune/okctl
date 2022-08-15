package core

import (
	"context"
	"fmt"
	"io"
	"path"
	"testing"

	"github.com/oslokommune/okctl/pkg/paths"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/clients/kubectl"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSetupNamespaceSync(t *testing.T) {
	testCases := []struct {
		name                   string
		expectedKubectlApplies int
	}{
		{
			name:                   "Should produce an ArgoCD application manifest for namespaces in the correct location",
			expectedKubectlApplies: 1,
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

			kubectlClient := &mockKubectlClient{}

			err := service.SetupNamespacesSync(ctx, kubectlClient, cluster)
			assert.NoError(t, err)

			argocdManifestPath := path.Join(
				absoluteRepositoryDir,
				cluster.Github.OutputPath,
				cluster.Metadata.Name,
				paths.DefaultArgoCDClusterConfigDir,
				fmt.Sprintf("%s.yaml", defaultArgoCDNamespacesManifestName),
			)

			result, err := fs.ReadFile(argocdManifestPath)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, result)

			assert.Equal(t, tc.expectedKubectlApplies, kubectlClient.numberOfApplies)
		})
	}
}

func TestSetupApplicationSync(t *testing.T) {
	testCases := []struct {
		name                   string
		expectedKubectlApplies int
	}{
		{
			name:                   "Should produce an ArgoCD application manifest for applications in the correct location",
			expectedKubectlApplies: 1,
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

			kubectlClient := &mockKubectlClient{}

			err := service.SetupApplicationsSync(ctx, kubectlClient, cluster)
			assert.NoError(t, err)

			argocdManifestPath := path.Join(
				absoluteRepositoryDir,
				cluster.Github.OutputPath,
				cluster.Metadata.Name,
				paths.DefaultArgoCDClusterConfigDir,
				fmt.Sprintf("%s.yaml", defaultArgoCDApplicationManifestName),
			)

			result, err := fs.ReadFile(argocdManifestPath)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, result)

			assert.Equal(t, tc.expectedKubectlApplies, kubectlClient.numberOfApplies)
		})
	}
}

type mockKubectlClient struct {
	numberOfApplies int
}

func (m *mockKubectlClient) Apply(_ io.Reader) error {
	m.numberOfApplies++

	return nil
}

func (m *mockKubectlClient) Delete(_ io.Reader) error                { panic("implement me") }
func (m *mockKubectlClient) Patch(_ kubectl.PatchOpts) error         { panic("implement me") }
func (m *mockKubectlClient) Exists(_ kubectl.Resource) (bool, error) { panic("implement me") }
