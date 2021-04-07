package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeCreateAWSLoadBalancerControllerPolicyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateAWSLoadBalancerControllerPolicyOpts

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

func decodeCreateAutoscalerPolicy(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateAutoscalerPolicy

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateBlockstoragePolicy(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateBlockstoragePolicy

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
