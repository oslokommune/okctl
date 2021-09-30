package version

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"sort"
)

type Versioner struct {
	github *githubb
}

// Info contains the version information
type Info struct {
	Version     string
	ShortCommit string
	BuildDate   string
}

// GetVersionInfo populates the version information
func (v Versioner) GetVersionInfo() (Info, error) {
	var semanticVersion string

	if Version == devVersion {
		// Version needs to be a valid semantic version, so we need to replace it with something else
		v, err := v.fetchSemanticDevVersion()
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
	}, nil
}

// String returns version info as JSON
func (v Versioner) String() (string, error) {
	versionInfo, err := v.GetVersionInfo()
	if err != nil {
		return "", fmt.Errorf("getting version info: %w", err)
	}

	data, err := json.Marshal(versionInfo)
	if err != nil {
		return "", fmt.Errorf("marshalling version info json: %w", err)
	}

	return string(data), nil
}

func (v Versioner) fetchSemanticDevVersion() (*semver.Version, error) {
	releases, err := v.github.ListReleases("oslokommune", "okctl")
	if err != nil {
		return nil, fmt.Errorf("listing releases: %w", err)
	}

	sort.SliceStable(releases, func(i, j int) bool {
		iVersion, err := semver.NewVersion(releases[i].GetTagName())
		if err != nil {
			return false
		}

		jVersion, err := semver.NewVersion(releases[j].GetTagName())
		if err != nil {
			return false
		}

		return iVersion.LessThan(jVersion)
	})

	newestVersionString := releases[len(releases)-1].GetTagName()
	newestVersion, err := semver.NewVersion(newestVersionString)
	if err != nil {
		return nil, fmt.Errorf("parsing version string '%s': %w", newestVersionString, err)
	}

	return newestVersion, nil
}

func New(ctx context.Context) Versioner {
	return Versioner{
		github: newGithub(ctx),
	}
}
