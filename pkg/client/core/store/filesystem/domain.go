package filesystem

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type domainStore struct {
	paths Paths
	fs    *afero.Afero
}

// HostedZone contains the outputs we will store
type HostedZone struct {
	ID           api.ID
	FQDN         string
	Domain       string
	HostedZoneID string
	NameServers  []string
	StackName    string
	IsDelegated  bool
	Primary      bool
}

func (s *domainStore) RemoveHostedZone(domain string) (*store.Report, error) {
	report, err := store.NewFileSystem(path.Join(s.paths.BaseDir, domain), s.fs).
		Remove(s.paths.OutputFile).
		Remove(s.paths.CloudFormationFile).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to remove hosted zone: %s", err)
	}

	return report, nil
}

func (s *domainStore) SaveHostedZone(d *client.HostedZone) (*store.Report, error) {
	p := HostedZone{
		ID:           d.HostedZone.ID,
		FQDN:         d.HostedZone.FQDN,
		Domain:       d.HostedZone.Domain,
		HostedZoneID: d.HostedZone.HostedZoneID,
		NameServers:  d.HostedZone.NameServers,
		StackName:    d.HostedZone.StackName,
		IsDelegated:  d.IsDelegated,
		Primary:      d.Primary,
	}

	report, err := store.NewFileSystem(path.Join(s.paths.BaseDir, d.HostedZone.Domain), s.fs).
		StoreStruct(s.paths.OutputFile, &p, store.ToJSON()).
		StoreBytes(s.paths.CloudFormationFile, d.HostedZone.CloudFormationTemplate).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to store hosted zone: %w", err)
	}

	return report, nil
}

func (s *domainStore) GetHostedZone(domain string) (*client.HostedZone, error) {
	hz := &HostedZone{}

	var template []byte

	_, err := store.NewFileSystem(path.Join(s.paths.BaseDir, domain), s.fs).
		GetStruct(s.paths.OutputFile, hz, store.FromJSON()).
		GetBytes(s.paths.CloudFormationFile, func(_ string, data []byte) {
			template = data
		}).
		Do()
	if err != nil {
		return nil, err
	}

	return &client.HostedZone{
		IsDelegated: hz.IsDelegated,
		Primary:     hz.Primary,
		HostedZone: &api.HostedZone{
			ID:                     hz.ID,
			FQDN:                   hz.FQDN,
			Domain:                 hz.Domain,
			HostedZoneID:           hz.HostedZoneID,
			NameServers:            hz.NameServers,
			StackName:              hz.StackName,
			CloudFormationTemplate: template,
		},
	}, nil
}

// NewDomainStore returns an initialised domain store
func NewDomainStore(paths Paths, fs *afero.Afero) client.DomainStore {
	return &domainStore{
		paths: paths,
		fs:    fs,
	}
}
