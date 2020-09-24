package userpooldomain_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/recordset"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
	"github.com/oslokommune/okctl/pkg/cfn/components/userpool"
	"github.com/oslokommune/okctl/pkg/cfn/components/userpooldomain"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	up := userpool.New("test", "test")
	ph := recordset.New("placeholder", "1.1.1.1", "auth.test.com", "GHFJE378FAKE")

	testCases := []struct {
		name     string
		golden   string
		resource cloudformation.Resource
	}{
		{
			name:     "Validate output",
			golden:   "user-pool.json",
			resource: userpooldomain.New("auth.oslo.systems", "arn://certificate/HAF93FAKE", up, ph).Resource(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.resource.(*cognito.UserPoolDomain).MarshalJSON()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
