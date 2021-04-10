// Package truncate implements truncating variables
package truncate

import "fmt"

const truncateFormat = "XXXtruncated%dbytesXXX"

// String truncates a string to the minimum of its length and the given max length
func String(s *string, maxLength int) string {
	if s == nil {
		return ""
	}

	// This way of doing substrings assumes ASCII and doesn't support UTF-8.
	// See https://stackoverflow.com/a/56129336
	truncateLength := min(maxLength, len(*s))
	truncated := (*s)[:truncateLength]

	if len(*s) > truncateLength {
		bytesTruncated := len(*s) - truncateLength
		truncated += fmt.Sprintf(truncateFormat, bytesTruncated)
	}

	return truncated
}

// Bytes truncates a byte array to the minimum of its length and the given max length
func Bytes(b []byte, maxLength int) []byte {
	if b == nil {
		return []byte{}
	}

	truncateLength := min(maxLength, len(b))

	truncated := b[:truncateLength]

	if len(b) > truncateLength {
		bytesTruncated := len(b) - truncateLength
		truncateInfo := fmt.Sprintf(truncateFormat, bytesTruncated)
		truncated = append(truncated, truncateInfo...)
	}

	return truncated
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}
