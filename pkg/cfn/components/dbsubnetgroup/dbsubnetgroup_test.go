package dbsubnetgroup_test

import (
	"testing"

	"github.com/awslabs/goformation/v4/cloudformation/rds"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components/dbsubnetgroup"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

type testSubnet struct{}

func (t testSubnet) Ref() string {
	return cloudformation.Ref("DbSubnetSomething")
}

func TestNew(t *testing.T) {
	testCases := []struct {
		name     string
		golden   string
		resource cloudformation.Resource
	}{
		{
			name:   "Validate output",
			golden: "dbsubnetgroup.json",
			resource: dbsubnetgroup.New(
				[]cfn.Referencer{testSubnet{}},
			).Resource(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.resource.(*rds.DBSubnetGroup).MarshalJSON()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
