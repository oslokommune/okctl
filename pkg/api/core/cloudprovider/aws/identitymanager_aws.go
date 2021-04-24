package aws

import (
	"fmt"

	"github.com/gosimple/slug"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
	"github.com/oslokommune/okctl/pkg/cognito"
)

const (
	// DefaultCloudFrontHostedZoneID is the default hosted zone id for all cloud front distributions
	DefaultCloudFrontHostedZoneID = "Z2FDTNDATAQYW2"
	// DefaultCloudFrontACMRegion is the default region for ACM certificates used with a cloud front distribution
	DefaultCloudFrontACMRegion = "us-east-1"
)

type identityManagerCloudProvider struct {
	provider v1alpha1.CloudProvider
}

func (s *identityManagerCloudProvider) DeleteIdentityPoolClient(opts api.DeleteIdentityPoolClientOpts) error {
	err := cfn.NewRunner(s.provider).Delete(cfn.NewStackNamer().IdentityPoolClient(opts.ID.ClusterName, opts.Purpose))
	if err != nil {
		return fmt.Errorf("deleting identity pool client: %w", err)
	}

	return nil
}

func (s *identityManagerCloudProvider) CreateIdentityPoolClient(opts api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error) {
	b := cfn.New(components.NewUserPoolClient(
		opts.Purpose,
		opts.ID.ClusterName,
		opts.CallbackURL,
		opts.UserPoolID,
	))

	stackName := cfn.NewStackNamer().IdentityPoolClient(opts.ID.ClusterName, opts.Purpose)

	template, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf("building identity pool client template: %w", err)
	}

	r := cfn.NewRunner(s.provider)

	err = r.CreateIfNotExists(opts.ID.ClusterName, stackName, template, nil, defaultTimeOut)
	if err != nil {
		return nil, fmt.Errorf("creating identity pool client cloud formation stack: %w", err)
	}

	c := &api.IdentityPoolClient{
		ID:                      opts.ID,
		UserPoolID:              opts.UserPoolID,
		Purpose:                 opts.Purpose,
		CallbackURL:             opts.CallbackURL,
		StackName:               stackName,
		CloudFormationTemplates: template,
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		fmt.Sprintf("%sClientID", opts.Purpose): cfn.String(&c.ClientID),
	})
	if err != nil {
		return nil, fmt.Errorf("retrieving identity pool client outputs: %w", err)
	}

	secret, err := cognito.New(s.provider).UserPoolClientSecret(c.ClientID, opts.UserPoolID)
	if err != nil {
		return nil, fmt.Errorf("retrieving client secret: %w", err)
	}

	c.ClientSecret = secret

	return c, nil
}

func (s *identityManagerCloudProvider) DeleteIdentityPool(opts api.DeleteIdentityPoolOpts) error {
	r := cfn.NewRunner(s.provider)

	err := r.Delete(cfn.NewStackNamer().AliasRecordSet(opts.ID.ClusterName, slug.Make(opts.Domain)))
	if err != nil {
		return fmt.Errorf("deleting alias record set for identity pool: %w", err)
	}

	err = r.Delete(cfn.NewStackNamer().IdentityPool(opts.ID.ClusterName))
	if err != nil {
		return fmt.Errorf("deleting identity pool: %w", err)
	}

	return nil
}

// nolint: funlen
func (s *identityManagerCloudProvider) CreateIdentityPool(certificateARN string, opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error) {
	b := cfn.New(components.NewUserPool(
		opts.ID.ClusterName,
		opts.AuthDomain,
		opts.HostedZoneID,
		certificateARN,
	),
	)

	stackName := cfn.NewStackNamer().IdentityPool(opts.ID.ClusterName)

	template, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf("building identity pool cloud formation template: %w", err)
	}

	r := cfn.NewRunner(s.provider)

	err = r.CreateIfNotExists(opts.ID.ClusterName, stackName, template, nil, defaultTimeOut)
	if err != nil {
		return nil, fmt.Errorf("creating identity pool cloud formation stack: %w", err)
	}

	d, err := cognito.New(s.provider).UserPoolDomainInfo(opts.AuthDomain)
	if err != nil {
		return nil, fmt.Errorf("getting cloudfront auth domain info: %w", err)
	}

	ba := cfn.New(components.NewAliasRecordSet("Auth", d.CloudFrontDomainName, DefaultCloudFrontHostedZoneID, d.UserPoolDomain, opts.HostedZoneID))

	aliasTemplate, err := ba.Build()
	if err != nil {
		return nil, fmt.Errorf("building alias cloud formation template: %w", err)
	}

	aliasStackName := cfn.NewStackNamer().AliasRecordSet(opts.ID.ClusterName, slug.Make(d.UserPoolDomain))

	err = r.CreateIfNotExists(opts.ID.ClusterName, aliasStackName, aliasTemplate, nil, defaultTimeOut)
	if err != nil {
		return nil, fmt.Errorf("creating alias cloud formation stack: %w", err)
	}

	pool := &api.IdentityPool{
		ID:                      opts.ID,
		AuthDomain:              opts.AuthDomain,
		HostedZoneID:            opts.HostedZoneID,
		StackName:               stackName,
		CloudFormationTemplates: template,
		RecordSetAlias: &api.RecordSetAlias{
			AliasDomain:            d.CloudFrontDomainName,
			AliasHostedZones:       DefaultCloudFrontHostedZoneID,
			StackName:              aliasStackName,
			CloudFormationTemplate: aliasTemplate,
		},
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		"UserPool": cfn.String(&pool.UserPoolID),
	})
	if err != nil {
		return nil, fmt.Errorf("retrieving identity pool outputs: %w", err)
	}

	// Skipping this for now, since we need to support the flow differently
	// than we are doing today.
	// err = cognito.New(s.provider).EnableMFA(pool.UserPoolID)
	// if err != nil {
	// 	return nil, fmt.Errorf("enabling mfa on the user pool: %w", err)
	// }

	return pool, nil
}

func (s *identityManagerCloudProvider) CreateIdentityPoolUser(opts api.CreateIdentityPoolUserOpts) (*api.IdentityPoolUser, error) {
	b := cfn.New(components.NewUserPoolUser(
		opts.Email,
		opts.UserPoolID,
	),
	)
	stackName := cfn.NewStackNamer().IdentityPoolUser(opts.ID.ClusterName, slug.Make(opts.Email))

	template, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf("building identity pool user cloud formation template: %w", err)
	}

	r := cfn.NewRunner(s.provider)

	err = r.CreateIfNotExists(opts.ID.ClusterName, stackName, template, nil, defaultTimeOut)
	if err != nil {
		return nil, fmt.Errorf("creating identity pool user cloud formation stack: %w", err)
	}

	user := &api.IdentityPoolUser{
		ID:                     opts.ID,
		Email:                  opts.Email,
		UserPoolID:             opts.UserPoolID,
		StackName:              stackName,
		CloudFormationTemplate: template,
	}

	return user, nil
}

// NewIdentityManagerCloudProvider returns an initialised cloud layer
func NewIdentityManagerCloudProvider(provider v1alpha1.CloudProvider) api.IdentityManagerCloudProvider {
	return &identityManagerCloudProvider{
		provider: provider,
	}
}
