// Package github provides a client for interacting with the Github API
package github

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/google/go-github/v32/github"
	githubAuth "github.com/oslokommune/okctl/pkg/credentials/github"
	"golang.org/x/oauth2"
)

// DefaultOrg is the default organisation used with okctl
const DefaultOrg = "oslokommune"

// Githuber invokes the github API
type Githuber interface {
	Teams(org string) ([]*Team, error)
	Repositories(org string) ([]*Repository, error)
	CreateDeployKey(org, repository, title, publicKey string) (*Key, error)
	ListTeamMembers(team client.GithubTeam) ([]client.GithubTeamMember, error)
}

// Github contains the state for interacting with the github API
type Github struct {
	Ctx    context.Context
	Client *github.Client
}

// ListTeamMembers lists members of a github team
func (g *Github) ListTeamMembers(team client.GithubTeam) ([]client.GithubTeamMember, error) {
	ctx := g.Ctx
	teamID := team.TeamID

	org, _, err := g.Client.Organizations.Get(ctx, team.Organisation)
	if err != nil {
		return nil, err
	}

	users, _, err := g.Client.Teams.ListTeamMembersByID(ctx, *org.ID, teamID, nil)
	if err != nil {
		fmt.Println(err)
	}

	members := []client.GithubTeamMember{}

	for _, u := range users {
		// Yes we really do need to do a call for every user,
		user, _, _ := g.Client.Users.Get(ctx, *u.Login)

		email := ""
		if user.Email != nil {
			email = *user.Email
		}

		members = append(members, client.GithubTeamMember{
			Login: *user.Login,
			Name:  *user.Name,
			Email: email,
		})
	}

	return members, nil
}

// Ensure that Github implements Githuber
var _ Githuber = &Github{}

// Repository shadows github.Repository
type Repository = github.Repository

// Team shadows github.Team
type Team = github.Team

// Key shadows github.Key
type Key = github.Key

// New returns an initialised github API client
func New(ctx context.Context, auth githubAuth.Authenticator) (*Github, error) {
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
		Ctx:    ctx,
		Client: client,
	}, nil
}

// Teams fetches all teams within the given organisation
func (g *Github) Teams(org string) ([]*Team, error) {
	opts := &github.ListOptions{
		PerPage: 10, // nolint: gomnd
	}

	var allTeams []*github.Team

	for {
		teams, resp, err := g.Client.Teams.ListTeams(g.Ctx, org, opts)
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
func (g *Github) Repositories(org string) ([]*Repository, error) {
	opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 10, // nolint: gomnd
		},
	}

	var allRepos []*github.Repository

	for {
		repos, resp, err := g.Client.Repositories.ListByOrg(g.Ctx, org, opts)
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
func (g *Github) CreateDeployKey(org, repository, title, publicKey string) (*Key, error) {
	key, _, err := g.Client.Repositories.CreateKey(g.Ctx, org, repository, &github.Key{
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
