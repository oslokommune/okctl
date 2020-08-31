// Package github provides a client for interacting with the Github API
package github

import "github.com/oslokommune/okctl/pkg/credentials/github"

// Github contains the state for interacting with the github API
type Github struct {
	auth github.Authenticator
}

// New returns an initialised github API client
func New(auth github.Authenticator) *Github {
	return &Github{
		auth: auth,
	}
}

// Teams fetches all teams within the oslokommune org
func (g *Github) Teams() error {
	return nil
}

// Repositories fetches all the repositories
func (g *Github) Repositories() error {
	return nil
}

// DeployKey ...
func (g *Github) DeployKey() error {
	return nil
}

// RemoteIsThis ...
func (g *Github) RemoteIsThis(repo string) error {
	return nil
}
