package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeCreateExternalSecretsPolicyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateExternalSecretsPolicyOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateAlbIngressControllerPolicyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateAlbIngressControllerPolicyOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateExternalDNSPolicyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateExternalDNSPolicyOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
