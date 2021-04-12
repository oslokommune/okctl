package storm

import (
	"errors"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type argoCDState struct {
	node stormpkg.Node
}

// ArgoCD contains state about an argo cd deployment
type ArgoCD struct {
	Metadata `storm:"inline"`
	Name     string `storm:"unique"`

	ID         ID
	ArgoDomain string
	ArgoURL    string
	AuthDomain string
}

// NewArgoCD returns storm compatible state
func NewArgoCD(a *client.ArgoCD, meta Metadata) *ArgoCD {
	return &ArgoCD{
		Metadata:   meta,
		Name:       "argocd",
		ID:         NewID(a.ID),
		ArgoDomain: a.ArgoDomain,
		ArgoURL:    a.ArgoURL,
		AuthDomain: a.AuthDomain,
	}
}

// Convert to client.ArgoCD
func (a *ArgoCD) Convert() *client.ArgoCD {
	return &client.ArgoCD{
		ID:         a.ID.Convert(),
		ArgoDomain: a.ArgoDomain,
		ArgoURL:    a.ArgoURL,
		AuthDomain: a.AuthDomain,
	}
}

func (a *argoCDState) SaveArgoCD(cd *client.ArgoCD) error {
	return a.node.Save(NewArgoCD(cd, NewMetadata()))
}

func (a *argoCDState) GetArgoCD() (*client.ArgoCD, error) {
	cd := &ArgoCD{}

	err := a.node.One("Name", "argocd", cd)
	if err != nil {
		return nil, err
	}

	return cd.Convert(), nil
}

func (a *argoCDState) RemoveArgoCD() error {
	cd := &ArgoCD{}

	err := a.node.One("Name", "argocd", cd)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return a.node.DeleteStruct(cd)
}

// NewArgoCDState returns an initialised state client
func NewArgoCDState(node stormpkg.Node) client.ArgoCDState {
	return &argoCDState{
		node: node,
	}
}
