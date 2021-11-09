package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeClusterCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var cluster api.ClusterCreateOpts

	err := json.NewDecoder(r.Body).Decode(&cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func decodeClusterDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var cluster api.ClusterDeleteOpts

	err := json.NewDecoder(r.Body).Decode(&cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func decodeGetClusterSecurityGroupIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.ClusterSecurityGroupIDGetOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
