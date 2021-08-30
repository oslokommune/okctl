package storm

import (
	"errors"
	"fmt"
	"time"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/breeze"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type argoCDState struct {
	node breeze.Client
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
	existing, err := a.getArgoCD()
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return a.node.Save(NewArgoCD(cd, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return a.node.Save(NewArgoCD(cd, existing.Metadata))
}

func (a *argoCDState) getArgoCD() (*ArgoCD, error) {
	cd := &ArgoCD{}

	err := a.node.One("Name", "argocd", cd)
	if err != nil {
		return nil, err
	}

	return cd, nil
}

func (a *argoCDState) GetArgoCD() (*client.ArgoCD, error) {
	cd, err := a.getArgoCD()
	if err != nil {
		return nil, err
	}

	return cd.Convert(), nil
}

func (a *argoCDState) HasArgoCD() (bool, error) {
	_, err := a.GetArgoCD()
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf(constant.QueryStateError, err)
	}

	return true, nil
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
func NewArgoCDState(node breeze.Client) client.ArgoCDState {
	return &argoCDState{
		node: node,
	}
}
