package client

import (
	"context"

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
func (r *GithubRepository) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.DeployKey, validation.Required),
		validation.Field(&r.Organisation, validation.Required),
		validation.Field(&r.FullName, validation.Required),
		validation.Field(&r.GitURL, validation.Required),
		validation.Field(&r.Repository, validation.Required),
	)
}

// ReadyGithubInfrastructureRepositoryOpts contains required inputs
type ReadyGithubInfrastructureRepositoryOpts struct {
	ID           api.ID
	Organisation string
	Repository   string // +optional
}

// Validate the inputs
func (o ReadyGithubInfrastructureRepositoryOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Organisation, validation.Required),
	)
}

// SelectedGithubRepository contains the selected github repo
type SelectedGithubRepository struct {
	ID           api.ID
	Organisation string
	Repository   string
	FullName     string
	GitURL       string
}

// SelectGithubInfrastructureRepositoryOpts contains required inputs
type SelectGithubInfrastructureRepositoryOpts struct {
	ID           api.ID
	Organisation string
	Repository   string
}

// GithubSecret represents an SSM secret parameter
type GithubSecret struct {
	Name    string
	Path    string
	Version int64
}

// Validate the data
func (s *GithubSecret) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Path, validation.Required),
	)
}

// GithubOauthApp is a github oauth app
type GithubOauthApp struct {
	ID           api.ID
	Organisation string
	Name         string
	SiteURL      string
	CallbackURL  string
	ClientID     string
	ClientSecret *GithubSecret
	Team         *GithubTeam
}

// Validate the data
func (a *GithubOauthApp) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Organisation, validation.Required),
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.SiteURL, validation.Required),
		validation.Field(&a.CallbackURL, validation.Required),
		validation.Field(&a.ClientID, validation.Required),
		validation.Field(&a.ClientSecret, validation.Required),
		validation.Field(&a.Team, validation.Required),
	)
}

// CreateGithubOauthAppOpts contains required inputs
type CreateGithubOauthAppOpts struct {
	ID           api.ID
	Organisation string
	Team         *GithubTeam // +optional
	Name         string
	SiteURL      string
	CallbackURL  string
}

// Validate the inputs
func (o CreateGithubOauthAppOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Organisation, validation.Required),
		validation.Field(&o.Name, validation.Required),
		validation.Field(&o.SiteURL, validation.Required),
		validation.Field(&o.CallbackURL, validation.Required),
	)
}

// GithubTeam is a github team
type GithubTeam struct {
	ID           api.ID
	Organisation string
	Name         string
}

// Validate the data
func (t *GithubTeam) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Name, validation.Required),
	)
}

// SelectGithubTeam contains required inputs
type SelectGithubTeam struct {
	ID           api.ID
	Organisation string
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
func (k *GithubDeployKey) Validate() error {
	return validation.ValidateStruct(k,
		validation.Field(&k.Organisation, validation.Required),
		validation.Field(&k.Repository, validation.Required),
		validation.Field(&k.Identifier, validation.Required),
		validation.Field(&k.Title, validation.Required),
		validation.Field(&k.PublicKey, validation.Required),
		validation.Field(&k.PrivateKeySecret, validation.Required),
	)
}

// CreateGithubDeployKey contains required inputs
type CreateGithubDeployKey struct {
	ID           api.ID
	Organisation string
	Repository   string
	Title        string
}

// GithubService is a business logic implementation
type GithubService interface {
	ReadyGithubInfrastructureRepository(ctx context.Context, opts ReadyGithubInfrastructureRepositoryOpts) (*GithubRepository, error)
	CreateGithubOauthApp(ctx context.Context, opts CreateGithubOauthAppOpts) (*GithubOauthApp, error)
}

// GithubAPI invokes the Github API
type GithubAPI interface {
	SelectGithubInfrastructureRepository(opts SelectGithubInfrastructureRepositoryOpts) (*SelectedGithubRepository, error)
	CreateGithubDeployKey(opts CreateGithubDeployKey) (*GithubDeployKey, error)
	SelectGithubTeam(opts SelectGithubTeam) (*GithubTeam, error)
	CreateGithubOauthApp(opts CreateGithubOauthAppOpts) (*GithubOauthApp, error)
}

// GithubReport is the report layer
type GithubReport interface {
	ReadyGithubInfrastructureRepository(repository *GithubRepository, report *store.Report) error
	CreateGithubOauthApp(app *GithubOauthApp, report *store.Report) error
}

// GithubState is the state layer
type GithubState interface {
	SaveGithubInfrastructureRepository(repository *GithubRepository) (*store.Report, error)
	GetGithubInfrastructureRepository(id api.ID) *GithubRepository
	SaveGithubOauthApp(app *GithubOauthApp) (*store.Report, error)
	GetGithubOauthApp(appName string, id api.ID) *GithubOauthApp
}
