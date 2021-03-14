package aws

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
)

type componentCloudProvider struct {
	provider v1alpha1.CloudProvider
}

const (
	postgresTimeOutInMinutes = 45
)

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
	})

	b := cfn.New(composer)

	template, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf("building cloud formation template: %w", err)
	}

	r := cfn.NewRunner(c.provider)

	err = r.CreateIfNotExists(
		opts.StackName,
		template,
		[]string{cfn.CapabilityNamedIam},
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
		CloudFormationTemplate:  template,
	}

	err = r.Outputs(opts.StackName, map[string]cfn.ProcessOutputFn{
		fmt.Sprintf("%sEndpointAddress", composer.NameResource("RDSPostgres")): cfn.String(&p.EndpointAddress),
		fmt.Sprintf("%sEndpointPort", composer.NameResource("RDSPostgres")):    cfn.Int(&p.EndpointPort),
		fmt.Sprintf("%sGroupId", composer.NameResource("RDSPostgresOutgoing")): cfn.String(&p.OutgoingSecurityGroupID),
		composer.NameResource("RDSInstanceAdmin"):                              cfn.String(&p.SecretsManagerAdminSecretARN),
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
