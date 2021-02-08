package version

import "encoding/json"

// Info contains the version information
type Info struct {
	Version     string
	ShortCommit string
	BuildDate   string
}

// GetVersionInfo populates the version information
func GetVersionInfo() Info {
	return Info{
		Version:     Version,
		ShortCommit: ShortCommit,
		BuildDate:   BuildDate,
	}
}

// String returns version info as JSON
func String() string {
	if data, err := json.Marshal(GetVersionInfo()); err == nil {
		return string(data)
	}

	return ""
}
