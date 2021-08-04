package commands

import (
	"fmt"
	"testing"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/stretchr/testify/assert"
)

func TestCleanOsEnvVars(t *testing.T) {
	testCases := []struct {
		name            string
		withVariables   []string
		expectVariables []string
	}{
		{
			name: "Should do nothing with no blacklisted variables",

			withVariables:   []string{"PATH=/bin", "EDITOR=vim"},
			expectVariables: []string{"PATH=/bin", "EDITOR=vim"},
		},
		{
			name: "Should remove blacklisted variables",

			withVariables: []string{
				"PATH=/bin",
				fmt.Sprintf("%s_%s=somevalue", constant.EnvPrefix, constant.EnvClusterDeclaration),
				"EDITOR=vim",
			},
			expectVariables: []string{"EDITOR=vim", "PATH=/bin"},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result := CleanOsEnvVars(tc.withVariables)

			assert.Equal(t, len(tc.expectVariables), len(result))

			for _, expectedItem := range tc.expectVariables {
				assert.True(t, contains(result, expectedItem))
			}
		})
	}
}

func contains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}
