package core

import (
	"context"
	"fmt"

	"github.com/gosimple/slug"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type identityManagerService struct {
	api   client.IdentityManagerAPI
	state client.IdentityManagerState

	cert client.CertificateService
}

func (s *identityManagerService) DeleteIdentityPoolClient(_ context.Context, opts client.DeleteIdentityPoolClientOpts) error {
	err := s.api.DeleteIdentityPoolClient(api.DeleteIdentityPoolClientOpts{
		ID:      opts.ID,
		Purpose: opts.Purpose,
	})
	if err != nil {
		return err
	}

	stackName := cfn.NewStackNamer().IdentityPoolClient(opts.ID.ClusterName, opts.Purpose)

	err = s.state.RemoveIdentityPoolClient(stackName)
	if err != nil {
		return err
	}

	return nil
}

// DeleteIdentityPool and all users
func (s *identityManagerService) DeleteIdentityPool(ctx context.Context, id api.ID) error {
	stackName := cfn.NewStackNamer().IdentityPool(id.ClusterName)

	pool, err := s.state.GetIdentityPool(stackName)
	if err != nil {
		return err
	}

	err = s.cert.DeleteCognitoCertificate(ctx, client.DeleteCognitoCertificateOpts{
		ID:     id,
		Domain: pool.AuthDomain,
	})
	if err != nil {
		return err
	}

	err = s.api.DeleteIdentityPool(api.DeleteIdentityPoolOpts{
		ID:         id,
		UserPoolID: pool.UserPoolID,
		Domain:     pool.AuthDomain,
	})
	if err != nil {
		return err
	}

	err = s.state.RemoveIdentityPool(stackName)
	if err != nil {
		return err
	}

	return nil
}

func (s *identityManagerService) CreateIdentityPoolUser(_ context.Context, opts client.CreateIdentityPoolUserOpts) (*client.IdentityPoolUser, error) {
	err := opts.Validate()
	if err != nil {
		return nil, err
	}

	u, err := s.api.CreateIdentityPoolUser(api.CreateIdentityPoolUserOpts{
		ID:         opts.ID,
		Email:      opts.Email,
		UserPoolID: opts.UserPoolID,
	})
	if err != nil {
		return nil, fmt.Errorf("creating identity pool user: %w", err)
	}

	user := &client.IdentityPoolUser{
		ID:                     u.ID,
		Email:                  u.Email,
		UserPoolID:             u.UserPoolID,
		StackName:              u.StackName,
		CloudFormationTemplate: u.CloudFormationTemplate,
	}

	err = s.state.SaveIdentityPoolUser(user)
	if err != nil {
		return nil, fmt.Errorf("updating identity pool user state: %w", err)
	}

	return user, nil
}

func (s *identityManagerService) DeleteIdentityPoolUser(_ context.Context, opts client.DeleteIdentityPoolUserOpts) error {
	err := opts.Validate()
	if err != nil {
		return err
	}

	err = s.api.DeleteIdentityPoolUser(api.DeleteIdentityPoolUserOpts{
		ClusterID: opts.ClusterID,
		UserEmail: opts.UserEmail,
	})
	if err != nil {
		return err
	}

	err = s.state.RemoveIdentityPoolUser(cfn.NewStackNamer().IdentityPoolUser(
		opts.ClusterID.ClusterName,
		slug.Make(opts.UserEmail),
	))
	if err != nil {
		return fmt.Errorf("reomving identity pool user from state: %w", err)
	}

	return nil
}

func (s *identityManagerService) CreateIdentityPoolClient(_ context.Context, opts client.CreateIdentityPoolClientOpts) (*client.IdentityPoolClient, error) {
	c, err := s.api.CreateIdentityPoolClient(api.CreateIdentityPoolClientOpts{
		ID:          opts.ID,
		UserPoolID:  opts.UserPoolID,
		Purpose:     opts.Purpose,
		CallbackURL: opts.CallbackURL,
	})
	if err != nil {
		return nil, fmt.Errorf("creating identity pool client: %w", err)
	}

	cl := &client.IdentityPoolClient{
		ID:                      c.ID,
		UserPoolID:              c.UserPoolID,
		Purpose:                 c.Purpose,
		CallbackURL:             c.CallbackURL,
		ClientID:                c.ClientID,
		ClientSecret:            c.ClientSecret,
		StackName:               c.StackName,
		CloudFormationTemplates: c.CloudFormationTemplates,
	}

	err = s.state.SaveIdentityPoolClient(cl)
	if err != nil {
		return nil, fmt.Errorf("storing identity pool client state: %w", err)
	}

	return cl, nil
}

func (s *identityManagerService) CreateIdentityPool(_ context.Context, opts client.CreateIdentityPoolOpts) (*client.IdentityPool, error) {
	p, err := s.api.CreateIdentityPool(api.CreateIdentityPoolOpts{
		ID:           opts.ID,
		AuthDomain:   opts.AuthDomain,
		AuthFQDN:     opts.AuthFQDN,
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return nil, fmt.Errorf("creating identity pool: %w", err)
	}

	pool := &client.IdentityPool{
		ID:                      p.ID,
		UserPoolID:              p.UserPoolID,
		AuthDomain:              p.AuthDomain,
		HostedZoneID:            p.HostedZoneID,
		StackName:               p.StackName,
		CloudFormationTemplates: p.CloudFormationTemplates,
		Certificate: &client.Certificate{
			ID:                     p.Certificate.ID,
			FQDN:                   p.Certificate.FQDN,
			Domain:                 p.Certificate.Domain,
			HostedZoneID:           p.Certificate.HostedZoneID,
			ARN:                    p.Certificate.CertificateARN,
			StackName:              p.Certificate.StackName,
			CloudFormationTemplate: p.Certificate.CloudFormationTemplate,
		},
		RecordSetAlias: &client.RecordSetAlias{
			AliasDomain:            p.RecordSetAlias.AliasDomain,
			AliasHostedZones:       p.RecordSetAlias.AliasHostedZones,
			StackName:              p.RecordSetAlias.StackName,
			CloudFormationTemplate: p.RecordSetAlias.CloudFormationTemplate,
		},
	}

	err = s.state.SaveIdentityPool(pool)
	if err != nil {
		return nil, fmt.Errorf("saving identity pool state: %w", err)
	}

	return pool, nil
}

// NewIdentityManagerService returns an initialised service
func NewIdentityManagerService(
	api client.IdentityManagerAPI,
	state client.IdentityManagerState,
	cert client.CertificateService,
) client.IdentityManagerService {
	return &identityManagerService{
		api:   api,
		state: state,
		cert:  cert,
	}
}
