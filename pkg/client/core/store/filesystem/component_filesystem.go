package filesystem

import (
	"path"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type componentStore struct {
	paths Paths
	fs    *afero.Afero
}

// PostgresDatabase contains the state we store
// in the output file
type PostgresDatabase struct {
	ID                           api.ID
	ApplicationName              string
	UserName                     string
	StackName                    string
	AdminSecretFriendlyName      string
	EndpointAddress              string
	EndpointPort                 int
	OutgoingSecurityGroupID      string
	SecretsManagerAdminSecretARN string
	Namespace                    string
	AdminSecretName              string
	DatabaseConfigMapName        string
}

func (c *componentStore) SavePostgresDatabase(db *client.PostgresDatabase) (*store.Report, error) {
	pg := &PostgresDatabase{
		ID:                           db.ID,
		ApplicationName:              db.ApplicationName,
		UserName:                     db.UserName,
		StackName:                    db.StackName,
		AdminSecretFriendlyName:      db.AdminSecretFriendlyName,
		EndpointAddress:              db.EndpointAddress,
		EndpointPort:                 db.EndpointPort,
		OutgoingSecurityGroupID:      db.OutgoingSecurityGroupID,
		SecretsManagerAdminSecretARN: db.SecretsManagerAdminSecretARN,
		Namespace:                    db.Namespace,
		AdminSecretName:              db.AdminSecretName,
		DatabaseConfigMapName:        db.DatabaseConfigMapName,
	}

	report, err := store.NewFileSystem(path.Join(c.paths.BaseDir, db.ApplicationName), c.fs).
		StoreStruct(c.paths.OutputFile, pg, store.ToJSON()).
		StoreBytes(c.paths.CloudFormationFile, []byte(db.CloudFormationTemplate)).
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (c *componentStore) RemovePostgresDatabase(applicationName string) (*store.Report, error) {
	return store.NewFileSystem(path.Join(c.paths.BaseDir, applicationName), c.fs).
		Remove(c.paths.OutputFile).
		Remove(c.paths.CloudFormationFile).
		AlterStore(store.SetBaseDir(c.paths.BaseDir)).
		RemoveDir(applicationName).
		Do()
}

// NewComponentStore returns an initialised component store
func NewComponentStore(paths Paths, fs *afero.Afero) client.ComponentStore {
	return &componentStore{
		paths: paths,
		fs:    fs,
	}
}
