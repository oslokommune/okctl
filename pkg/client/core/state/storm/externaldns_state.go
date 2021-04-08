package storm

import (
	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type externalDNSState struct {
	node stormpkg.Node
}

// ExternalDNS contains storm compatible state
type ExternalDNS struct {
	Metadata `storm:"inline"`

	Name string `storm:"unique,index"`
	Kube *ExternalDNSKube
}

// NewExternalDNS returns storm compatible state
func NewExternalDNS(e *client.ExternalDNS, meta Metadata) *ExternalDNS {
	return &ExternalDNS{
		Metadata: meta,
		Name:     e.Name,
		Kube:     NewExternalDNSKube(e.Kube),
	}
}

// Convert to client.ExternalDNS
func (e *ExternalDNS) Convert() *client.ExternalDNS {
	return &client.ExternalDNS{
		Name: e.Name,
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

func (e *externalDNSState) GetExternalDNS(name string) (*client.ExternalDNS, error) {
	ex := &ExternalDNS{}

	err := e.node.One("Name", name, ex)
	if err != nil {
		return nil, err
	}

	return ex.Convert(), nil
}

func (e *externalDNSState) RemoveExternalDNS(name string) error {
	ex := &ExternalDNS{}

	err := e.node.One("Name", name, ex)
	if err != nil {
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