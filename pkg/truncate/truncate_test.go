package truncate_test

import (
	"github.com/oslokommune/okctl/pkg/truncate"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTruncate(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{
			name:      "Should truncate string",
			input:     "1234567890",
			maxLength: 5,
			expected:  "12345",
		},
		{
			name:      "Should keep string if it's equal to maxLength",
			input:     "1234567890",
			maxLength: 10,
			expected:  "1234567890",
		},
		{
			name:      "Should keep string if it's over maxLength",
			input:     "1234567890",
			maxLength: 11,
			expected:  "1234567890",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Test Truncate string
			truncatedString := truncate.Truncate(&tc.input, tc.maxLength)
			assert.Equal(t, tc.expected, truncatedString)

			// Test Truncate bytes
			inputBytes := []byte(tc.input)
			truncatedBytes := truncate.TruncateBytes(inputBytes, tc.maxLength)
			expectedBytes := []byte(tc.expected)
			assert.Equal(t, expectedBytes, truncatedBytes)
		})
	}
}
