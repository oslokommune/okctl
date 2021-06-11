package storm

import (
	"errors"
	"time"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/breeze"
	"github.com/oslokommune/okctl/pkg/client"
)

type domainState struct {
	node breeze.Client
}

// HostedZone contains storm compatible state
type HostedZone struct {
	Metadata `storm:"inline"`

	ID                     ID
	IsDelegated            bool
	Primary                bool
	Managed                bool
	FQDN                   string
	Domain                 string `storm:"unique"`
	HostedZoneID           string
	NameServers            []string
	StackName              string
	CloudFormationTemplate string
}

// NewHostedZone constructs a storm compatible HostedZone
func NewHostedZone(hz *client.HostedZone, meta Metadata) *HostedZone {
	return &HostedZone{
		Metadata:               meta,
		ID:                     NewID(hz.ID),
		IsDelegated:            hz.IsDelegated,
		Primary:                hz.Primary,
		Managed:                hz.Managed,
		FQDN:                   hz.FQDN,
		Domain:                 hz.Domain,
		HostedZoneID:           hz.HostedZoneID,
		NameServers:            hz.NameServers,
		StackName:              hz.StackName,
		CloudFormationTemplate: string(hz.CloudFormationTemplate),
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
		CloudFormationTemplate: []byte(hz.CloudFormationTemplate),
	}
}

func (d *domainState) SaveHostedZone(hz *client.HostedZone) error {
	existing, err := d.getHostedZone(hz.Domain)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return d.node.Save(NewHostedZone(hz, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return d.node.Save(NewHostedZone(hz, existing.Metadata))
}

func (d *domainState) UpdateHostedZone(zone *client.HostedZone) error {
	return d.SaveHostedZone(zone)
}

func (d *domainState) RemoveHostedZone(domain string) error {
	hz := &HostedZone{}

	err := d.node.One("Domain", domain, hz)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return d.node.DeleteStruct(hz)
}

func (d *domainState) GetHostedZone(domain string) (*client.HostedZone, error) {
	hz, err := d.getHostedZone(domain)
	if err != nil {
		return nil, err
	}

	return hz.Convert(), nil
}

func (d *domainState) getHostedZone(domain string) (*HostedZone, error) {
	hz := &HostedZone{}

	err := d.node.One("Domain", domain, hz)
	if err != nil {
		return nil, err
	}

	return hz, nil
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

	return nil, stormpkg.ErrNotFound
}

func (d *domainState) HasPrimaryHostedZone() (bool, error) {
	_, err := d.GetPrimaryHostedZone()
	if err == nil {
		return true, nil
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return false, nil
	}

	return false, err
}

// NewDomainState returns an initialised state store
func NewDomainState(db breeze.Client) client.DomainState {
	return &domainState{
		node: db,
	}
}
