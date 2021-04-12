package storm

import (
	"errors"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/helm"
	"helm.sh/helm/v3/pkg/release"
)

type helmState struct {
	node stormpkg.Node
}

// Helm contains storm compatible state
type Helm struct {
	Metadata `storm:"inline"`

	ReleaseName string `storm:"unique"`
	ID          ID
	Release     *release.Release
	Chart       *helm.Chart
}

// NewHelm returns storm compatible state
func NewHelm(h *client.Helm, meta Metadata) *Helm {
	return &Helm{
		Metadata:    meta,
		ReleaseName: h.Chart.ReleaseName,
		ID:          NewID(h.ID),
		Release:     h.Release,
		Chart:       h.Chart,
	}
}

// Convert to client.Helm
func (h *Helm) Convert() *client.Helm {
	return &client.Helm{
		ID:      h.ID.Convert(),
		Release: h.Release,
		Chart:   h.Chart,
	}
}

func (h *helmState) SaveHelmRelease(helm *client.Helm) error {
	return h.node.Save(NewHelm(helm, NewMetadata()))
}

func (h *helmState) RemoveHelmRelease(releaseName string) error {
	r := &Helm{}

	err := h.node.One("ReleaseName", releaseName, r)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return h.node.DeleteStruct(r)
}

func (h *helmState) GetHelmRelease(releaseName string) (*client.Helm, error) {
	r := &Helm{}

	err := h.node.One("ReleaseName", releaseName, r)
	if err != nil {
		return nil, err
	}

	return r.Convert(), nil
}

// NewHelmState returns an initialised helm state
func NewHelmState(node stormpkg.Node) client.HelmState {
	return &helmState{
		node: node,
	}
}
