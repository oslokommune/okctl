package components_test

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/oslokommune/okctl/pkg/cfn/components"
)

func TestRootDomain(t *testing.T) {
	testCases := []struct {
		name   string
		domain string
		expect string
	}{
		{
			name:   "Basic domain",
			domain: "test.oslo.systems",
			expect: "oslo.systems",
		},
		{
			name:   "No sub-domain",
			domain: "oslo.systems",
			expect: "oslo.systems",
		},
		{
			name:   "Deeply nested",
			domain: "auth.test.oslo.systems",
			expect: "test.oslo.systems",
		},
		{
			name:   "Just a word",
			domain: "systems",
			expect: "systems",
		},
		{
			name:   "Empty",
			domain: "",
			expect: "",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got := components.RootDomain(tc.domain)
			assert.Equal(t, tc.expect, got)
		})
	}
}
