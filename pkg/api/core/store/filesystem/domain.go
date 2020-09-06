package filesystem

import (
	"fmt"
	"path"

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

// Domain contains the outputs we will store
type Domain struct {
	ID           api.ID
	FQDN         string
	Domain       string
	HostedZoneID string
	NameServers  []string
	StackName    string
}

func (d *domainStore) SaveDomain(domain *api.Domain) error {
	p := Domain{
		ID:           domain.ID,
		FQDN:         domain.FQDN,
		Domain:       domain.Domain,
		HostedZoneID: domain.HostedZoneID,
		NameServers:  domain.NameServers,
		StackName:    domain.StackName,
	}

	for i, cluster := range d.repoState.Clusters {
		if cluster.Environment == domain.ID.Environment {
			cluster.HostedZone.Domain = domain.Domain
			cluster.HostedZone.FQDN = domain.Domain
			d.repoState.Clusters[i] = cluster
		}
	}

	_, err := store.NewFileSystem(path.Join(d.paths.BaseDir, domain.Domain), d.fs).
		StoreStruct(d.paths.OutputFile, &p, store.ToJSON()).
		StoreBytes(d.paths.CloudFormationFile, domain.CloudFormationTemplate).
		AlterStore(store.SetBaseDir(d.repoPaths.BaseDir)).
		StoreStruct(d.repoPaths.ConfigFile, d.repoState, store.ToYAML()).
		Do()
	if err != nil {
		return fmt.Errorf("failed to store domain: %w", err)
	}

	return nil
}

// NewDomainStore returns an initialised domain store
func NewDomainStore(repoState *repository.Data, paths, repoPaths Paths, fs *afero.Afero) api.DomainStore {
	return &domainStore{
		paths:     paths,
		repoPaths: repoPaths,
		repoState: repoState,
		fs:        fs,
	}
}
