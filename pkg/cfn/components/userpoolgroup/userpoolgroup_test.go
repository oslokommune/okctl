package userpoolgroup_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/userpoolgroup"
	"github.com/sebdah/goldie/v2"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
	"github.com/oslokommune/okctl/pkg/cfn/components/userpool"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	up := userpool.New("test")

	testCases := []struct {
		name     string
		golden   string
		resource cloudformation.Resource
	}{
		{
			name:     "Validate output",
			golden:   "user-group.json",
			resource: userpoolgroup.New("group", "description", up).Resource(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.resource.(*cognito.UserPoolGroup).MarshalJSON()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
