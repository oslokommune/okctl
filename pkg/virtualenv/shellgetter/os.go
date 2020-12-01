package shellgetter

// Os specifies which operating system is being referred
type Os int

const (
	// OsLinux represents the Linux OS
	OsLinux Os = iota

	// OsDarwin represents the Darwin OS
	OsDarwin

	// OsUnknown represents an unknown OS
	OsUnknown
)
