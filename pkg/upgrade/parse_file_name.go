package upgrade

import (
	"fmt"
	"strings"
)

type okctlUpgradeFile struct {
	filename  string
	version   string
	os        string
	arch      string
	extension string
}

const (
	expectedSubstringsInOkctlUpgradeFilename    = 4
	expectedMinimumSubstringsInArchAndExtension = 2
)

// parseOkctlUpgradeFilename converts a string like 'okctl-upgrade_0.0.63_Darwin_amd64.tar.gz' to an okctlUpgradeFile
// struct.
func parseOkctlUpgradeFilename(filename string) (okctlUpgradeFile, error) {
	filenameParts := strings.Split(filename, "_")
	if len(filenameParts) != expectedSubstringsInOkctlUpgradeFilename {
		return okctlUpgradeFile{}, fmt.Errorf(
			"expected 4 substrings when splitting on underscore (_), got %d in string '%s'",
			len(filenameParts), filename,
		)
	}

	version := filenameParts[1]          // 0.0.63
	os := filenameParts[2]               // Darwin
	archAndExtension := filenameParts[3] // amd64.tar.gz

	archAndExtensionParts := strings.Split(archAndExtension, ".")
	if len(archAndExtensionParts) < expectedMinimumSubstringsInArchAndExtension {
		return okctlUpgradeFile{}, fmt.Errorf(
			"expected at least %d substrings when splitting on dot (.), got %d in string '%s'",
			expectedMinimumSubstringsInArchAndExtension, len(archAndExtensionParts), filename,
		)
	}

	arch := archAndExtensionParts[0]            // amd64
	extension := archAndExtension[len(arch)+1:] // tar.gz

	return okctlUpgradeFile{
		filename:  filename,
		version:   version,
		os:        os,
		arch:      arch,
		extension: extension,
	}, nil
}
