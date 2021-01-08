package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetParameterSecret matches the REST API route
const TargetParameterSecret = "parameters/secret/"

type parameterAPI struct {
	client *HTTPClient
}

func (a *parameterAPI) DeleteSecret(opts api.DeleteSecretOpts) error {
	return a.client.DoDelete(TargetParameterSecret, &opts)
}

func (a *parameterAPI) CreateSecret(opts api.CreateSecretOpts) (*api.SecretParameter, error) {
	into := &api.SecretParameter{}
	return into, a.client.DoPost(TargetParameterSecret, &opts, into)
}

// NewParameterAPI returns an initialised REST API client
func NewParameterAPI(client *HTTPClient) client.ParameterAPI {
	return &parameterAPI{
		client: client,
	}
}
