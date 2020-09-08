package github_test

import (
	"context"
	"net/http"
	"testing"

	"gopkg.in/h2non/gock.v1"

	ghPkg "github.com/google/go-github/v32/github"
	githubAuth "github.com/oslokommune/okctl/pkg/credentials/github"
	"github.com/oslokommune/okctl/pkg/github"
	"github.com/stretchr/testify/assert"
)

func TestGithubRepositories(t *testing.T) {
	repositories := []*ghPkg.Repository{
		{
			ID:      github.Int64Ptr(12345),
			Name:    github.StringPtr("something"),
			Private: github.BoolPtr(true),
		},
	}

	testCases := []struct {
		name   string
		github *github.Github
		expect []*github.Repository
	}{
		{
			name: "Should work",
			github: func() *github.Github {
				gh, err := github.New(context.Background(), githubAuth.New(githubAuth.NewInMemoryPersister(), &http.Client{}, githubAuth.NewAuthStatic(&githubAuth.Credentials{
					AccessToken: "meh",
					Type:        githubAuth.CredentialsTypePersonalAccessToken,
				})))
				assert.NoError(t, err)

				return gh
			}(),
			expect: repositories,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			gock.New(tc.github.Client.BaseURL.String()).
				Get("/orgs/oslokommune/repos").
				MatchParam("per_page", "10").
				Reply(http.StatusOK).
				JSON(repositories)

			got, err := tc.github.Repositories(github.DefaultOrg)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, got)
		})
	}
}

func TestGithubTeams(t *testing.T) {
	teams := []*ghPkg.Team{
		{
			ID:      github.Int64Ptr(12345),
			Name:    github.StringPtr("myTeam"),
			Privacy: github.StringPtr("secret"),
		},
	}

	testCases := []struct {
		name   string
		github *github.Github
		expect []*github.Team
	}{
		{
			name: "Should work",
			github: func() *github.Github {
				gh, err := github.New(context.Background(), githubAuth.New(githubAuth.NewInMemoryPersister(), &http.Client{}, githubAuth.NewAuthStatic(&githubAuth.Credentials{
					AccessToken: "meh",
					Type:        githubAuth.CredentialsTypePersonalAccessToken,
				})))
				assert.NoError(t, err)

				return gh
			}(),
			expect: teams,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			gock.New(tc.github.Client.BaseURL.String()).
				Get("/orgs/oslokommune/teams").
				MatchParam("per_page", "10").
				Reply(http.StatusOK).
				JSON(teams)

			got, err := tc.github.Teams(github.DefaultOrg)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, got)
		})
	}
}

func TestGithubCreateDeployKey(t *testing.T) {
	key := &ghPkg.Key{
		ID:       github.Int64Ptr(12345),
		Key:      github.StringPtr("ssh-rsa 1234567abc"),
		URL:      github.StringPtr("https://"),
		Title:    github.StringPtr("myTitle"),
		ReadOnly: github.BoolPtr(true),
	}

	testCases := []struct {
		name   string
		github *github.Github
		expect *github.Key
	}{
		{
			name: "Should work",
			github: func() *github.Github {
				gh, err := github.New(context.Background(), githubAuth.New(githubAuth.NewInMemoryPersister(), &http.Client{}, githubAuth.NewAuthStatic(&githubAuth.Credentials{
					AccessToken: "meh",
					Type:        githubAuth.CredentialsTypePersonalAccessToken,
				})))
				assert.NoError(t, err)

				return gh
			}(),
			expect: key,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			defer gock.Off()

			gock.New(tc.github.Client.BaseURL.String()).
				Post("/repos/oslokommune/myRepo/keys").
				Reply(http.StatusOK).
				JSON(key)

			got, err := tc.github.CreateDeployKey(github.DefaultOrg, "myRepo", "myTitle", "ssh-rsa 1234567abc")
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, got)
		})
	}
}
