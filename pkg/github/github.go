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

// DefaultAWSInfrastructureRepository is the name of the infrastructure as code repository which should be given
// nameserver delegation requests
const DefaultAWSInfrastructureRepository = "origo-aws-infrastructure"

// DefaultAWSInfrastructurePrimaryBranch is the name of the primary branch (due to git/github moving away from "master"
// to "main")
const DefaultAWSInfrastructurePrimaryBranch = "master"

// Githuber invokes the github API
type Githuber interface {
	Teams(org string) ([]*Team, error)
	Repositories(org string) ([]*Repository, error)
	CreateDeployKey(org, repository, title, publicKey string) (*Key, error)
	DeleteDeployKey(org, repository string, id int64) error
	CreatePullRequest(r *PullRequest) error
	ListReleases(owner, repo string) ([]*RepositoryRelease, error)
}

// Github contains the state for interacting with the github API
type Github struct {
	Ctx    context.Context
	Client *github.Client
}

// Ensure that Github implements Githuber
var _ Githuber = &Github{}

// Repository shadows github.Repository
type Repository = github.Repository

// Team shadows github.Team
type Team = github.Team

// Key shadows github.Key
type Key = github.Key

type RepositoryRelease = github.RepositoryRelease

type ReleaseAsset = github.ReleaseAsset

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
		return nil, fmt.Errorf("creating deploy key: %w", err)
	}

	return key, nil
}

// DeleteDeployKey removes a read-only deploy key
func (g *Github) DeleteDeployKey(org, repository string, identifier int64) error {
	_, err := g.Client.Repositories.DeleteKey(g.Ctx, org, repository, identifier)
	if err != nil {
		return fmt.Errorf("deleting deploy key: %w", err)
	}

	return nil
}

// PullRequest contains data about the PR
type PullRequest struct {
	Organisation      string
	Repository        string
	SourceBranch      string
	DestinationBranch string
	Title             string
	Body              string
	Labels            []string
}

// CreatePullRequest creates a pull request from sourceBranch to destinationBranch
func (g *Github) CreatePullRequest(r *PullRequest) error {
	pr, _, err := g.Client.PullRequests.Create(g.Ctx, r.Organisation, r.Repository, &github.NewPullRequest{
		Title: StringPtr(r.Title),
		Head:  StringPtr(r.SourceBranch),
		Base:  StringPtr(r.DestinationBranch),
		Body:  StringPtr(r.Body),
	})
	if err != nil {
		return fmt.Errorf("creating github pull request: %w", err)
	}

	if len(r.Labels) > 0 {
		_, _, err = g.Client.Issues.AddLabelsToIssue(g.Ctx, r.Organisation, r.Repository, pr.GetNumber(), r.Labels)
		if err != nil {
			return fmt.Errorf("adding labels to pull request: %w", err)
		}
	}

	return nil
}

const ListReleasesPageSize = 100

func (g *Github) ListReleases(owner, repo string) ([]*RepositoryRelease, error) {
	opts := &github.ListOptions{
		PerPage: ListReleasesPageSize,
	}

	var allReleases []*RepositoryRelease

	for {
		// Documentation: https://docs.github.com/en/rest/reference/repos#list-release-assets
		releases, response, err := g.Client.Repositories.ListReleases(g.Ctx, owner, repo, opts)
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
