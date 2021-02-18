package storageclass_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/kube/manifests/storageclass"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name        string
		golden      string
		storage     string
		params      *storageclass.EBSParameters
		annotations map[string]string
		expect      interface{}
		expectErr   bool
	}{
		{
			name:        "Should work",
			golden:      "storageclass.yaml",
			storage:     "ebs-sc",
			params:      storageclass.NewEBSParameters(),
			annotations: storageclass.DefaultStorageClassAnnotation(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t)

			got, err := storageclass.New(tc.storage, tc.params, tc.annotations)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				dm := got.StorageClassManifest()
				db, err := yaml.Marshal(dm)
				assert.NoError(t, err)
				g.Assert(t, tc.golden, db)
			}
		})
	}
}
