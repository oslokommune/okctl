package filesystem

import (
	"fmt"
	"path"

	"github.com/gosimple/slug"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type identityManagerStore struct {
	poolPaths  Paths
	certPaths  Paths
	aliasPaths Paths
	fs         *afero.Afero
}

// IdentityPool contains the output state
type IdentityPool struct {
	ID             api.ID
	UserPoolID     string
	AuthDomain     string
	HostedZoneID   string
	StackName      string
	Clients        []IdentityClient
	Certificate    Certificate
	RecordSetAlias RecordSetAlias
}

// IdentityClient contains the output state
type IdentityClient struct {
	Purpose     string
	CallbackURL string
	ClientID    string
}

// RecordSetAlias contains the output state
type RecordSetAlias struct {
	AliasDomain      string
	AliasHostedZones string
	StackName        string
}

func (s *identityManagerStore) SaveIdentityPool(pool *api.IdentityPool) (*store.Report, error) {
	p := &IdentityPool{
		ID:           pool.ID,
		UserPoolID:   pool.UserPoolID,
		AuthDomain:   pool.AuthDomain,
		HostedZoneID: pool.HostedZoneID,
		StackName:    pool.StackName,
		Clients:      nil,
		Certificate: Certificate{
			ID:             pool.ID,
			FQDN:           pool.Certificate.FQDN,
			Domain:         pool.Certificate.Domain,
			HostedZoneID:   pool.Certificate.HostedZoneID,
			CertificateARN: pool.Certificate.CertificateARN,
			StackName:      pool.Certificate.StackName,
		},
		RecordSetAlias: RecordSetAlias{
			AliasDomain:      pool.RecordSetAlias.AliasDomain,
			AliasHostedZones: pool.RecordSetAlias.AliasHostedZones,
			StackName:        pool.RecordSetAlias.StackName,
		},
	}

	for _, c := range pool.Clients {
		p.Clients = append(p.Clients, IdentityClient{
			Purpose:     c.Purpose,
			CallbackURL: c.CallbackURL,
			ClientID:    c.ClientID,
		})
	}

	report, err := store.NewFileSystem(s.poolPaths.BaseDir, s.fs).
		StoreStruct(s.poolPaths.OutputFile, p, store.ToJSON()).
		StoreBytes(s.poolPaths.CloudFormationFile, pool.CloudFormationTemplates).
		AlterStore(store.SetBaseDir(path.Join(s.certPaths.BaseDir, slug.Make(p.Certificate.Domain)))).
		StoreBytes(s.certPaths.CloudFormationFile, pool.Certificate.CloudFormationTemplate).
		AlterStore(store.SetBaseDir(path.Join(s.aliasPaths.CloudFormationFile, slug.Make(p.AuthDomain)))).
		StoreBytes(s.aliasPaths.CloudFormationFile, pool.RecordSetAlias.CloudFormationTemplate).
		Do()
	if err != nil {
		return nil, fmt.Errorf("writing identity pool: %w", err)
	}

	return report, nil
}

// NewIdentityManagerStore returns an initialised store
func NewIdentityManagerStore(poolPaths, certPaths, aliasPaths Paths, fs *afero.Afero) client.IdentityManagerStore {
	return &identityManagerStore{
		poolPaths:  poolPaths,
		certPaths:  certPaths,
		aliasPaths: aliasPaths,
		fs:         fs,
	}
}
