package truncate_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/truncate"
	"github.com/stretchr/testify/assert"
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
			expected:  "12345XXXtruncated5bytesXXX",
		},
		{
			name:      "Should truncate some other string",
			input:     "1234567890",
			maxLength: 7,
			expected:  "1234567XXXtruncated3bytesXXX",
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
			// Test String string
			truncatedString := truncate.String(&tc.input, tc.maxLength)
			assert.Equal(t, tc.expected, truncatedString)

			// Test String bytes
			inputBytes := []byte(tc.input)
			truncatedBytes := truncate.Bytes(inputBytes, tc.maxLength)

			expectedBytes := []byte(tc.expected)
			assert.Equal(t, expectedBytes, truncatedBytes)
		})
	}
}

func TestTruncateNil(t *testing.T) {
	t.Run("Should return empty values if receiving nil", func(t *testing.T) {
		assert.Equal(t, "", truncate.String(nil, 5))
		assert.Equal(t, []byte{}, truncate.Bytes(nil, 5))
	})
}

func TestTruncateSideEffects(t *testing.T) {
	t.Run("Should not modify original byte array", func(t *testing.T) {
		b := []byte("hello")
		expected := []byte("hello")

		truncate.Bytes(b, 3)

		assert.Equal(t, expected, b)
	})
}
