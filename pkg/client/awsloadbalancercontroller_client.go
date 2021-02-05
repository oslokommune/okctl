package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// AWSLoadBalancerController contains the state of an
// alb ingress controller deployment
type AWSLoadBalancerController struct {
	Policy         *api.ManagedPolicy
	ServiceAccount *api.ServiceAccount
	Chart          *api.Helm
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

// AWSLoadBalancerControllerAPI defines the API layer
type AWSLoadBalancerControllerAPI interface {
	CreateAWSLoadBalancerControllerPolicy(opts api.CreateAWSLoadBalancerControllerPolicyOpts) (*api.ManagedPolicy, error)
	DeleteAWSLoadBalancerControllerPolicy(id api.ID) error
	CreateAWSLoadBalancerControllerServiceAccount(opts api.CreateAWSLoadBalancerControllerServiceAccountOpts) (*api.ServiceAccount, error)
	DeleteAWSLoadBalancerControllerServiceAccount(id api.ID) error
	CreateAWSLoadBalancerControllerHelmChart(opts api.CreateAWSLoadBalancerControllerHelmChartOpts) (*api.Helm, error)
}

// AWSLoadBalancerControllerStore defines the storage layer
type AWSLoadBalancerControllerStore interface {
	SaveAWSLoadBalancerController(controller *AWSLoadBalancerController) (*store.Report, error)
	RemoveAWSLoadBalancerController(id api.ID) (*store.Report, error)
}

// AWSLoadBalancerControllerReport defines the report layer
type AWSLoadBalancerControllerReport interface {
	ReportCreateAWSLoadBalancerController(controller *AWSLoadBalancerController, report *store.Report) error
	ReportDeleteAWSLoadBalancerController(report *store.Report) error
}
