package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

const (
	// TargetAWSLoadBalancerControllerPolicy is the REST API route
	TargetAWSLoadBalancerControllerPolicy = "managedpolicies/awsloadbalancercontroller/"
	// TargetAWSLoadBalancerControllerServiceAccount is the REST API route
	TargetAWSLoadBalancerControllerServiceAccount = "serviceaccounts/awsloadbalancercontroller/"
	// TargetAWSLoadBalancerControllerHelm is the REST API route
	TargetAWSLoadBalancerControllerHelm = "helm/awsloadbalancercontroller/"
)

type awsLoadBalancerControllerAPI struct {
	client *HTTPClient
}

func (a *awsLoadBalancerControllerAPI) DeleteAWSLoadBalancerControllerPolicy(id api.ID) error {
	return a.client.DoDelete(TargetAWSLoadBalancerControllerPolicy, &id)
}

func (a *awsLoadBalancerControllerAPI) DeleteAWSLoadBalancerControllerServiceAccount(id api.ID) error {
	return a.client.DoDelete(TargetAWSLoadBalancerControllerServiceAccount, &id)
}

func (a *awsLoadBalancerControllerAPI) CreateAWSLoadBalancerControllerPolicy(opts api.CreateAWSLoadBalancerControllerPolicyOpts) (*api.ManagedPolicy, error) {
	into := &api.ManagedPolicy{}
	return into, a.client.DoPost(TargetAWSLoadBalancerControllerPolicy, &opts, into)
}

// nolint: lll
func (a *awsLoadBalancerControllerAPI) CreateAWSLoadBalancerControllerServiceAccount(opts api.CreateAWSLoadBalancerControllerServiceAccountOpts) (*api.ServiceAccount, error) {
	into := &api.ServiceAccount{}
	return into, a.client.DoPost(TargetAWSLoadBalancerControllerServiceAccount, &opts, into)
}

func (a *awsLoadBalancerControllerAPI) CreateAWSLoadBalancerControllerHelmChart(opts api.CreateAWSLoadBalancerControllerHelmChartOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, a.client.DoPost(TargetAWSLoadBalancerControllerHelm, &opts, into)
}

// NewAWSLoadBalancerControllerAPI returns an initialised REST API client
func NewAWSLoadBalancerControllerAPI(client *HTTPClient) client.AWSLoadBalancerControllerAPI {
	return &awsLoadBalancerControllerAPI{
		client: client,
	}
}
