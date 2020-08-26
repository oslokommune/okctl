package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeCreateExternalSecretsServiceAccount(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateExternalSecretsServiceAccountOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateAlbIngressControllerServiceAccount(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateAlbIngressControllerServiceAccountOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateExternalDnsServiceAccount(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateExternalDnsServiceAccountOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
