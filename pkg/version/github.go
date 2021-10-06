package version

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
)

// Github knows how to do GitHub operations
type Github interface {
	ListReleases(ctx context.Context, owner, repo string) ([]*RepositoryRelease, error)
}

type httpGithuber struct {
	Client *github.Client
}

// RepositoryRelease shadows github.RepositoryRelease
type RepositoryRelease = github.RepositoryRelease

const listReleasesPageSize = 100

// ListReleases lists the given repository's releases
func (g *httpGithuber) ListReleases(ctx context.Context, owner, repo string) ([]*RepositoryRelease, error) {
	opts := &github.ListOptions{
		PerPage: listReleasesPageSize,
	}

	var allReleases []*RepositoryRelease

	for {
		// Documentation: https://docs.github.com/en/rest/reference/repos#list-release-assets
		releases, response, err := g.Client.Repositories.ListReleases(ctx, owner, repo, opts)
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

func newGithub() Github {
	return &httpGithuber{
		Client: github.NewClient(nil),
	}
}
