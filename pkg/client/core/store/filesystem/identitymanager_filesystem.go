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
	poolPaths   Paths
	certPaths   Paths
	aliasPaths  Paths
	clientPaths Paths
	userPaths   Paths
	fs          *afero.Afero
}

// IdentityPool contains the output state
type IdentityPool struct {
	ID             api.ID
	UserPoolID     string
	AuthDomain     string
	HostedZoneID   string
	StackName      string
	Certificate    client.Certificate
	RecordSetAlias RecordSetAlias
}

// RecordSetAlias contains the output state
type RecordSetAlias struct {
	AliasDomain      string
	AliasHostedZones string
	StackName        string
}

// IdentityPoolClient contains the output state
type IdentityPoolClient struct {
	ID          api.ID
	UserPoolID  string
	Purpose     string
	CallbackURL string
	ClientID    string
	StackName   string
}

// IdentityPoolUser contains the output state
type IdentityPoolUser struct {
	ID         api.ID
	Email      string
	UserPoolID string
	StackName  string
}

// RemoveIdentityPool removes alias, cert, and identitypool files from file store
func (s *identityManagerStore) RemoveIdentityPool(id api.ID) (*store.Report, error) {
	// nolint: godox
	// TODO: send domain info in here somehow, or read the file system
	authDomain := "auth-" + id.ClusterName + "-oslo-systems/"

	report, err := store.NewFileSystem(s.aliasPaths.BaseDir, s.fs).
		Remove(authDomain + s.aliasPaths.CloudFormationFile).
		Remove(authDomain + s.aliasPaths.OutputFile).
		RemoveDir("").
		AlterStore(store.SetBaseDir(s.certPaths.BaseDir)).
		Remove(authDomain + s.certPaths.CloudFormationFile).
		Remove(authDomain + s.certPaths.OutputFile).
		RemoveDir("").
		AlterStore(store.SetBaseDir(s.poolPaths.BaseDir)).
		Remove(s.poolPaths.CloudFormationFile).
		Remove(s.poolPaths.OutputFile).
		RemoveDir("").
		AlterStore(store.SetBaseDir(s.poolPaths.BaseDir)).
		RemoveDir("").
		Do()
	if err != nil {
		return nil, err
	}

	return report, err
}

func (s *identityManagerStore) RemoveIdentityPoolClient(opts api.DeleteIdentityPoolClientOpts) (*store.Report, error) {
	return store.NewFileSystem(path.Join(s.clientPaths.BaseDir, opts.Purpose), s.fs).
		Remove(s.clientPaths.OutputFile).
		Remove(s.clientPaths.CloudFormationFile).
		Do()
}

func (s *identityManagerStore) SaveIdentityPoolClient(client *api.IdentityPoolClient) (*store.Report, error) {
	c := &IdentityPoolClient{
		ID:          client.ID,
		UserPoolID:  client.UserPoolID,
		Purpose:     client.Purpose,
		CallbackURL: client.CallbackURL,
		ClientID:    client.ClientID,
		StackName:   client.StackName,
	}

	report, err := store.NewFileSystem(path.Join(s.clientPaths.BaseDir, client.Purpose), s.fs).
		StoreStruct(s.clientPaths.OutputFile, c, store.ToJSON()).
		StoreBytes(s.clientPaths.CloudFormationFile, client.CloudFormationTemplates).
		Do()
	if err != nil {
		return nil, fmt.Errorf("storing the identity client: %w", err)
	}

	return report, nil
}

func (s *identityManagerStore) SaveIdentityPool(pool *api.IdentityPool) (*store.Report, error) {
	p := &IdentityPool{
		ID:           pool.ID,
		UserPoolID:   pool.UserPoolID,
		AuthDomain:   pool.AuthDomain,
		HostedZoneID: pool.HostedZoneID,
		StackName:    pool.StackName,
		Certificate: client.Certificate{
			ID:           pool.ID,
			FQDN:         pool.Certificate.FQDN,
			Domain:       pool.Certificate.Domain,
			HostedZoneID: pool.Certificate.HostedZoneID,
			ARN:          pool.Certificate.CertificateARN,
			StackName:    pool.Certificate.StackName,
		},
		RecordSetAlias: RecordSetAlias{
			AliasDomain:      pool.RecordSetAlias.AliasDomain,
			AliasHostedZones: pool.RecordSetAlias.AliasHostedZones,
			StackName:        pool.RecordSetAlias.StackName,
		},
	}

	report, err := store.NewFileSystem(s.poolPaths.BaseDir, s.fs).
		StoreStruct(s.poolPaths.OutputFile, p, store.ToJSON()).
		StoreBytes(s.poolPaths.CloudFormationFile, pool.CloudFormationTemplates).
		AlterStore(store.SetBaseDir(path.Join(s.certPaths.BaseDir, slug.Make(p.Certificate.Domain)))).
		StoreBytes(s.certPaths.CloudFormationFile, pool.Certificate.CloudFormationTemplate).
		AlterStore(store.SetBaseDir(path.Join(s.aliasPaths.BaseDir, slug.Make(p.AuthDomain)))).
		StoreBytes(s.aliasPaths.CloudFormationFile, pool.RecordSetAlias.CloudFormationTemplate).
		Do()
	if err != nil {
		return nil, fmt.Errorf("writing identity pool: %w", err)
	}

	return report, nil
}

func (s *identityManagerStore) SaveIdentityPoolUser(user *api.IdentityPoolUser) (*store.Report, error) {
	u := &IdentityPoolUser{
		ID:         user.ID,
		Email:      user.Email,
		UserPoolID: user.UserPoolID,
		StackName:  user.StackName,
	}

	report, err := store.NewFileSystem(path.Join(s.userPaths.BaseDir, slug.Make(u.Email)), s.fs).
		StoreStruct(s.userPaths.OutputFile, u, store.ToJSON()).
		StoreBytes(s.userPaths.CloudFormationFile, user.CloudFormationTemplate).
		Do()
	if err != nil {
		return nil, fmt.Errorf("writing identity user: %w", err)
	}

	return report, nil
}

// NewIdentityManagerStore returns an initialised store
func NewIdentityManagerStore(poolPaths, certPaths, aliasPaths, clientPaths, userPaths Paths, fs *afero.Afero) client.IdentityManagerStore {
	return &identityManagerStore{
		clientPaths: clientPaths,
		poolPaths:   poolPaths,
		certPaths:   certPaths,
		aliasPaths:  aliasPaths,
		userPaths:   userPaths,
		fs:          fs,
	}
}
