package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
)

// AWSLoadBalancerController contains the state of an
// alb ingress controller deployment
type AWSLoadBalancerController struct {
	Policy         *ManagedPolicy
	ServiceAccount *ServiceAccount
	Chart          *Helm
}

// CreateAWSLoadBalancerControllerOpts defines the required inputs
type CreateAWSLoadBalancerControllerOpts struct {
	ID    api.ID
	VPCID string
}

// AWSLoadBalancerControllerService defines the service layer
type AWSLoadBalancerControllerService interface {
	CreateAWSLoadBalancerController(ctx context.Context, opts CreateAWSLoadBalancerControllerOpts) (*AWSLoadBalancerController, error)
	DeleteAWSLoadBalancerController(ctx context.Context, id api.ID) error
}
