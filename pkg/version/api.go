package version

import (
	"encoding/json"
	"fmt"
	"github.com/oslokommune/okctl/pkg/version/developmentversion"
	"github.com/oslokommune/okctl/pkg/version/types"
)

// GetVersionInfo returns the current version
func GetVersionInfo() types.Info {
	var semanticVersion string

	if Version == DevVersion {
		semanticVersion = developmentversion.GetVersionInfo()
	} else {
		semanticVersion = Version
	}

	return types.Info{
		Version:     semanticVersion,
		ShortCommit: ShortCommit,
		BuildDate:   BuildDate,
	}
}

// String returns the current version as JSON
func String() string {
	versionInfo := GetVersionInfo()

	data, err := json.Marshal(versionInfo)
	if err != nil {
		panic(fmt.Sprintf("Could not marshal versionInfo: %s", err.Error()))
	}

	return string(data)
}
