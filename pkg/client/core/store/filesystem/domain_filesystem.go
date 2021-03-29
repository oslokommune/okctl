package filesystem

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/spf13/afero"
)

type domainStore struct {
	paths Paths
	fs    *afero.Afero
}

func (s *domainStore) RemoveHostedZone(domain string) (*store.Report, error) {
	report, err := store.NewFileSystem(path.Join(s.paths.BaseDir, domain), s.fs).
		Remove(s.paths.CloudFormationFile).
		AlterStore(store.SetBaseDir(s.paths.BaseDir)).
		RemoveDir(domain).
		RemoveDir("").
		Do()
	if err != nil {
		return nil, fmt.Errorf("removing hosted zone: %s", err)
	}

	return report, nil
}

func (s *domainStore) SaveHostedZone(d *client.HostedZone) (*store.Report, error) {
	report, err := store.NewFileSystem(path.Join(s.paths.BaseDir, d.Domain), s.fs).
		StoreBytes(s.paths.CloudFormationFile, d.CloudFormationTemplate).
		Do()
	if err != nil {
		return nil, fmt.Errorf("storing hosted zone: %w", err)
	}

	return report, nil
}

// NewDomainStore returns an initialised domain store
func NewDomainStore(paths Paths, fs *afero.Afero) client.DomainStore {
	return &domainStore{
		paths: paths,
		fs:    fs,
	}
}
