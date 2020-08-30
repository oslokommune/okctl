package filesystem

import (
	"encoding/json"
	"fmt"
	"path"

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
	Repository     string
	Environment    string
	FQDN           string
	Domain         string
	HostedZoneID   string
	CertificateARN string
	StackName      string
}

// nolint: funlen
func (c *certificateStore) SaveCertificate(certificate *api.Certificate) error {
	cert := Certificate{
		Repository:     certificate.Repository,
		Environment:    certificate.Environment,
		FQDN:           certificate.FQDN,
		Domain:         certificate.Domain,
		HostedZoneID:   certificate.HostedZoneID,
		CertificateARN: certificate.CertificateARN,
		StackName:      certificate.CertificateARN,
	}

	data, err := json.Marshal(cert)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	err = c.fs.MkdirAll(path.Join(c.paths.BaseDir, certificate.Domain), 0o744)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err = c.fs.WriteFile(path.Join(c.paths.BaseDir, certificate.Domain, c.paths.OutputFile), data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write outputs: %w", err)
	}

	err = c.fs.WriteFile(path.Join(c.paths.BaseDir, certificate.Domain, c.paths.CloudFormationFile), certificate.CloudFormationTemplate, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write cloud formation template: %w", err)
	}

	for i, cluster := range c.repoState.Clusters {
		if cluster.Environment == certificate.Environment {
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

	err = c.fs.MkdirAll(c.repoPaths.BaseDir, 0o744)
	if err != nil {
		return err
	}

	state, err := c.repoState.YAML()
	if err != nil {
		return err
	}

	return c.fs.WriteFile(path.Join(c.repoPaths.BaseDir, c.repoPaths.ConfigFile), state, 0o644)
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
