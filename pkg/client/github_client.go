package client

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/state"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
)

// GithubRepository is a github repository
type GithubRepository struct {
	ID           api.ID
	Organisation string
	Repository   string
	FullName     string
	GitURL       string
	DeployKey    *GithubDeployKey
}

// Validate the github repository
func (r GithubRepository) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.DeployKey, validation.Required),
		validation.Field(&r.Organisation, validation.Required),
		validation.Field(&r.FullName, validation.Required),
		validation.Field(&r.GitURL, validation.Required),
		validation.Field(&r.Repository, validation.Required),
	)
}

// NewGithubRepository initializes a new GithubRepository
func NewGithubRepository(clusterID api.ID, host, organization, name string) *GithubRepository {
	fullName := fmt.Sprintf("%s/%s", organization, name)

	return &GithubRepository{
		ID:           clusterID,
		Organisation: organization,
		Repository:   name,
		FullName:     fullName,
		GitURL:       fmt.Sprintf("%s:%s", host, fullName),
		DeployKey:    nil,
	}
}

// GithubSecret represents an SSM secret parameter
type GithubSecret struct {
	Name    string
	Path    string
	Version int64
}

// Validate the data
func (s GithubSecret) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Path, validation.Required),
	)
}

// GithubDeployKey is a github deploy key
type GithubDeployKey struct {
	ID               api.ID
	Organisation     string
	Repository       string
	Identifier       int64
	Title            string
	PublicKey        string
	PrivateKeySecret *GithubSecret
}

// Validate the data
func (k GithubDeployKey) Validate() error {
	return validation.ValidateStruct(&k,
		validation.Field(&k.Organisation, validation.Required),
		validation.Field(&k.Repository, validation.Required),
		validation.Field(&k.Identifier, validation.Required),
		validation.Field(&k.Title, validation.Required),
		validation.Field(&k.PublicKey, validation.Required),
		validation.Field(&k.PrivateKeySecret, validation.Required),
	)
}

// CreateGithubDeployKeyOpts contains required inputs
type CreateGithubDeployKeyOpts struct {
	ID           api.ID
	Organisation string
	Repository   string
	Title        string
}

// GithubService is a business logic implementation
type GithubService interface {
	CreateRepositoryDeployKey(ctx context.Context, repository *GithubRepository) (*GithubDeployKey, error)
}

// GithubAPI invokes the Github API
type GithubAPI interface {
	CreateRepositoryDeployKey(opts CreateGithubDeployKeyOpts) (*GithubDeployKey, error)
}

// GithubState is the state layer
type GithubState interface {
	SaveRepositoryDeployKey(repository *GithubRepository) (*store.Report, error)
	GetRepositoryDeployKey(id api.ID) state.GithubRepository
}
