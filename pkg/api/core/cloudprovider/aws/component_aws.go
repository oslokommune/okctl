package aws

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn/components/lambdafunction"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
)

type componentCloudProvider struct {
	provider v1alpha1.CloudProvider
}

func (c *componentCloudProvider) CreateS3Bucket(opts *api.CreateS3BucketOpts) (*api.S3Bucket, error) {
	composition := components.NewS3BucketComposer(opts.Name, opts.ID.Repository, opts.ID.Environment)

	template, err := cfn.New(composition).Build()
	if err != nil {
		return nil, fmt.Errorf("building the cloud formation template: %w", err)
	}

	r := cfn.NewRunner(c.provider)

	err = r.CreateIfNotExists(
		opts.StackName,
		template,
		nil,
		defaultTimeOut,
	)
	if err != nil {
		return nil, fmt.Errorf("creating cloud formation stack: %w", err)
	}

	b := &api.S3Bucket{
		ID:                     opts.ID,
		StackName:              opts.StackName,
		CloudFormationTemplate: string(template),
	}

	err = r.Outputs(opts.StackName, map[string]cfn.ProcessOutputFn{
		composition.ResourceBucketNameOutput(): cfn.String(&b.Name),
	})
	if err != nil {
		return nil, fmt.Errorf("collecting stack outputs: %w", err)
	}

	return b, nil
}

func (c *componentCloudProvider) DeleteS3Bucket(opts *api.DeleteS3BucketOpts) error {
	return cfn.NewRunner(c.provider).Delete(opts.StackName)
}

const (
	postgresTimeOutInMinutes = 45
)

// nolint: funlen
func (c *componentCloudProvider) CreatePostgresDatabase(opts *api.CreatePostgresDatabaseOpts) (*api.PostgresDatabase, error) {
	composer := components.NewRDSPostgresComposer(components.RDSPostgresComposerOpts{
		ApplicationDBName: opts.ApplicationName,
		AWSAccountID:      opts.ID.AWSAccountID,
		Repository:        opts.ID.Repository,
		Environment:       opts.ID.Environment,
		DBSubnetGroupName: opts.DBSubnetGroupName,
		UserName:          opts.UserName,
		VpcID:             opts.VpcID,
		VPCDBSubnetIDs:    opts.DBSubnetIDs,
		VPCDBSubnetCIDRs:  opts.DBSubnetCIDRs,
		Bucket:            opts.RotaterBucket,
		Key:               opts.RotaterKey,
	})

	b := cfn.New(composer)

	template, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf("building cloud formation template: %w", err)
	}

	template, err = lambdafunction.PatchRotateLambda(
		composer.NameResource("RDSPostgresLambdaRotateFunction"),
		composer.NameResource("SecretsManagerVPCEndpoint"),
		template,
	)
	if err != nil {
		return nil, fmt.Errorf("patching cloud formation template: %w", err)
	}

	r := cfn.NewRunner(c.provider)

	err = r.CreateIfNotExists(
		opts.StackName,
		template,
		[]string{cfn.CapabilityNamedIam, cfn.CapabilityAutoExpand},
		postgresTimeOutInMinutes,
	)
	if err != nil {
		return nil, fmt.Errorf("creating cloud formation stack: %w", err)
	}

	p := &api.PostgresDatabase{
		ID:                      opts.ID,
		ApplicationName:         opts.ApplicationName,
		UserName:                opts.UserName,
		StackName:               opts.StackName,
		AdminSecretFriendlyName: composer.AdminSecretFriendlyName(),
		CloudFormationTemplate:  string(template),
	}

	err = r.Outputs(opts.StackName, map[string]cfn.ProcessOutputFn{
		fmt.Sprintf("%sEndpointAddress", composer.NameResource("RDSPostgres")):     cfn.String(&p.EndpointAddress),
		fmt.Sprintf("%sEndpointPort", composer.NameResource("RDSPostgres")):        cfn.Int(&p.EndpointPort),
		fmt.Sprintf("%sGroupId", composer.NameResource("RDSPostgresOutgoing")):     cfn.String(&p.OutgoingSecurityGroupID),
		composer.NameResource("RDSInstanceAdmin"):                                  cfn.String(&p.SecretsManagerAdminSecretARN),
		composer.NameResource("RDSPostgresLambdaManagedPolicy"):                    cfn.String(&p.LambdaPolicyARN),
		fmt.Sprintf("%sArn", composer.NameResource("RDSPostgresLambdaRotateRole")): cfn.String(&p.LambdaRoleARN),
		composer.NameResource("RDSPostgresLambdaRotateFunction"):                   cfn.String(&p.LambdaFunctionARN),
	})
	if err != nil {
		return nil, fmt.Errorf("collecting stack outputs: %w", err)
	}

	return p, nil
}

func (c *componentCloudProvider) DeletePostgresDatabase(opts *api.DeletePostgresDatabaseOpts) error {
	return cfn.NewRunner(c.provider).Delete(opts.StackName)
}

// NewComponentCloudProvider returns an initialised component cloud provider
func NewComponentCloudProvider(provider v1alpha1.CloudProvider) api.ComponentCloudProvider {
	return &componentCloudProvider{
		provider: provider,
	}
}
