package manager_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/manager"
	"github.com/oslokommune/okctl/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	testCases := []struct {
		name        string
		manager     *manager.Manager
		expect      interface{}
		expectError bool
	}{
		{
			name: "Should work",
			manager: manager.
				New(
					mock.NewCloudProvider().
						DescribeStacksEmpty().
						CreateStackSuccess().
						DescribeStacksResponse(),
				).
				WithBuilder(mock.NewGoodBuilder()),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := tc.manager.CreateIfNotExists(10)
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
