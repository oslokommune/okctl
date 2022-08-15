package argocd_test

import (
	"errors"
	"path"
	"strings"
	"testing"

	"github.com/oslokommune/okctl/pkg/argocd"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const validApplication = `apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: mock-application
  namespace: argocd
spec:
  destination:
    namespace: mock-namespace
    server: https://kubernetes.default.svc
  project: default
  source:
    path: infrastructure/applications/mock-application/overlays/mock-cluster
    repoURL: git@github.com:mock-org/mock-iac-repo
    targetRevision: HEAD
  syncPolicy:
    automated:
      prune: false
      selfHeal: false

---
`

func TestIsArgoCDApplication(t *testing.T) {
	testCases := []struct {
		name           string
		withFs         *afero.Afero
		withTargetPath string
		expectBool     bool
		expectErr      error
	}{
		{
			name:           "Should return true with no error upon a valid ArgoCD application file",
			withTargetPath: path.Join("/", "my-application.yaml"),
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				err := fs.WriteReader(path.Join("/", "my-application.yaml"), strings.NewReader(validApplication))
				assert.NoError(t, err)

				return fs
			}(),
			expectBool: true,
			expectErr:  nil,
		},
		{
			name:           "Should return false and no error upon invalid file",
			withTargetPath: path.Join("/", "README.md"),
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				targetPath := path.Join("/", "README.md")

				err := fs.WriteReader(targetPath, strings.NewReader("# A title\n## Header\n\nThis is not an Argo app"))
				assert.NoError(t, err)

				return fs
			}(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := argocd.IsArgoCDApplication(tc.withFs, tc.withTargetPath)

			assert.Equal(t, tc.expectBool, result)
			assert.True(t, errors.Is(tc.expectErr, err))
		})
	}
}
