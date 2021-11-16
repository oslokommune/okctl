package version

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/oslokommune/okctl/pkg/version/github"
)

// GetVersionInfo returns the current version
func GetVersionInfo() Info {
	var semanticVersion string

	if Version == devVersion {
		v, err := github.FetchVersion(context.Background())
		if err != nil {
			semanticVersion = "0.0.10"
			fmt.Printf("Warning: Could not get version, using hard coded version '%s' instead\n", semanticVersion)
		} else {
			semanticVersion = v.String()
		}
	} else {
		semanticVersion = Version
	}

	return Info{
		Version:     semanticVersion,
		ShortCommit: ShortCommit,
		BuildDate:   BuildDate,
	}
}

// String returns the current version
func String() string {
	versionInfo := GetVersionInfo()

	data, err := json.Marshal(versionInfo)
	if err != nil {
		panic(fmt.Sprintf("Could not marshal versionInfo: %s", err.Error()))
	}

	return string(data)
}
