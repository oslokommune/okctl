package storm

import (
	"errors"
	"time"

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
	existing, err := h.getHelmRelease(helm.Chart.ReleaseName)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return h.node.Save(NewHelm(helm, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return h.node.Save(NewHelm(helm, existing.Metadata))
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
	r, err := h.getHelmRelease(releaseName)
	if err != nil {
		return nil, err
	}

	return r.Convert(), nil
}

func (h *helmState) getHelmRelease(releaseName string) (*Helm, error) {
	r := &Helm{}

	err := h.node.One("ReleaseName", releaseName, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// NewHelmState returns an initialised helm state
func NewHelmState(node stormpkg.Node) client.HelmState {
	return &helmState{
		node: node,
	}
}
