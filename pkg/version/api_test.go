package version_test

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/version"
	"github.com/oslokommune/okctl/pkg/version/developmentversion"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

//nolint:funlen
func TestGetVersionInfoInDevelopment(t *testing.T) {
	testCases := []struct {
		name               string
		withGithubReleases []*developmentversion.RepositoryRelease
		expectVersion      string
		withError          bool
	}{
		{
			name: "Should get github release version",
			withGithubReleases: []*developmentversion.RepositoryRelease{
				{
					TagName: github.StringPtr("0.0.50"),
				},
			},
			expectVersion: "0.0.50",
		},
		{
			name: "Should get version from latest github release",
			withGithubReleases: []*developmentversion.RepositoryRelease{
				{
					TagName: github.StringPtr("0.0.10"),
				},
				{
					TagName: github.StringPtr("0.0.40"),
				},
				{
					TagName: github.StringPtr("0.0.30"),
				},
			},
			expectVersion: "0.0.40",
		},
		{
			name:          "Should return hard coded version if github fails",
			withError:     true,
			expectVersion: developmentversion.HardCodedVersion,
		},
		{
			name:          "Should return hard coded version if there are no releases",
			expectVersion: developmentversion.HardCodedVersion,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			var listReleases developmentversion.ListReleasesFn

			if !tc.withError {
				listReleases = func(ctx context.Context, owner string, repo string) ([]*developmentversion.RepositoryRelease, error) {
					return tc.withGithubReleases, nil
				}
			} else {
				listReleases = func(ctx context.Context, owner string, repo string) ([]*developmentversion.RepositoryRelease, error) {
					return nil, fmt.Errorf("some error")
				}
			}
			defer developmentversion.Reset()

			developmentversion.ListReleases = listReleases

			// When
			v := version.GetVersionInfo().Version

			// Then
			assert.Equal(t, tc.expectVersion, v)
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		name               string
		withGithubReleases []*developmentversion.RepositoryRelease
	}{
		{
			name: "Should get JSON marshalled version info",
			withGithubReleases: []*developmentversion.RepositoryRelease{
				{
					TagName: github.StringPtr("0.0.20"),
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			listReleases := func(ctx context.Context, owner string, repo string) ([]*developmentversion.RepositoryRelease, error) {
				return tc.withGithubReleases, nil
			}
			developmentversion.ListReleases = listReleases

			defer developmentversion.Reset()

			// When
			v := version.String()

			// Then
			g := goldie.New(t)
			g.Assert(t, tc.name, []byte(v))
		})
	}
}
