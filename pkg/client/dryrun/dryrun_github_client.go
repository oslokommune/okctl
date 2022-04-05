package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/github"
	"io"
)

type githubService struct {
	out io.Writer
}

func (g githubService) CreateGithubRepository(_ context.Context, _ client.CreateGithubRepositoryOpts) (*client.GithubRepository, error) {
	return &client.GithubRepository{}, nil
}

func (g githubService) DeleteGithubRepository(_ context.Context, _ client.DeleteGithubRepositoryOpts) error {
	return nil
}

func (g githubService) CreateRepositoryDeployKey(_ client.CreateGithubDeployKeyOpts) (*client.GithubDeployKey, error) {
	fmt.Fprintf(g.out, formatCreate("Github Deploy Key"))

	return &client.GithubDeployKey{}, nil
}

func (g githubService) DeleteRepositoryDeployKey(_ client.DeleteGithubDeployKeyOpts) error {
	fmt.Fprintf(g.out, formatDelete("Github Deploy Key"))

	return nil
}

func (g githubService) ListReleases(_, _ string) ([]*github.RepositoryRelease, error) {
	panic("implement me")
}
