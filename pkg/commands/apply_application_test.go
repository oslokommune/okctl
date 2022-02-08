package commands

import (
	"bytes"
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestApplyApplicationSuccessMessage(t *testing.T) {
	testCases := []struct {
		name        string
		application v1alpha1.Application
	}{
		{
			name: "Should get expected success message when using image URI",
			application: v1alpha1.Application{
				Metadata: v1alpha1.ApplicationMeta{
					Name: "my-app",
				},
				Image: v1alpha1.ApplicationImage{
					URI: "ubuntu",
				},
			},
		},
		{
			name: "Should get expected success message when using image name",
			application: v1alpha1.Application{
				Metadata: v1alpha1.ApplicationMeta{
					Name: "my-app",
				},
				Image: v1alpha1.ApplicationImage{
					Name: "my-image",
				},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var actual bytes.Buffer

			err := WriteApplyApplicationSuccessMessage(WriteApplyApplicationSucessMessageOpts{
				Out:         &actual,
				Application: tc.application,
				Cluster: v1alpha1.Cluster{
					Metadata: v1alpha1.ClusterMeta{
						Name: "test-cluster",
					},
					Github: v1alpha1.ClusterGithub{
						OutputPath: "infrastructure",
					},
				},
			})
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, t.Name(), actual.Bytes())
		})
	}
}

func TestDeleteApplicationSucessMessage(t *testing.T) {
	testCases := []struct {
		name            string
		withCluster     v1alpha1.Cluster
		withApplication v1alpha1.Application
	}{
		{
			name:            "Should produce expected success message",
			withCluster:     v1alpha1.Cluster{},
			withApplication: v1alpha1.Application{Metadata: v1alpha1.ApplicationMeta{Name: "test"}},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			var result bytes.Buffer

			err := WriteDeleteApplicationSuccessMessage(WriteDeleteApplicationSuccessMessageOpts{
				Out:         &result,
				Cluster:     tc.withCluster,
				Application: tc.withApplication,
			})
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, t.Name(), result.Bytes())
		})
	}
}
