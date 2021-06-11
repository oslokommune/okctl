package storm

import (
	"errors"
	"fmt"
	"time"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/breeze"
	"github.com/oslokommune/okctl/pkg/client"
)

type externalDNSState struct {
	node breeze.Client
}

// ExternalDNS contains storm compatible state
type ExternalDNS struct {
	Metadata `storm:"inline"`

	Name string `storm:"unique"`
	Kube *ExternalDNSKube
}

// NewExternalDNS returns storm compatible state
func NewExternalDNS(e *client.ExternalDNS, meta Metadata) *ExternalDNS {
	return &ExternalDNS{
		Metadata: meta,
		Name:     "external-dns",
		Kube:     NewExternalDNSKube(e.Kube),
	}
}

// Convert to client.ExternalDNS
func (e *ExternalDNS) Convert() *client.ExternalDNS {
	return &client.ExternalDNS{
		Kube: e.Kube.Convert(),
	}
}

// ExternalDNSKube contains storm compatible state
type ExternalDNSKube struct {
	ID           ID
	HostedZoneID string
	DomainFilter string
	Manifests    map[string]string
}

// NewExternalDNSKube returns storm compatible state
func NewExternalDNSKube(k *client.ExternalDNSKube) *ExternalDNSKube {
	manifests := map[string]string{}

	for key, v := range k.Manifests {
		manifests[key] = string(v)
	}

	return &ExternalDNSKube{
		ID:           NewID(k.ID),
		HostedZoneID: k.HostedZoneID,
		DomainFilter: k.DomainFilter,
		Manifests:    manifests,
	}
}

// Convert to client.ExternalDNSKube
func (k *ExternalDNSKube) Convert() *client.ExternalDNSKube {
	manifests := map[string][]byte{}

	for key, v := range k.Manifests {
		manifests[key] = []byte(v)
	}

	return &client.ExternalDNSKube{
		ID:           k.ID.Convert(),
		HostedZoneID: k.HostedZoneID,
		DomainFilter: k.DomainFilter,
		Manifests:    manifests,
	}
}

func (e *externalDNSState) SaveExternalDNS(dns *client.ExternalDNS) error {
	existing, err := e.getExternalDNS()
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return e.node.Save(NewExternalDNS(dns, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return e.node.Save(NewExternalDNS(dns, existing.Metadata))
}

func (e *externalDNSState) GetExternalDNS() (*client.ExternalDNS, error) {
	ex, err := e.getExternalDNS()
	if err != nil {
		return nil, err
	}

	return ex.Convert(), nil
}

func (e *externalDNSState) getExternalDNS() (*ExternalDNS, error) {
	ex := &ExternalDNS{}

	err := e.node.One("Name", "external-dns", ex)
	if err != nil {
		return nil, err
	}

	return ex, nil
}

func (e *externalDNSState) RemoveExternalDNS() error {
	ex := &ExternalDNS{}

	err := e.node.One("Name", "external-dns", ex)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return e.node.DeleteStruct(ex)
}

func (e *externalDNSState) HasExternalDNS() (bool, error) {
	_, err := e.getExternalDNS()
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("querying state: %w", err)
	}

	return true, nil
}

// NewExternalDNSState returns an initialised state
func NewExternalDNSState(node breeze.Client) client.ExternalDNSState {
	return &externalDNSState{
		node: node,
	}
}
