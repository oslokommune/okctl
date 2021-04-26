package storm

import (
	"errors"
	"time"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type componentState struct {
	node stormpkg.Node
}

// PostgresDatabase contains storm compatible state
type PostgresDatabase struct {
	Metadata `storm:"inline"`

	ID                           ID
	ApplicationName              string
	UserName                     string
	StackName                    string `storm:"unique"`
	AdminSecretFriendlyName      string
	EndpointAddress              string
	EndpointPort                 int
	OutgoingSecurityGroupID      string
	SecretsManagerAdminSecretARN string
	LambdaPolicyARN              string
	LambdaRoleARN                string
	LambdaFunctionARN            string
	CloudFormationTemplate       string
	Namespace                    string
	AdminSecretName              string
	AdminSecretARN               string
	DatabaseConfigMapName        string
	RotaterBucket                *S3Bucket
}

// NewPostgresDatabase returns storm compatible state
func NewPostgresDatabase(d *client.PostgresDatabase, meta Metadata) *PostgresDatabase {
	return &PostgresDatabase{
		Metadata:                     meta,
		ID:                           NewID(d.ID),
		ApplicationName:              d.ApplicationName,
		UserName:                     d.UserName,
		StackName:                    d.StackName,
		AdminSecretFriendlyName:      d.AdminSecretFriendlyName,
		EndpointAddress:              d.EndpointAddress,
		EndpointPort:                 d.EndpointPort,
		OutgoingSecurityGroupID:      d.OutgoingSecurityGroupID,
		SecretsManagerAdminSecretARN: d.SecretsManagerAdminSecretARN,
		LambdaPolicyARN:              d.LambdaPolicyARN,
		LambdaRoleARN:                d.LambdaRoleARN,
		LambdaFunctionARN:            d.LambdaFunctionARN,
		CloudFormationTemplate:       d.CloudFormationTemplate,
		Namespace:                    d.Namespace,
		AdminSecretName:              d.AdminSecretName,
		AdminSecretARN:               d.SecretsManagerAdminSecretARN,
		DatabaseConfigMapName:        d.DatabaseConfigMapName,
		RotaterBucket:                NewS3Bucket(d.RotaterBucket),
	}
}

// Convert to client.PostgresDatabase
func (d *PostgresDatabase) Convert() *client.PostgresDatabase {
	return &client.PostgresDatabase{
		ID:                           d.ID.Convert(),
		ApplicationName:              d.ApplicationName,
		UserName:                     d.UserName,
		StackName:                    d.StackName,
		AdminSecretFriendlyName:      d.AdminSecretFriendlyName,
		EndpointAddress:              d.EndpointAddress,
		EndpointPort:                 d.EndpointPort,
		OutgoingSecurityGroupID:      d.OutgoingSecurityGroupID,
		SecretsManagerAdminSecretARN: d.SecretsManagerAdminSecretARN,
		LambdaPolicyARN:              d.LambdaPolicyARN,
		LambdaRoleARN:                d.LambdaRoleARN,
		LambdaFunctionARN:            d.LambdaFunctionARN,
		CloudFormationTemplate:       d.CloudFormationTemplate,
		Namespace:                    d.Namespace,
		AdminSecretName:              d.AdminSecretName,
		AdminSecretARN:               d.AdminSecretARN,
		DatabaseConfigMapName:        d.DatabaseConfigMapName,
		RotaterBucket:                d.RotaterBucket.Convert(),
	}
}

// S3Bucket contains storm compatible state
type S3Bucket struct {
	Name                   string
	StackName              string
	CloudFormationTemplate string
}

// NewS3Bucket returns storm compatible state
func NewS3Bucket(b *client.S3Bucket) *S3Bucket {
	return &S3Bucket{
		Name:                   b.Name,
		StackName:              b.StackName,
		CloudFormationTemplate: b.CloudFormationTemplate,
	}
}

// Convert to client.S3Bucket
func (b *S3Bucket) Convert() *client.S3Bucket {
	return &client.S3Bucket{
		Name:                   b.Name,
		StackName:              b.StackName,
		CloudFormationTemplate: b.CloudFormationTemplate,
	}
}

func (c *componentState) SavePostgresDatabase(database *client.PostgresDatabase) error {
	existing, err := c.getPostgresDatabase(database.StackName)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return c.node.Save(NewPostgresDatabase(database, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return c.node.Save(NewPostgresDatabase(database, existing.Metadata))
}

func (c *componentState) RemovePostgresDatabase(stackName string) error {
	db := &PostgresDatabase{}

	err := c.node.One("StackName", stackName, db)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return c.node.DeleteStruct(db)
}

func (c *componentState) GetPostgresDatabase(stackName string) (*client.PostgresDatabase, error) {
	db, err := c.getPostgresDatabase(stackName)
	if err != nil {
		return nil, err
	}

	return db.Convert(), nil
}

func (c *componentState) getPostgresDatabase(stackName string) (*PostgresDatabase, error) {
	db := &PostgresDatabase{}

	err := c.node.One("StackName", stackName, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (c *componentState) GetPostgresDatabases() ([]*client.PostgresDatabase, error) {
	var dbs []*PostgresDatabase

	err := c.node.All(&dbs)
	if err != nil {
		return nil, err
	}

	ret := make([]*client.PostgresDatabase, len(dbs))

	for i, db := range dbs {
		ret[i] = db.Convert()
	}

	return ret, nil
}

// NewComponentState returns an initialised state client
func NewComponentState(node stormpkg.Node) client.ComponentState {
	return &componentState{
		node: node,
	}
}
