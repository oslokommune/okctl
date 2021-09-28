package core

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/oslokommune/okctl/pkg/api"
)

func makeCreateSecurityGroupEndpoint(s api.SecurityGroupService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		into := request.(*api.CreateSecurityGroupOpts)

		return s.CreateSecurityGroup(ctx, *into)
	}
}

func makeGetSecurityGroupEndpoint(s api.SecurityGroupService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		into := request.(*api.GetSecurityGroupOpts)

		return s.GetSecurityGroup(ctx, *into)
	}
}

func makeDeleteSecurityGroupEndpoint(s api.SecurityGroupService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		into := request.(*api.DeleteSecurityGroupOpts)

		return &Empty{}, s.DeleteSecurityGroup(ctx, *into)
	}
}

func makeAddSecurityGroupRuleEndpoint(s api.SecurityGroupService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		into := request.(*api.AddRuleOpts)

		return s.AddRule(ctx, *into)
	}
}

func makeRemoveSecurityGroupRuleEndpoint(s api.SecurityGroupService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		into := request.(*api.RemoveRuleOpts)

		return &Empty{}, s.RemoveRule(ctx, *into)
	}
}
