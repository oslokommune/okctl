package truncate

// Truncate truncates a string to the minimum of its length and the given max length
func Truncate(s *string, maxLength int) string {
	// This way of doing substrings assumes ASCII and doesn't support UTF-8.
	// See https://stackoverflow.com/a/56129336
	truncateLength := min(maxLength, len(*s))
	return (*s)[:truncateLength]
}
// Truncate truncates a byte array to the minimum of its length and the given max length
func TruncateBytes(b []byte, maxLength int) []byte {
	// This way of doing substrings assumes ASCII and doesn't support UTF-8.
	// See https://stackoverflow.com/a/56129336
	truncateLength := min(maxLength, len(b))
	return b[:truncateLength]
}

func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}
