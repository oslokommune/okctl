package storm

import (
	"errors"
	"time"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type parameterState struct {
	node stormpkg.Node
}

// SecretParameter contains storm compatible state
type SecretParameter struct {
	Metadata `storm:"inline"`

	ID      ID
	Name    string `storm:"unique"`
	Path    string
	Version int64
}

// NewSecretParameter returns a storm compatible SecretParameter
func NewSecretParameter(p *client.SecretParameter, meta Metadata) *SecretParameter {
	return &SecretParameter{
		Metadata: meta,
		ID:       NewID(p.ID),
		Name:     p.Name,
		Path:     p.Path,
		Version:  p.Version,
	}
}

// Convert to a client.SecretParameter
func (p *SecretParameter) Convert() *client.SecretParameter {
	return &client.SecretParameter{
		ID:      p.ID.Convert(),
		Name:    p.Name,
		Path:    p.Path,
		Version: p.Version,
		Content: "",
	}
}

func (p *parameterState) SaveSecret(parameter *client.SecretParameter) error {
	existing, err := p.getSecret(parameter.Name)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return p.node.Save(NewSecretParameter(parameter, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return p.node.Save(NewSecretParameter(parameter, existing.Metadata))
}

func (p *parameterState) GetSecret(name string) (*client.SecretParameter, error) {
	s, err := p.getSecret(name)
	if err != nil {
		return nil, err
	}

	return s.Convert(), nil
}

func (p *parameterState) getSecret(name string) (*SecretParameter, error) {
	s := &SecretParameter{}

	err := p.node.One("Name", name, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (p *parameterState) RemoveSecret(name string) error {
	s := &SecretParameter{}

	err := p.node.One("Name", name, s)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return p.node.DeleteStruct(s)
}

// NewParameterState returns an initialised state
func NewParameterState(node stormpkg.Node) client.ParameterState {
	return &parameterState{
		node: node,
	}
}
