package managedpolicy_test

import (
	"testing"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/iam"
	"github.com/oslokommune/okctl/pkg/cfn/components/managedpolicy"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	document := struct {
		Name string `json:"Document"`
	}{
		Name: "TestDocument",
	}

	testCases := []struct {
		name     string
		golden   string
		resource cloudformation.Resource
	}{
		{
			name:   "Valid output",
			golden: "managed-policy.json",
			resource: managedpolicy.New(
				"MyPolicy",
				"somePolicyName",
				"some desc",
				&document).Resource(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.resource.(*iam.ManagedPolicy).MarshalJSON()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
