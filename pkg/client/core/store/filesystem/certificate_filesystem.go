package filesystem

import (
	"path"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/spf13/afero"
)

type certificateStore struct {
	paths Paths
	fs    *afero.Afero
}

func (s *certificateStore) RemoveCertificate(domain string) error {
	_, err := store.NewFileSystem(path.Join(s.paths.BaseDir, domain), s.fs).
		Remove(s.paths.CloudFormationFile).
		AlterStore(store.SetBaseDir(s.paths.BaseDir)).
		RemoveDir(domain).
		RemoveDir("").
		Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *certificateStore) SaveCertificate(c *client.Certificate) error {
	_, err := store.NewFileSystem(path.Join(s.paths.BaseDir, c.Domain), s.fs).
		StoreBytes(s.paths.CloudFormationFile, c.CloudFormationTemplate).
		Do()
	if err != nil {
		return err
	}

	return nil
}

// NewCertificateStore returns an initialised certificate store
func NewCertificateStore(paths Paths, fs *afero.Afero) client.CertificateStore {
	return &certificateStore{
		paths: paths,
		fs:    fs,
	}
}
