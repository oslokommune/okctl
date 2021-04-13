package client

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

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

// GithubDeployKey is a github deploy key
type GithubDeployKey struct {
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

// DeleteGithubRepositoryOpts contains the required inputs
type DeleteGithubRepositoryOpts struct {
	ID           api.ID
	Organisation string
	Name         string
}

// CreateGithubRepositoryOpts contains the required inputs
type CreateGithubRepositoryOpts struct {
	ID           api.ID
	Host         string
	Organization string
	Name         string
}

// Validate the inputs
func (o CreateGithubRepositoryOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Host, validation.Required),
		validation.Field(&o.Organization, validation.Required),
		validation.Field(&o.Name, validation.Required),
	)
}

// CreateGithubDeployKeyOpts contains required inputs
type CreateGithubDeployKeyOpts struct {
	ID           api.ID
	Organisation string
	Repository   string
	Title        string
}

// DeleteGithubDeployKeyOpts contains the required inputs
type DeleteGithubDeployKeyOpts struct {
	ID           api.ID
	Organisation string
	Repository   string
	Identifier   int64
}

// GithubService is a business logic implementation
type GithubService interface {
	CreateGithubRepository(ctx context.Context, opts CreateGithubRepositoryOpts) (*GithubRepository, error)
	DeleteGithubRepository(ctx context.Context, opts DeleteGithubRepositoryOpts) error
}

// GithubAPI invokes the Github API
type GithubAPI interface {
	CreateRepositoryDeployKey(opts CreateGithubDeployKeyOpts) (*GithubDeployKey, error)
	DeleteRepositoryDeployKey(opts DeleteGithubDeployKeyOpts) error
}

// GithubState is the state layer
type GithubState interface {
	SaveGithubRepository(repository *GithubRepository) error
	GetGithubRepository(fullName string) (*GithubRepository, error)
	RemoveGithubRepository(fullName string) error
}
