// Package github knows how to fetch a version from github
package github

import (
	"context"
	"fmt"
	"sort"

	"github.com/Masterminds/semver"

	"github.com/google/go-github/v32/github"
)

var cachedVersion *semver.Version //nolint:gochecknoglobals // Ignoring global since version.GetVersionInfo() is global already

// FetchVersion fetches the latest version from GitHub
func FetchVersion(ctx context.Context) (*semver.Version, error) {
	if cachedVersion != nil {
		fmt.Println("Using cache")
		return cachedVersion, nil
	}

	var err error
	cachedVersion, err = doFetchVersion(ctx)

	return cachedVersion, err
}

func doFetchVersion(ctx context.Context) (*semver.Version, error) {
	releases, err := listReleases(ctx, "oslokommune", "okctl")
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
