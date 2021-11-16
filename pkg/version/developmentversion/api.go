// Package developmentversion knows how to fetch the current version in development mode
package developmentversion

import (
	"context"
	"fmt"
	sortPkg "sort"

	"github.com/Masterminds/semver"

	"github.com/google/go-github/v32/github"
)

var cachedVersion string //nolint:gochecknoglobals // Ignoring global since version.GetVersionInfo() is global already

// GetVersionInfo returns the current version
func GetVersionInfo() string {
	if len(cachedVersion) > 0 {
		fmt.Println("Using cache")
		return cachedVersion
	}

	cachedVersion = getGithubOrHardCodedVersion()

	return cachedVersion
}

// Reset resets all globals. Useful for testing.
func Reset() {
	cachedVersion = ""
	ListReleases = listReleases
}

// HardCodedVersion is the version used if contact with GitHub fails
const HardCodedVersion = "0.0.10"

func getGithubOrHardCodedVersion() string {
	ver, err := doFetchVersion()
	if err != nil {
		fmt.Printf("Warning: Could not get version, using hard coded version '%s' instead\n", HardCodedVersion)

		return HardCodedVersion
	}

	return ver.String()
}

// ListReleases knows how to list GitHub releases
var ListReleases ListReleasesFn = listReleases //nolint:gochecknoglobals

func doFetchVersion() (*semver.Version, error) {
	releases, err := ListReleases(context.Background(), "oslokommune", "okctl")
	if err != nil {
		return nil, fmt.Errorf("listing releases: %w", err)
	}

	sort(releases)

	if len(releases) == 0 {
		return nil, fmt.Errorf("expected one or more releases, but was zero")
	}

	newestVersionString := releases[len(releases)-1].GetTagName()

	newestVersion, err := semver.NewVersion(newestVersionString)
	if err != nil {
		return nil, fmt.Errorf("parsing version string '%s': %w", newestVersionString, err)
	}

	return newestVersion, nil
}

func sort(releases []*RepositoryRelease) {
	sortPkg.SliceStable(releases, func(i, j int) bool {
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
}

const listReleasesPageSize = 100

func listReleases(ctx context.Context, owner, repo string) ([]*RepositoryRelease, error) {
	client := github.NewClient(nil)

	opts := &github.ListOptions{
		PerPage: listReleasesPageSize,
	}

	var allReleases []*RepositoryRelease

	for {
		// Documentation: https://docs.github.com/en/rest/reference/repos#list-release-assets
		releases, response, err := client.Repositories.ListReleases(ctx, owner, repo, opts)
		if err != nil {
			return nil, fmt.Errorf("listing releases: %w", err)
		}

		allReleases = append(allReleases, releases...)

		if response.NextPage == 0 {
			break
		}

		opts.Page = response.NextPage
	}

	return allReleases, nil
}
