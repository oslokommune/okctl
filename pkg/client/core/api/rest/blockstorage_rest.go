package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

const (
	// TargetBlockstoragePolicy is the API route for the policy
	TargetBlockstoragePolicy = "managedpolicies/blockstorage/"
	// TargetBlockstorageServiceAccount is the API route for the service account
	TargetBlockstorageServiceAccount = "serviceaccounts/blockstorage/"
	// TargetBlockstorageHelm is the API route for helm
	TargetBlockstorageHelm = "helm/blockstorage/"
)

type blockstorageAPI struct {
	client *HTTPClient
}

func (a *blockstorageAPI) DeleteBlockstoragePolicy(id api.ID) error {
	return a.client.DoDelete(TargetBlockstoragePolicy, &id)
}

func (a *blockstorageAPI) DeleteBlockstorageServiceAccount(id api.ID) error {
	return a.client.DoDelete(TargetBlockstorageServiceAccount, &id)
}

func (a *blockstorageAPI) CreateBlockstoragePolicy(opts api.CreateBlockstoragePolicy) (*api.ManagedPolicy, error) {
	into := &api.ManagedPolicy{}
	return into, a.client.DoPost(TargetBlockstoragePolicy, &opts, into)
}

func (a *blockstorageAPI) CreateBlockstorageServiceAccount(opts api.CreateBlockstorageServiceAccountOpts) (*api.ServiceAccount, error) {
	into := &api.ServiceAccount{}
	return into, a.client.DoPost(TargetBlockstorageServiceAccount, &opts, into)
}

func (a *blockstorageAPI) CreateBlockstorageHelmChart(opts api.CreateBlockstorageHelmChartOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, a.client.DoPost(TargetBlockstorageHelm, &opts, into)
}

// NewBlockstorageAPI returns an initialised API client
func NewBlockstorageAPI(client *HTTPClient) client.BlockstorageAPI {
	return &blockstorageAPI{
		client: client,
	}
}
