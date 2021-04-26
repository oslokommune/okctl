package storm

import (
	"errors"
	"time"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type githubState struct {
	node stormpkg.Node
}

// GithubRepository contains storm compatible state
type GithubRepository struct {
	Metadata `storm:"inline"`

	ID           ID
	Organisation string
	Repository   string
	FullName     string `storm:"unique"`
	GitURL       string
	DeployKey    *GithubDeployKey
}

// NewGithubRepository returns storm compatible state
func NewGithubRepository(r *client.GithubRepository, meta Metadata) *GithubRepository {
	return &GithubRepository{
		Metadata:     meta,
		ID:           NewID(r.ID),
		Organisation: r.Organisation,
		Repository:   r.Repository,
		FullName:     r.FullName,
		GitURL:       r.GitURL,
		DeployKey:    NewGithubDeployKey(r.DeployKey),
	}
}

// Convert to client.GithubRepository
func (r *GithubRepository) Convert() *client.GithubRepository {
	return &client.GithubRepository{
		ID:           r.ID.Convert(),
		Organisation: r.Organisation,
		Repository:   r.Repository,
		FullName:     r.FullName,
		GitURL:       r.GitURL,
		DeployKey:    r.DeployKey.Convert(),
	}
}

// GithubDeployKey contains storm compatible state
type GithubDeployKey struct {
	Organisation     string
	Repository       string
	Identifier       int64
	Title            string
	PublicKey        string
	PrivateKeySecret *GithubSecret
}

// NewGithubDeployKey returns storm compatible state
func NewGithubDeployKey(k *client.GithubDeployKey) *GithubDeployKey {
	return &GithubDeployKey{
		Organisation:     k.Organisation,
		Repository:       k.Repository,
		Identifier:       k.Identifier,
		Title:            k.Title,
		PublicKey:        k.PublicKey,
		PrivateKeySecret: NewGithubSecret(k.PrivateKeySecret),
	}
}

// Convert to client.GithubDeployKey
func (k *GithubDeployKey) Convert() *client.GithubDeployKey {
	return &client.GithubDeployKey{
		Organisation:     k.Organisation,
		Repository:       k.Repository,
		Identifier:       k.Identifier,
		Title:            k.Title,
		PublicKey:        k.PublicKey,
		PrivateKeySecret: k.PrivateKeySecret.Convert(),
	}
}

// GithubSecret contains storm compatible state
type GithubSecret struct {
	Name    string
	Path    string
	Version int64
}

// NewGithubSecret returns storm compatible state
func NewGithubSecret(s *client.GithubSecret) *GithubSecret {
	return &GithubSecret{
		Name:    s.Name,
		Path:    s.Path,
		Version: s.Version,
	}
}

// Convert to client.GithubSecret
func (s *GithubSecret) Convert() *client.GithubSecret {
	return &client.GithubSecret{
		Name:    s.Name,
		Path:    s.Path,
		Version: s.Version,
	}
}

func (g *githubState) SaveGithubRepository(repository *client.GithubRepository) error {
	existing, err := g.getGithubRepository(repository.FullName)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return g.node.Save(NewGithubRepository(repository, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return g.node.Save(NewGithubRepository(repository, existing.Metadata))
}

func (g *githubState) RemoveGithubRepository(fullName string) error {
	r := &GithubRepository{}

	err := g.node.One("FullName", fullName, r)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return g.node.DeleteStruct(r)
}

func (g *githubState) GetGithubRepository(fullName string) (*client.GithubRepository, error) {
	r, err := g.getGithubRepository(fullName)
	if err != nil {
		return nil, err
	}

	return r.Convert(), nil
}

func (g *githubState) getGithubRepository(fullName string) (*GithubRepository, error) {
	r := &GithubRepository{}

	err := g.node.One("FullName", fullName, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// NewGithubState returns an initialised state client
func NewGithubState(node stormpkg.Node) client.GithubState {
	return &githubState{
		node: node,
	}
}
