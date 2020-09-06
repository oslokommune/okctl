package filesystem

import (
	"fmt"
	"path"

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

// nolint: funlen
func (c *certificateStore) SaveCertificate(certificate *api.Certificate) error {
	cert := Certificate{
		ID:             certificate.ID,
		FQDN:           certificate.FQDN,
		Domain:         certificate.Domain,
		HostedZoneID:   certificate.HostedZoneID,
		CertificateARN: certificate.CertificateARN,
		StackName:      certificate.CertificateARN,
	}

	for i, cluster := range c.repoState.Clusters {
		if cluster.Environment == certificate.ID.Environment {
			found := false

			for _, storedCert := range cluster.Certificates {
				if storedCert.ARN == certificate.CertificateARN {
					found = true
					break
				}
			}

			if !found {
				cluster.Certificates = append(cluster.Certificates, repository.Certificate{
					ARN:    certificate.CertificateARN,
					Domain: certificate.Domain,
					FQDN:   certificate.FQDN,
				})
			}

			c.repoState.Clusters[i] = cluster
		}
	}

	_, err := store.NewFileSystem(path.Join(c.paths.BaseDir, certificate.Domain), c.fs).
		StoreStruct(c.paths.OutputFile, &cert, store.ToJSON()).
		StoreBytes(c.paths.CloudFormationFile, certificate.CloudFormationTemplate).
		AlterStore(store.SetBaseDir(c.repoPaths.BaseDir)).
		StoreStruct(c.repoPaths.ConfigFile, c.repoState, store.ToYAML()).
		Do()
	if err != nil {
		return fmt.Errorf("failed to store certificate: %w", err)
	}

	return nil
}

// NewCertificateStore returns an initialised certificate store
func NewCertificateStore(repoState *repository.Data, paths, repoPaths Paths, fs *afero.Afero) api.CertificateStore {
	return &certificateStore{
		paths:     paths,
		repoPaths: repoPaths,
		repoState: repoState,
		fs:        fs,
	}
}
