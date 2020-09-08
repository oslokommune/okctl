package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

const (
	// TargetExternalSecretsPolicy is the API route for the policy
	TargetExternalSecretsPolicy = "managedpolicies/externalsecrets/"
	// TargetExternalSecretsServiceAccount is the API route for the service account
	TargetExternalSecretsServiceAccount = "serviceaccounts/externalsecrets/"
	// TargetExternalSecretsHelm is the API route for helm
	TargetExternalSecretsHelm = "helm/externalsecrets/"
)

type externalSecretsAPI struct {
	client *HTTPClient
}

func (a *externalSecretsAPI) CreateExternalSecretsPolicy(opts api.CreateExternalSecretsPolicyOpts) (*api.ManagedPolicy, error) {
	into := &api.ManagedPolicy{}
	return into, a.client.DoPost(TargetExternalSecretsPolicy, &opts, into)
}

func (a *externalSecretsAPI) CreateExternalSecretsServiceAccount(opts api.CreateExternalSecretsServiceAccountOpts) (*api.ServiceAccount, error) {
	into := &api.ServiceAccount{}
	return into, a.client.DoPost(TargetExternalSecretsServiceAccount, &opts, into)
}

func (a *externalSecretsAPI) CreateExternalSecretsHelmChart(opts api.CreateExternalSecretsHelmChartOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, a.client.DoPost(TargetExternalSecretsHelm, &opts, into)
}

// NewExternalSecretsAPI returns an initialised API client
func NewExternalSecretsAPI(client *HTTPClient) client.ExternalSecretsAPI {
	return &externalSecretsAPI{
		client: client,
	}
}
