package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

const (
	// TargetAutoscalerPolicy is the API route for the policy
	TargetAutoscalerPolicy = "managedpolicies/autoscaler/"
	// TargetAutoscalerServiceAccount is the API route for the service account
	TargetAutoscalerServiceAccount = "serviceaccounts/autoscaler/"
	// TargetAutoscalerHelm is the API route for helm
	TargetAutoscalerHelm = "helm/autoscaler/"
)

type autoscalerAPI struct {
	client *HTTPClient
}

func (a *autoscalerAPI) DeleteAutoscalerPolicy(id api.ID) error {
	return a.client.DoDelete(TargetAutoscalerPolicy, &id)
}

func (a *autoscalerAPI) DeleteAutoscalerServiceAccount(id api.ID) error {
	return a.client.DoDelete(TargetAutoscalerServiceAccount, &id)
}

func (a *autoscalerAPI) CreateAutoscalerPolicy(opts api.CreateAutoscalerPolicy) (*api.ManagedPolicy, error) {
	into := &api.ManagedPolicy{}
	return into, a.client.DoPost(TargetAutoscalerPolicy, &opts, into)
}

func (a *autoscalerAPI) CreateAutoscalerServiceAccount(opts api.CreateAutoscalerServiceAccountOpts) (*api.ServiceAccount, error) {
	into := &api.ServiceAccount{}
	return into, a.client.DoPost(TargetAutoscalerServiceAccount, &opts, into)
}

func (a *autoscalerAPI) CreateAutoscalerHelmChart(opts api.CreateAutoscalerHelmChartOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, a.client.DoPost(TargetAutoscalerHelm, &opts, into)
}

// NewAutoscalerAPI returns an initialised API client
func NewAutoscalerAPI(client *HTTPClient) client.AutoscalerAPI {
	return &autoscalerAPI{
		client: client,
	}
}
