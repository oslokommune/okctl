// Package github provides a client for interacting with the Github API
package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
	githubAuth "github.com/oslokommune/okctl/pkg/credentials/github"
	"golang.org/x/oauth2"
)

// DefaultOrg is the default organisation used with okctl
const DefaultOrg = "oslokommune"

// Github contains the state for interacting with the github API
type Github struct {
	Organisation string
	Ctx          context.Context
	Client       *github.Client
}

// Repository shadows github.Repository
type Repository = github.Repository

// Team shadows github.Team
type Team = github.Team

// Key shadows github.Key
type Key = github.Key

// New returns an initialised github API client
func New(org string, auth githubAuth.Authenticator) (*Github, error) {
	ctx := context.Background()

	credentials, err := auth.Raw()
	if err != nil {
		return nil, fmt.Errorf("failed to get github credentials: %w", err)
	}

	client := github.NewClient(
		oauth2.NewClient(ctx,
			oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: credentials.AccessToken,
				},
			),
		),
	)

	return &Github{
		Organisation: org,
		Ctx:          ctx,
		Client:       client,
	}, nil
}

// Teams fetches all teams within the given organisation
func (g *Github) Teams() ([]*Team, error) {
	opts := &github.ListOptions{
		PerPage: 10, // nolint: gomnd
	}

	var allTeams []*github.Team

	for {
		teams, resp, err := g.Client.Teams.ListTeams(g.Ctx, g.Organisation, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve teams: %w", err)
		}

		allTeams = append(allTeams, teams...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return allTeams, nil
}

// Repositories fetches all the repositories within the given organisation
func (g *Github) Repositories() ([]*Repository, error) {
	opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 10, // nolint: gomnd
		},
	}

	var allRepos []*github.Repository

	for {
		repos, resp, err := g.Client.Repositories.ListByOrg(g.Ctx, g.Organisation, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve repositories: %w", err)
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return allRepos, nil
}

// CreateDeployKey creates a read-only deploy key for the given owner/repo
func (g *Github) CreateDeployKey(repository, title, publicKey string) (*Key, error) {
	key, _, err := g.Client.Repositories.CreateKey(g.Ctx, g.Organisation, repository, &github.Key{
		Title:    &title,
		Key:      &publicKey,
		ReadOnly: BoolPtr(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create deploy key: %w", err)
	}

	return key, nil
}

// BoolPtr returns a pointer to the bool
func BoolPtr(v bool) *bool {
	return &v
}

// StringPtr returns a pointer to the string
func StringPtr(v string) *string {
	return &v
}

// Int64Ptr returns a pointer to the int64
func Int64Ptr(v int64) *int64 {
	return &v
}
