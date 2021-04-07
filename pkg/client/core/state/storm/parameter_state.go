package storm

import (
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
	Name    string `storm:"unique,index"`
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
	return p.node.Save(NewSecretParameter(parameter, NewMetadata()))
}

func (p *parameterState) GetSecret(name string) (*client.SecretParameter, error) {
	s := &SecretParameter{}

	err := p.node.One("Name", name, s)
	if err != nil {
		return nil, err
	}

	return s.Convert(), nil
}

func (p *parameterState) RemoveSecret(name string) error {
	s := &SecretParameter{}

	err := p.node.One("Name", name, s)
	if err != nil {
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
