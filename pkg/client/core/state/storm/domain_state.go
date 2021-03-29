package storm

import (
	"errors"
	"time"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

// ErrNotFound is a not found error
var ErrNotFound = errors.New("not found")

type domainState struct {
	node stormpkg.Node
}

// HostedZone contains storm compatible state
type HostedZone struct {
	ID                     *ID
	IsDelegated            bool
	Primary                bool
	Managed                bool
	FQDN                   string
	Domain                 string `storm:"unique,index"`
	HostedZoneID           string
	NameServers            []string
	StackName              string
	CloudFormationTemplate []byte

	Metadata `storm:"inline"`
}

// NewHostedZone constructs a storm compatible HostedZone
func NewHostedZone(hz *client.HostedZone, meta Metadata) *HostedZone {
	return &HostedZone{
		ID:                     NewID(hz.ID),
		IsDelegated:            hz.IsDelegated,
		Primary:                hz.Primary,
		Managed:                hz.Managed,
		FQDN:                   hz.FQDN,
		Domain:                 hz.Domain,
		HostedZoneID:           hz.HostedZoneID,
		NameServers:            hz.NameServers,
		StackName:              hz.StackName,
		CloudFormationTemplate: hz.CloudFormationTemplate,
		Metadata:               meta,
	}
}

// Convert to a client.HostedZone
func (hz *HostedZone) Convert() *client.HostedZone {
	return &client.HostedZone{
		ID:                     hz.ID.Convert(),
		IsDelegated:            hz.IsDelegated,
		Primary:                hz.Primary,
		Managed:                hz.Managed,
		FQDN:                   hz.FQDN,
		Domain:                 hz.Domain,
		HostedZoneID:           hz.HostedZoneID,
		NameServers:            hz.NameServers,
		StackName:              hz.StackName,
		CloudFormationTemplate: hz.CloudFormationTemplate,
	}
}

func (d *domainState) SaveHostedZone(hz *client.HostedZone) error {
	return d.node.Save(NewHostedZone(hz, Metadata{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Deleted:   false,
	}))
}

func (d *domainState) UpdateHostedZone(zone *client.HostedZone) error {
	return d.node.Update(zone)
}

func (d *domainState) RemoveHostedZone(domain string) error {
	hz := &HostedZone{}

	err := d.node.One("Domain", domain, hz)
	if err != nil {
		return err
	}

	return d.node.DeleteStruct(hz)
}

func (d *domainState) GetHostedZone(domain string) (*client.HostedZone, error) {
	var hz HostedZone

	err := d.node.One("Domain", domain, &hz)
	if err != nil {
		return nil, err
	}

	return hz.Convert(), nil
}

func (d *domainState) GetHostedZones() ([]*client.HostedZone, error) {
	var hzs []*HostedZone

	err := d.node.AllByIndex("UpdatedAt", &hzs)
	if err != nil {
		return nil, err
	}

	var converted []*client.HostedZone

	for _, hz := range hzs {
		if !hz.Deleted {
			converted = append(converted, hz.Convert())
		}
	}

	return converted, nil
}

func (d *domainState) GetPrimaryHostedZone() (*client.HostedZone, error) {
	hzs, err := d.GetHostedZones()
	if err != nil {
		return nil, err
	}

	for _, hz := range hzs {
		if hz.Primary {
			return hz, nil
		}
	}

	return nil, ErrNotFound
}

// NewDomainState returns an initialised state store
func NewDomainState(db stormpkg.Node) client.DomainState {
	return &domainState{
		node: db,
	}
}