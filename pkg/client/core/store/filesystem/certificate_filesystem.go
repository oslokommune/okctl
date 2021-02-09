package filesystem

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type certificateStore struct {
	paths Paths
	fs    *afero.Afero
}

func (s *certificateStore) RemoveCertificate(domain string) (*store.Report, error) {
	return store.NewFileSystem(path.Join(s.paths.BaseDir, domain), s.fs).
		Remove(s.paths.OutputFile).
		Remove(s.paths.CloudFormationFile).
		Do()
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

func (s *certificateStore) GetCertificate(domain string) (*api.Certificate, error) {
	cert := &Certificate{}

	var template []byte

	_, err := store.NewFileSystem(path.Join(s.paths.BaseDir, domain), s.fs).
		GetStruct(s.paths.OutputFile, cert, store.FromJSON()).
		GetBytes(s.paths.CloudFormationFile, func(_ string, data []byte) {
			template = data
		}).
		Do()
	if err != nil {
		return nil, err
	}

	return &api.Certificate{
		ID:                     cert.ID,
		FQDN:                   cert.FQDN,
		Domain:                 cert.Domain,
		HostedZoneID:           cert.HostedZoneID,
		CertificateARN:         cert.CertificateARN,
		StackName:              cert.StackName,
		CloudFormationTemplate: template,
	}, nil
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

	report, err := store.NewFileSystem(path.Join(s.paths.BaseDir, c.Domain), s.fs).
		StoreStruct(s.paths.OutputFile, &cert, store.ToJSON()).
		StoreBytes(s.paths.CloudFormationFile, c.CloudFormationTemplate).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to store certificate: %w", err)
	}

	return report, nil
}

// NewCertificateStore returns an initialised certificate store
func NewCertificateStore(paths Paths, fs *afero.Afero) client.CertificateStore {
	return &certificateStore{
		paths: paths,
		fs:    fs,
	}
}
