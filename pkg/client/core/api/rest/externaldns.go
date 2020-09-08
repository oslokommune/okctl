package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

const (
	// TargetExternalDNSPolicy matches the REST API route
	TargetExternalDNSPolicy = "managedpolicies/externaldns/"
	// TargetExternalDNSServiceAccount matches the REST API route
	TargetExternalDNSServiceAccount = "serviceaccounts/externaldns/"
	// TargetKubeExternalDNS matches the REST API route
	TargetKubeExternalDNS = "kube/externaldns/"
)

type externalDNSAPI struct {
	client *HTTPClient
}

func (a *externalDNSAPI) CreateExternalDNSPolicy(opts api.CreateExternalDNSPolicyOpts) (*api.ManagedPolicy, error) {
	into := &api.ManagedPolicy{}
	return into, a.client.DoPost(TargetExternalDNSPolicy, &opts, into)
}

func (a *externalDNSAPI) CreateExternalDNSServiceAccount(opts api.CreateExternalDNSServiceAccountOpts) (*api.ServiceAccount, error) {
	into := &api.ServiceAccount{}
	return into, a.client.DoPost(TargetExternalDNSServiceAccount, &opts, into)
}

func (a *externalDNSAPI) CreateExternalDNSKubeDeployment(opts api.CreateExternalDNSKubeDeploymentOpts) (*api.Kube, error) {
	into := &api.Kube{}
	return into, a.client.DoPost(TargetKubeExternalDNS, &opts, into)
}

// NewExternalDNSAPI returns an initialised REST API client
func NewExternalDNSAPI(client *HTTPClient) client.ExternalDNSAPI {
	return &externalDNSAPI{
		client: client,
	}
}
