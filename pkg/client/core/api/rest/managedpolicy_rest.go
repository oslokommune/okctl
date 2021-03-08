package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// ManagedPoliciesTarget defines the REST API Route
const ManagedPoliciesTarget = "managedpolicies/"

type managedPolicyAPI struct {
	client *HTTPClient
}

func (m *managedPolicyAPI) CreatePolicy(opts api.CreatePolicyOpts) (*api.ManagedPolicy, error) {
	into := &api.ManagedPolicy{}
	return into, m.client.DoPost(ManagedPoliciesTarget, &opts, into)
}

func (m *managedPolicyAPI) DeletePolicy(opts api.DeletePolicyOpts) error {
	return m.client.DoDelete(ManagedPoliciesTarget, &opts)
}

// NewManagedPolicyAPI returns an initialised API client
func NewManagedPolicyAPI(client *HTTPClient) client.ManagedPolicyAPI {
	return &managedPolicyAPI{
		client: client,
	}
}
