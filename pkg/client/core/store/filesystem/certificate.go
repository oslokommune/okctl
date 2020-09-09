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

type certificateStore struct {
	paths     Paths
	repoPaths Paths
	repoState *repository.Data
	fs        *afero.Afero
}

// Certificate contains the data we store to the outputs
type Certificate struct {
	ID             api.ID
	FQDN           string
	Domain         string
	HostedZoneID   string
	CertificateARN string
	StackName      string
}

func (s *certificateStore) SaveCertificate(c *api.Certificate) (*store.Report, error) {
	cert := Certificate{
		ID:             c.ID,
		FQDN:           c.FQDN,
		Domain:         c.Domain,
		HostedZoneID:   c.HostedZoneID,
		CertificateARN: c.CertificateARN,
		StackName:      c.CertificateARN,
	}

	cluster, ok := s.repoState.Clusters[c.ID.Environment]
	if !ok {
		return nil, fmt.Errorf("found no cluster for environment: %s", c.ID.Environment)
	}

	if cluster.Certificates == nil {
		cluster.Certificates = map[string]string{}
	}

	cluster.Certificates[c.Domain] = c.CertificateARN

	report, err := store.NewFileSystem(path.Join(s.paths.BaseDir, c.Domain), s.fs).
		StoreStruct(s.paths.OutputFile, &cert, store.ToJSON()).
		StoreBytes(s.paths.CloudFormationFile, c.CloudFormationTemplate).
		AlterStore(store.SetBaseDir(s.repoPaths.BaseDir)).
		StoreStruct(s.repoPaths.ConfigFile, s.repoState, store.ToYAML()).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to store certificate: %w", err)
	}

	return report, nil
}

// NewCertificateStore returns an initialised certificate store
func NewCertificateStore(repoState *repository.Data, paths, repoPaths Paths, fs *afero.Afero) client.CertificateStore {
	return &certificateStore{
		paths:     paths,
		repoPaths: repoPaths,
		repoState: repoState,
		fs:        fs,
	}
}
