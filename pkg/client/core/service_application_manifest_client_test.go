package core

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/jsonpatch"
	"github.com/oslokommune/okctl/pkg/scaffold/resources"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

type expectedFile struct {
	Filename string
	Content  []byte
}

//nolint:funlen
func TestApplicationManifestService_SaveManifest(t *testing.T) {
	testCases := []struct {
		name string

		withOpts    client.SaveManifestOpts
		expectFiles []expectedFile
	}{
		{
			name: "Should work",
			withOpts: client.SaveManifestOpts{
				ApplicationName: "test",
				Filename:        "service.yaml",
				Content: func() []byte {
					raw, _ := yaml.Marshal(resources.CreateOkctlService(
						v1alpha1.Application{Metadata: v1alpha1.ApplicationMeta{Name: "test"}},
						"main",
					))

					return raw
				}(),
			},
			expectFiles: []expectedFile{
				{
					Filename: "/test/base/kustomization.yaml",
					Content: func() []byte {
						raw, _ := yaml.Marshal(resources.Kustomization{
							Resources: []string{"service.yaml"},
						})

						return raw
					}(),
				},
				{
					Filename: "/test/base/service.yaml",
					Content: func() []byte {
						raw, _ := yaml.Marshal(resources.CreateOkctlService(
							v1alpha1.Application{Metadata: v1alpha1.ApplicationMeta{Name: "test"}},
							"main",
						))

						return raw
					}(),
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			fs := &afero.Afero{Fs: afero.NewMemMapFs()}
			service := NewApplicationManifestService(fs, "/")

			err := service.SaveManifest(context.Background(), tc.withOpts)
			assert.NoError(t, err)

			for _, f := range tc.expectFiles {
				content, err := fs.ReadFile(f.Filename)
				assert.NoError(t, err)

				assert.Equal(t, f.Content, content)
			}
		})
	}
}

//nolint:funlen
func TestApplicationManifestService_SavePatch(t *testing.T) {
	testCases := []struct {
		name string

		withOpts    client.SavePatchOpts
		expectFiles []expectedFile
	}{
		{
			name: "Should work",
			withOpts: client.SavePatchOpts{
				ApplicationName: "testapp",
				ClusterName:     "testcluster",
				Kind:            "Service",
				Patch: jsonpatch.Patch{
					Operations: []jsonpatch.Operation{
						{
							Type:  jsonpatch.OperationTypeAdd,
							Path:  "/spec/nothing",
							Value: "0",
						},
					},
				},
			},
			expectFiles: []expectedFile{
				{
					Filename: "/testapp/overlays/testcluster/kustomization.yaml",
					Content: func() []byte {
						raw, _ := yaml.Marshal(resources.Kustomization{
							Resources: []string{"../../base"},
							Patches: []resources.PatchReference{
								{
									Path: "service-patch.json",
									Target: resources.PatchTarget{
										Kind: "Service",
									},
								},
							},
						})

						return raw
					}(),
				},
				{
					Filename: "/testapp/overlays/testcluster/service-patch.json",
					Content: func() []byte {
						raw, _ := json.Marshal(jsonpatch.Patch{
							Operations: []jsonpatch.Operation{
								{
									Type:  jsonpatch.OperationTypeAdd,
									Path:  "/spec/nothing",
									Value: "0",
								},
							},
						})

						return raw
					}(),
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			fs := &afero.Afero{Fs: afero.NewMemMapFs()}
			service := NewApplicationManifestService(fs, "/")

			err := service.SavePatch(context.Background(), tc.withOpts)
			assert.NoError(t, err)

			for _, f := range tc.expectFiles {
				content, err := fs.ReadFile(f.Filename)
				assert.NoError(t, err)

				assert.Equal(t, f.Content, content)
			}
		})
	}
}
