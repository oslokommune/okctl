package rest

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// SecurityGroupsTarget defines the REST API Route
const (
	SecurityGroupsTarget      = "securitygroups/"
	SecurityGroupsRulesTarget = SecurityGroupsTarget + "rules/"
)

type securityGroupAPI struct {
	client *HTTPClient
}

func (s *securityGroupAPI) CreateSecurityGroup(_ context.Context, opts api.CreateSecurityGroupOpts) (api.SecurityGroup, error) {
	into := api.SecurityGroup{}
	return into, s.client.DoPost(SecurityGroupsTarget, &opts, &into)
}

func (s *securityGroupAPI) GetSecurityGroup(_ context.Context, opts api.GetSecurityGroupOpts) (api.SecurityGroup, error) {
	into := api.SecurityGroup{}
	return into, s.client.DoGet(SecurityGroupsTarget, &opts, &into)
}

func (s *securityGroupAPI) DeleteSecurityGroup(_ context.Context, opts api.DeleteSecurityGroupOpts) error {
	return s.client.DoDelete(SecurityGroupsTarget, &opts)
}

func (s *securityGroupAPI) AddRule(_ context.Context, opts api.AddRuleOpts) (api.Rule, error) {
	into := api.Rule{}
	return into, s.client.DoPost(SecurityGroupsRulesTarget, &opts, &into)
}

func (s *securityGroupAPI) RemoveRule(_ context.Context, opts api.RemoveRuleOpts) error {
	return s.client.DoDelete(SecurityGroupsRulesTarget, &opts)
}

// NewSecurityGroupAPI returns an initialised API client
func NewSecurityGroupAPI(client *HTTPClient) client.SecurityGroupAPI {
	return &securityGroupAPI{
		client: client,
	}
}
