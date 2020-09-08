package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

const (
	// TargetAlbIngressControllerPolicy is the REST API route
	TargetAlbIngressControllerPolicy = "managedpolicies/albingresscontroller/"
	// TargetAlbIngressControllerServiceAccount is the REST API route
	TargetAlbIngressControllerServiceAccount = "serviceaccounts/albingresscontroller/"
	// TargetAlbIngressControllerHelm is the REST API route
	TargetAlbIngressControllerHelm = "helm/albingresscontroller/"
)

type albIngressControllerAPI struct {
	client *HTTPClient
}

func (a *albIngressControllerAPI) CreateAlbIngressControllerPolicy(opts api.CreateAlbIngressControllerPolicyOpts) (*api.ManagedPolicy, error) {
	into := &api.ManagedPolicy{}
	return into, a.client.DoPost(TargetAlbIngressControllerPolicy, &opts, into)
}

func (a *albIngressControllerAPI) CreateAlbIngressControllerServiceAccount(opts api.CreateAlbIngressControllerServiceAccountOpts) (*api.ServiceAccount, error) {
	into := &api.ServiceAccount{}
	return into, a.client.DoPost(TargetAlbIngressControllerServiceAccount, &opts, into)
}

func (a *albIngressControllerAPI) CreateAlbIngressControllerHelmChart(opts api.CreateAlbIngressControllerHelmChartOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, a.client.DoPost(TargetAlbIngressControllerHelm, &opts, into)
}

// NewALBIngressControllerAPI returns an initialised REST API client
func NewALBIngressControllerAPI(client *HTTPClient) client.ALBIngressControllerAPI {
	return &albIngressControllerAPI{
		client: client,
	}
}
