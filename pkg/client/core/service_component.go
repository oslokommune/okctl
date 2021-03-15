package core

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/oslokommune/okctl/pkg/s3api"
	"github.com/oslokommune/okctl/pkg/static"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/spinner"
)

type componentService struct {
	spinner spinner.Spinner
	api     client.ComponentAPI
	store   client.ComponentStore
	state   client.ComponentState
	report  client.ComponentReport

	manifest client.ManifestService

	provider v1alpha1.CloudProvider
}

func rotaterBucketName(clusterName, applicationName string) string {
	return fmt.Sprintf("%s-%s-lambdas", clusterName, applicationName)
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

// nolint: funlen
func (c *componentService) CreatePostgresDatabase(ctx context.Context, opts client.CreatePostgresDatabaseOpts) (*client.PostgresDatabase, error) {
	err := c.spinner.Start("postgres")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = c.spinner.Stop()
	}()

	rotaterBucket, err := c.api.CreateS3Bucket(api.CreateS3BucketOpts{
		ID:        opts.ID,
		Name:      rotaterBucketName(opts.ID.ClusterName, opts.ApplicationName),
		StackName: cfn.NewStackNamer().S3Bucket(opts.ApplicationName, opts.ID.Repository, opts.ID.Environment),
	})
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	// Create a zip archive of the lambda code
	w := zip.NewWriter(buf)
	f, err := w.Create(postgresRotaterLambdaKey)
	if err != nil {
		return nil, err
	}

	_, err = f.Write([]byte(static.SecretsManagerPostgresRotationSingleUserBody))
	if err != nil {
		return nil, err
	}

	err = s3api.New(c.provider).PutObject(
		rotaterBucket.Name,
		postgresRotaterLambdaKey,
		bytes.NewReader(buf.Bytes()),
	)
	if err != nil {
		return nil, err
	}

	pg, err := c.api.CreatePostgresDatabase(api.CreatePostgresDatabaseOpts{
		ID:                opts.ID,
		ApplicationName:   opts.ApplicationName,
		UserName:          opts.UserName,
		StackName:         cfn.NewStackNamer().RDSPostgres(opts.ApplicationName, opts.ID.Repository, opts.ID.Environment),
		VpcID:             opts.VpcID,
		DBSubnetGroupName: opts.DBSubnetGroupName,
		DBSubnetIDs:       opts.DBSubnetIDs,
		DBSubnetCIDRs:     opts.DBSubnetCIDRs,
		RotaterBucket:     rotaterBucket.Name,
		RotaterKey:        postgresRotaterLambdaKey,
	})
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
		ID: opts.ID,
		Manifests: []api.Manifest{
			{
				Name:      adminSecretName(opts.ApplicationName),
				Namespace: opts.Namespace,
				Backend:   api.BackendTypeSecretsManager,
				Data: []api.Data{
					{
						Key:      pg.AdminSecretFriendlyName,
						Name:     "username",
						Property: "username",
					},
					{
						Key:      pg.AdminSecretFriendlyName,
						Name:     "password",
						Property: "password",
					},
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
			"endpointAddress": pg.EndpointAddress,
			"endpointPort":    strconv.Itoa(pg.EndpointPort),
		},
	})
	if err != nil {
		return nil, err
	}

	postgres := &client.PostgresDatabase{
		Namespace:             opts.Namespace,
		PostgresDatabase:      pg,
		AdminSecretName:       adminSecretName(opts.ApplicationName),
		DatabaseConfigMapName: pgConfigMapName(opts.ApplicationName),
		RotaterBucket:         rotaterBucket,
	}

	r1, err := c.store.SavePostgresDatabase(postgres)
	if err != nil {
		return nil, err
	}

	r2, err := c.state.SavePostgresDatabase(postgres)
	if err != nil {
		return nil, err
	}

	err = c.report.ReportCreatePostgresDatabase(postgres, []*store.Report{r1, r2})
	if err != nil {
		return nil, err
	}

	return postgres, nil
}

// nolint: funlen
func (c *componentService) DeletePostgresDatabase(ctx context.Context, opts client.DeletePostgresDatabaseOpts) error {
	err := c.spinner.Start("postgres")
	if err != nil {
		return err
	}

	defer func() {
		_ = c.spinner.Stop()
	}()

	err = c.api.DeletePostgresDatabase(api.DeletePostgresDatabaseOpts{
		ID:        opts.ID,
		StackName: cfn.NewStackNamer().RDSPostgres(opts.ApplicationName, opts.ID.Repository, opts.ID.Environment),
	})
	if err != nil {
		return err
	}

	err = c.manifest.DeleteConfigMap(ctx, client.DeleteConfigMapOpts{
		ID:        opts.ID,
		Name:      pgConfigMapName(opts.ApplicationName),
		Namespace: opts.Namespace,
	})
	if err != nil {
		return err
	}

	err = c.manifest.DeleteExternalSecret(ctx, client.DeleteExternalSecretOpts{
		ID: opts.ID,
		Secrets: map[string]string{
			pgConfigMapName(opts.ApplicationName): opts.Namespace,
		},
	})
	if err != nil {
		return err
	}

	err = s3api.New(c.provider).DeleteObject(
		rotaterBucketName(opts.ID.ClusterName, opts.ApplicationName),
		postgresRotaterLambdaKey,
	)
	if err != nil {
		return err
	}

	err = c.api.DeleteS3Bucket(api.DeleteS3BucketOpts{
		ID:        opts.ID,
		StackName: cfn.NewStackNamer().S3Bucket(opts.ApplicationName, opts.ID.Repository, opts.ID.Environment),
	})
	if err != nil {
		return err
	}

	r1, err := c.store.RemovePostgresDatabase(opts.ApplicationName)
	if err != nil {
		return err
	}

	r2, err := c.state.RemovePostgresDatabase(opts.ApplicationName)
	if err != nil {
		return err
	}

	err = c.report.ReportDeletePostgresDatabase(opts.ApplicationName, []*store.Report{r1, r2})
	if err != nil {
		return err
	}

	return nil
}

// NewComponentService returns an initialised component service
func NewComponentService(
	spin spinner.Spinner,
	api client.ComponentAPI,
	store client.ComponentStore,
	state client.ComponentState,
	report client.ComponentReport,
	manifest client.ManifestService,
	provider v1alpha1.CloudProvider,
) client.ComponentService {
	return &componentService{
		spinner:  spin,
		api:      api,
		store:    store,
		state:    state,
		report:   report,
		manifest: manifest,
		provider: provider,
	}
}
