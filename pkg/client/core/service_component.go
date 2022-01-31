package core

import (
	"bytes"
	"context"
	stderrors "errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/smapi"

	"github.com/oslokommune/okctl/pkg/ec2api"

	"github.com/oslokommune/okctl/pkg/static/rotater"

	"github.com/oslokommune/okctl/pkg/iamapi"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/s3api"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type componentService struct {
	service  api.ComponentService
	state    client.ComponentState
	manifest client.ManifestService
	provider v1alpha1.CloudProvider
}

func rotaterBucketName(clusterName, applicationName string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+") // nolint: gocritic
	if err != nil {
		return "", err
	}

	return reg.ReplaceAllString(fmt.Sprintf("%s%slambdas", clusterName, applicationName), ""), nil
}

func adminSecretName(applicationName string) string {
	return fmt.Sprintf("%s-postgres-admin", applicationName)
}

func pgConfigMapName(applicationName string) string {
	return fmt.Sprintf("%s-postgres", applicationName)
}

const (
	postgresRotaterLambdaKey = "rotate-postgres-single.zip"
)

//nolint:funlen,gocyclo
func (c *componentService) CreatePostgresDatabase(ctx context.Context, opts client.CreatePostgresDatabaseOpts) (*client.PostgresDatabase, error) {
	bucketName, err := rotaterBucketName(opts.ID.ClusterName, opts.ApplicationName)
	if err != nil {
		return nil, err
	}

	bucket, err := c.service.CreateS3Bucket(ctx, &api.CreateS3BucketOpts{
		ID:        opts.ID,
		Name:      bucketName,
		StackName: cfn.NewStackNamer().S3Bucket(opts.ApplicationName, opts.ID.ClusterName),
	})
	if err != nil {
		return nil, err
	}

	err = s3api.New(c.provider).PutObject(
		bucket.Name,
		postgresRotaterLambdaKey,
		bytes.NewReader(rotater.LambdaFunctionZip),
	)
	if err != nil {
		return nil, err
	}

	pg, err := c.service.CreatePostgresDatabase(ctx, &api.CreatePostgresDatabaseOpts{
		ID:                opts.ID,
		ApplicationName:   opts.ApplicationName,
		UserName:          opts.UserName,
		StackName:         cfn.NewStackNamer().RDSPostgres(opts.ApplicationName, opts.ID.ClusterName),
		VpcID:             opts.VpcID,
		DBSubnetGroupName: opts.DBSubnetGroupName,
		DBSubnetIDs:       opts.DBSubnetIDs,
		DBSubnetCIDRs:     opts.DBSubnetCIDRs,
		RotaterBucket:     bucket.Name,
		RotaterKey:        postgresRotaterLambdaKey,
	})
	if err != nil {
		return nil, err
	}

	err = iamapi.New(c.provider).AttachRolePolicy(pg.LambdaPolicyARN, pg.LambdaRoleARN)
	if err != nil {
		return nil, err
	}

	// We cannot enable the secret rotation until all the policies have been created,
	// we couldn't find a way of doing this in cloud formation. Therefore we do it
	// here, after the policy has been attached to the role
	err = smapi.New(c.provider).RotateSecret(pg.LambdaFunctionARN, pg.SecretsManagerAdminSecretARN)
	if err != nil {
		return nil, err
	}

	_, err = c.manifest.CreateNamespace(ctx, api.CreateNamespaceOpts{
		ID:        opts.ID,
		Namespace: opts.Namespace,
	})
	if err != nil {
		return nil, err
	}

	_, err = c.manifest.CreateExternalSecret(ctx, client.CreateExternalSecretOpts{
		ID:        opts.ID,
		Name:      adminSecretName(opts.ApplicationName),
		Namespace: opts.Namespace,
		Manifest: api.Manifest{
			Name:      adminSecretName(opts.ApplicationName),
			Namespace: opts.Namespace,
			Backend:   api.BackendTypeSecretsManager,
			Data: []api.Data{
				{
					Key:      pg.AdminSecretFriendlyName,
					Name:     "PGUSER",
					Property: "username",
				},
				{
					Key:      pg.AdminSecretFriendlyName,
					Name:     "PGPASSWORD",
					Property: "password",
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = c.manifest.CreateConfigMap(ctx, client.CreateConfigMapOpts{
		ID:        opts.ID,
		Name:      pgConfigMapName(opts.ApplicationName),
		Namespace: opts.Namespace,
		Data: map[string]string{
			"PGHOST":     pg.EndpointAddress,
			"PGPORT":     strconv.Itoa(pg.EndpointPort),
			"PGDATABASE": opts.ApplicationName,
			"PGSSLMODE":  "disable", // We should enable this eventually
		},
	})
	if err != nil {
		return nil, err
	}

	postgres := &client.PostgresDatabase{
		ID:                           pg.ID,
		ApplicationName:              pg.ApplicationName,
		UserName:                     pg.UserName,
		StackName:                    pg.StackName,
		AdminSecretFriendlyName:      pg.AdminSecretFriendlyName,
		EndpointAddress:              pg.EndpointAddress,
		EndpointPort:                 pg.EndpointPort,
		OutgoingSecurityGroupID:      pg.OutgoingSecurityGroupID,
		SecretsManagerAdminSecretARN: pg.SecretsManagerAdminSecretARN,
		LambdaPolicyARN:              pg.LambdaPolicyARN,
		LambdaRoleARN:                pg.LambdaRoleARN,
		LambdaFunctionARN:            pg.LambdaFunctionARN,
		CloudFormationTemplate:       pg.CloudFormationTemplate,
		Namespace:                    opts.Namespace,
		AdminSecretName:              adminSecretName(opts.ApplicationName),
		AdminSecretARN:               pg.SecretsManagerAdminSecretARN,
		DatabaseConfigMapName:        pgConfigMapName(opts.ApplicationName),
		RotaterBucket: &client.S3Bucket{
			Name:                   bucket.Name,
			StackName:              bucket.StackName,
			CloudFormationTemplate: bucket.CloudFormationTemplate,
		},
	}

	err = ec2api.New(c.provider).AuthorizePodToNodeGroupTraffic(
		"ng-generic",
		pg.OutgoingSecurityGroupID,
		opts.VpcID,
	)
	if err != nil {
		return nil, err
	}

	err = c.state.SavePostgresDatabase(postgres)
	if err != nil {
		return nil, err
	}

	return postgres, nil
}

// nolint: funlen,gocyclo
func (c *componentService) DeletePostgresDatabase(ctx context.Context, opts client.DeletePostgresDatabaseOpts) error {
	stackName := cfn.NewStackNamer().RDSPostgres(opts.ApplicationName, opts.ID.ClusterName)

	db, err := c.state.GetPostgresDatabase(stackName)
	if err != nil {
		return err
	}

	err = ec2api.New(c.provider).RevokePodToNodeGroupTraffic(
		"ng-generic",
		db.OutgoingSecurityGroupID,
		opts.VpcID,
	)
	if err != nil && !stderrors.Is(err, ec2api.ErrNotFound) {
		return err
	}

	err = smapi.New(c.provider).CancelRotateSecret(db.AdminSecretARN)
	if err != nil {
		return err
	}

	err = iamapi.New(c.provider).DetachRolePolicy(db.LambdaPolicyARN, db.LambdaRoleARN)
	if err != nil {
		return err
	}

	err = c.service.DeletePostgresDatabase(ctx, &api.DeletePostgresDatabaseOpts{
		ID:        opts.ID,
		StackName: stackName,
	})
	if err != nil {
		return err
	}

	err = c.manifest.DeleteConfigMap(ctx, client.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      pgConfigMapName(opts.ApplicationName),
		Namespace: db.Namespace,
	})
	if err != nil {
		return err
	}

	err = c.manifest.DeleteExternalSecret(ctx, client.DeleteExternalSecretOpts{
		ID:   opts.ID,
		Name: adminSecretName(opts.ApplicationName),
		Secrets: map[string]string{
			adminSecretName(opts.ApplicationName): db.Namespace,
		},
	})
	if err != nil {
		return err
	}

	bucketName, err := rotaterBucketName(opts.ID.ClusterName, opts.ApplicationName)
	if err != nil {
		return err
	}

	err = s3api.New(c.provider).DeleteObject(
		bucketName,
		postgresRotaterLambdaKey,
	)
	if err != nil && !stderrors.Is(err, s3api.ErrBucketDoesNotExist) {
		return err
	}

	err = c.service.DeleteS3Bucket(ctx, &api.DeleteS3BucketOpts{
		ID:        opts.ID,
		StackName: cfn.NewStackNamer().S3Bucket(opts.ApplicationName, opts.ID.ClusterName),
	})
	if err != nil && errors.IsKind(err, errors.NotExist) {
		return err
	}

	err = c.manifest.DeleteNamespace(ctx, api.DeleteNamespaceOpts{
		ID:        opts.ID,
		Namespace: opts.Namespace,
	})
	if err != nil {
		return err
	}

	err = c.state.RemovePostgresDatabase(stackName)
	if err != nil {
		return err
	}

	return nil
}

func (c *componentService) GetPostgresDatabase(_ context.Context, opts client.GetPostgresDatabaseOpts) (*client.PostgresDatabase, error) {
	stackName := cfn.NewStackNamer().RDSPostgres(opts.DatabaseName, opts.ClusterID.ClusterName)

	db, err := c.state.GetPostgresDatabase(stackName)
	if err != nil {
		return nil, fmt.Errorf("fetching database from state: %w", err)
	}

	return db, nil
}

// NewComponentService returns an initialised component service
func NewComponentService(
	service api.ComponentService,
	state client.ComponentState,
	manifest client.ManifestService,
	provider v1alpha1.CloudProvider,
) client.ComponentService {
	return &componentService{
		service:  service,
		state:    state,
		manifest: manifest,
		provider: provider,
	}
}
