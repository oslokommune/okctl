package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeClusterCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var cluster api.ClusterCreateOpts

	return cluster, json.NewDecoder(r.Body).Decode(&cluster)
}

func decodeClusterDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var cluster api.ClusterDeleteOpts

	return cluster, json.NewDecoder(r.Body).Decode(&cluster)
}
