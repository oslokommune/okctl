package upgrade

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/github"
)

type githubServiceMock struct {
	releases []*github.RepositoryRelease
}

func (g githubServiceMock) CreateGithubRepository(context.Context, client.CreateGithubRepositoryOpts) (*client.GithubRepository, error) {
	panic("not needed by mock")
}

func (g githubServiceMock) DeleteGithubRepository(context.Context, client.DeleteGithubRepositoryOpts) error {
	panic("not needed by mock")
}

//goland:noinspection GoUnusedParameter
func (g githubServiceMock) ListReleases(owner, repo string) ([]*github.RepositoryRelease, error) {
	return g.releases, nil
}

// NewGithubServiceMock returns a new client.GithubService
func NewGithubServiceMock(releases []*github.RepositoryRelease) client.GithubService {
	return githubServiceMock{
		releases: releases,
	}
}
