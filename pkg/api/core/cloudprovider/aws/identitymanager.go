package aws

import (
	"fmt"

	"github.com/gosimple/slug"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
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

// nolint: funlen
func (s *identityManagerCloudProvider) CreateIdentityPool(certificateARN string, opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error) {
	clients := make([]components.UserPoolClient, len(opts.Clients))

	for i, c := range opts.Clients {
		clients[i] = components.UserPoolClient{
			Purpose:     c.Purpose,
			CallbackURL: c.CallbackURL,
		}
	}

	b := cfn.New(components.NewUserPoolWithClients(
		opts.ID.Environment, opts.ID.Repository, opts.AuthDomain, certificateARN, clients),
	)

	stackName := cfn.NewStackNamer().IdentityPool(opts.ID.Repository, opts.ID.Environment)

	template, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf("building identity pool cloud formation template: %w", err)
	}

	r := cfn.NewRunner(s.provider)

	err = r.CreateIfNotExists(stackName, template, nil, defaultTimeOut)
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

	aliasStackName := cfn.NewStackNamer().AliasRecordSet(opts.ID.Repository, opts.ID.Environment, slug.Make(d.UserPoolDomain))

	err = r.CreateIfNotExists(aliasStackName, aliasTemplate, nil, defaultTimeOut)
	if err != nil {
		return nil, fmt.Errorf("creating alias cloud formation stack: %w", err)
	}

	pool := &api.IdentityPool{
		ID:                      opts.ID,
		AuthDomain:              opts.AuthDomain,
		HostedZoneID:            opts.HostedZoneID,
		StackName:               stackName,
		CloudFormationTemplates: template,
		Clients:                 nil,
		Certificate:             nil,
		RecordSetAlias: &api.RecordSetAlias{
			AliasDomain:            d.CloudFrontDomainName,
			AliasHostedZones:       DefaultCloudFrontHostedZoneID,
			StackName:              aliasStackName,
			CloudFormationTemplate: aliasTemplate,
		},
	}

	outputs := make(map[string]cfn.ProcessOutputFn, len(clients))
	apiClients := make([]*api.IdentityClient, len(clients))

	for i, c := range opts.Clients {
		apiClients[i] = &api.IdentityClient{
			Purpose:     c.Purpose,
			CallbackURL: c.CallbackURL,
		}

		outputs[fmt.Sprintf("%sClientID", c.Purpose)] = cfn.String(&apiClients[i].ClientID)
	}

	outputs["UserPool"] = cfn.String(&pool.UserPoolID)

	err = r.Outputs(stackName, outputs)
	if err != nil {
		return nil, fmt.Errorf("retrieving identity pool outputs: %w", err)
	}

	for _, c := range apiClients {
		secret, err := cognito.New(s.provider).UserPoolClientSecret(c.ClientID, pool.UserPoolID)
		if err != nil {
			return nil, fmt.Errorf("retrieving client secret: %w", err)
		}

		c.ClientSecret = secret
	}

	pool.Clients = apiClients

	return pool, nil
}

// NewIdentityManagerCloudProvider returns an initialised cloud layer
func NewIdentityManagerCloudProvider(provider v1alpha1.CloudProvider) api.IdentityManagerCloudProvider {
	return &identityManagerCloudProvider{
		provider: provider,
	}
}
