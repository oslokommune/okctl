package git_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/git"
	"github.com/stretchr/testify/assert"
)

func TestGithubRemotes(t *testing.T) {
	testCases := []struct {
		name   string
		path   string
		expect string
	}{
		{
			name:   "Should work",
			path:   "../..",
			expect: "oslokommune/okctl",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := git.GithubRepoFullName("oslokommune", tc.path)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, got)
		})
	}
}
