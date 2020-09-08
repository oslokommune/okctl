package filesystem

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/spf13/afero"
)

type domainStore struct {
	paths     Paths
	repoPaths Paths
	repoState *repository.Data
	fs        *afero.Afero
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

	cluster, ok := s.repoState.Clusters[d.HostedZone.ID.Environment]
	if !ok {
		return nil, fmt.Errorf("failed to find cluster for env: %s", d.HostedZone.ID.Environment)
	}

	cluster.HostedZone[d.HostedZone.Domain] = &repository.HostedZone{
		IsCreated:   true,
		IsDelegated: d.IsDelegated,
		Primary:     d.Primary,
		Domain:      d.HostedZone.Domain,
		FQDN:        d.HostedZone.FQDN,
		NameServers: d.HostedZone.NameServers,
	}

	subDir := d.HostedZone.Domain
	if d.Primary {
		subDir = "primary"
	}

	report, err := store.NewFileSystem(path.Join(s.paths.BaseDir, subDir), s.fs).
		StoreStruct(s.paths.OutputFile, &p, store.ToJSON()).
		StoreBytes(s.paths.CloudFormationFile, d.HostedZone.CloudFormationTemplate).
		AlterStore(store.SetBaseDir(s.repoPaths.BaseDir)).
		StoreStruct(s.repoPaths.ConfigFile, s.repoState, store.ToYAML()).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to store hosted zone: %w", err)
	}

	return report, nil
}

func (s *domainStore) GetPrimaryHostedZone(id api.ID) (*client.HostedZone, error) {
	hz := &HostedZone{}

	var template []byte

	callback := func(_ string, data []byte) error {
		template = data
		return nil
	}

	_, err := store.NewFileSystem(path.Join(s.paths.BaseDir, "primary"), s.fs).
		GetStruct(s.paths.OutputFile, hz, store.FromJSON()).
		GetBytes(s.paths.CloudFormationFile, callback).
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
func NewDomainStore(repoState *repository.Data, paths, repoPaths Paths, fs *afero.Afero) client.DomainStore {
	return &domainStore{
		paths:     paths,
		repoPaths: repoPaths,
		repoState: repoState,
		fs:        fs,
	}
}
