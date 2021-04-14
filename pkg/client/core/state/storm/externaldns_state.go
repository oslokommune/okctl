package storm

import (
	"errors"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type externalDNSState struct {
	node stormpkg.Node
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
	Manifests    map[string][]byte
}

// NewExternalDNSKube returns storm compatible state
func NewExternalDNSKube(k *client.ExternalDNSKube) *ExternalDNSKube {
	return &ExternalDNSKube{
		ID:           NewID(k.ID),
		HostedZoneID: k.HostedZoneID,
		DomainFilter: k.DomainFilter,
		Manifests:    k.Manifests,
	}
}

// Convert to client.ExternalDNSKube
func (k *ExternalDNSKube) Convert() *client.ExternalDNSKube {
	return &client.ExternalDNSKube{
		ID:           k.ID.Convert(),
		HostedZoneID: k.HostedZoneID,
		DomainFilter: k.DomainFilter,
		Manifests:    k.Manifests,
	}
}

func (e *externalDNSState) SaveExternalDNS(dns *client.ExternalDNS) error {
	return e.node.Save(NewExternalDNS(dns, NewMetadata()))
}

func (e *externalDNSState) GetExternalDNS() (*client.ExternalDNS, error) {
	ex := &ExternalDNS{}

	err := e.node.One("Name", "external-dns", ex)
	if err != nil {
		return nil, err
	}

	return ex.Convert(), nil
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

// NewExternalDNSState returns an initialised state
func NewExternalDNSState(node stormpkg.Node) client.ExternalDNSState {
	return &externalDNSState{
		node: node,
	}
}
