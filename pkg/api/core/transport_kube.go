package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeCreateExternalDNSKubeDeployment(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateExternalDNSKubeDeploymentOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateExternalSecrets(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateExternalSecretsOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeDeleteNamespace(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.DeleteNamespaceOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateStorageClass(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateStorageClassOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeDeleteExternalSecrets(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.DeleteExternalSecretsOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
