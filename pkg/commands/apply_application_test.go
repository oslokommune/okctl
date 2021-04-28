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

			err := WriteApplyApplicationSuccessMessage(&actual, tc.application, "infrastructure")
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, t.Name(), actual.Bytes())
		})
	}
}
