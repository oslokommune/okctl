package upgrade

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/github"
)

type githubServiceMock struct {
	releases []*github.RepositoryRelease
}

func (g githubServiceMock) CreateGithubRepository(ctx context.Context, opts client.CreateGithubRepositoryOpts) (*client.GithubRepository, error) {
	panic("implement me")
}

func (g githubServiceMock) DeleteGithubRepository(ctx context.Context, opts client.DeleteGithubRepositoryOpts) error {
	panic("implement me")
}

func (g githubServiceMock) ListReleases(owner, repo string) ([]*github.RepositoryRelease, error) {
	panic("implement me")
}

// NewGithubServiceMock returns a new client.GithubService
func NewGithubServiceMock(releases []*github.RepositoryRelease) client.GithubService {
	return githubServiceMock{
		releases: releases,
	}
}
