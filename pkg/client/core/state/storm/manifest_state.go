package storm

import (
	"errors"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type manifestState struct {
	node stormpkg.Node
}

// KubernetesManifest contains storm compatible state
type KubernetesManifest struct {
	Metadata `storm:"inline"`

	ID        ID
	Name      string `storm:"unique"`
	Namespace string
	Type      string
	Content   []byte
}

// NewKubernetesManifest constructs a storm compatible KubernetesManifest
func NewKubernetesManifest(m *client.KubernetesManifest, meta Metadata) *KubernetesManifest {
	return &KubernetesManifest{
		Metadata:  meta,
		ID:        NewID(m.ID),
		Name:      m.Name,
		Namespace: m.Namespace,
		Type:      m.Type.String(),
		Content:   m.Content,
	}
}

// Convert to a client.KubernetesManifest
func (m *KubernetesManifest) Convert() *client.KubernetesManifest {
	return &client.KubernetesManifest{
		ID:        m.ID.Convert(),
		Name:      m.Name,
		Namespace: m.Namespace,
		Type:      client.ManifestType(m.Type),
		Content:   m.Content,
	}
}

func (s *manifestState) SaveKubernetesManifests(manifest *client.KubernetesManifest) error {
	return s.node.Save(NewKubernetesManifest(manifest, NewMetadata()))
}

func (s *manifestState) GetKubernetesManifests(name string) (*client.KubernetesManifest, error) {
	m := &KubernetesManifest{}

	err := s.node.One("Name", name, m)
	if err != nil {
		return nil, err
	}

	return m.Convert(), nil
}

func (s *manifestState) RemoveKubernetesManifests(name string) error {
	m := &KubernetesManifest{}

	err := s.node.One("Name", name, m)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return s.node.DeleteStruct(m)
}

// NewManifestState returns an initialised manifest state
func NewManifestState(node stormpkg.Node) client.ManifestState {
	return &manifestState{
		node: node,
	}
}
