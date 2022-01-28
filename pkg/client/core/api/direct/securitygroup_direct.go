package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type securityGroupAPIDirectClient struct {
	service api.SecurityGroupService
}

func (s *securityGroupAPIDirectClient) CreateSecurityGroup(_ context.Context, opts api.CreateSecurityGroupOpts) (api.SecurityGroup, error) {
	return s.service.CreateSecurityGroup(context.Background(), opts)
}

func (s *securityGroupAPIDirectClient) GetSecurityGroup(_ context.Context, opts api.GetSecurityGroupOpts) (api.SecurityGroup, error) {
	return s.service.GetSecurityGroup(context.Background(), opts)
}

func (s *securityGroupAPIDirectClient) DeleteSecurityGroup(_ context.Context, opts api.DeleteSecurityGroupOpts) error {
	return s.service.DeleteSecurityGroup(context.Background(), opts)
}

func (s *securityGroupAPIDirectClient) AddRule(_ context.Context, opts api.AddRuleOpts) (api.Rule, error) {
	return s.service.AddRule(context.Background(), opts)
}

func (s *securityGroupAPIDirectClient) RemoveRule(_ context.Context, opts api.RemoveRuleOpts) error {
	return s.service.RemoveRule(context.Background(), opts)
}

// NewSecurityGroupAPI returns an initialised API client
func NewSecurityGroupAPI(service api.SecurityGroupService) client.SecurityGroupAPI {
	return &securityGroupAPIDirectClient{
		service: service,
	}
}
