package aliasrecordset_test

import (
	"testing"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/route53"
	"github.com/oslokommune/okctl/pkg/cfn/components/aliasrecordset"
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
			name:   "Validate output",
			golden: "alias-recordset.json",
			resource: aliasrecordset.New(
				"DomainPoolAuth",
				"cloudfront-us-east-1.aws.com",
				"HJOJF678FAKE",
				"auth.oslo.systems",
				"GHFJE78FAKE",
			).Resource(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.resource.(*route53.RecordSet).MarshalJSON()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
