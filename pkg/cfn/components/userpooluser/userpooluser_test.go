package userpooluser_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/userpooluser"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name     string
		golden   string
		resource cloudformation.Resource
	}{
		{
			name:     "Validate output",
			golden:   "user-pool-user.json",
			resource: userpooluser.New("testperson@origo.oslokommune.no", "Testperson", "ABCJE378FAKE").Resource(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.resource.(*cognito.UserPoolUser).MarshalJSON()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
