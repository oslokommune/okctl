package userpoolusertogroupattachment_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/userpooluser"
	"github.com/oslokommune/okctl/pkg/cfn/components/userpoolusertogroupattachment"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	user := userpooluser.New("testperson@origo.oslo.kommune.no", "Testperson", "ABCJE378FAKE")

	testCases := []struct {
		name     string
		golden   string
		resource cloudformation.Resource
	}{
		{
			name:     "Validate output",
			golden:   "user-pool-user-to-group-attachment.json",
			resource: userpoolusertogroupattachment.New(user, "Username1337", "somegroup", "XXXXXXFAKE").Resource(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.resource.(*cognito.UserPoolUserToGroupAttachment).MarshalJSON()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
