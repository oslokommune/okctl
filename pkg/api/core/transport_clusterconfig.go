package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeClusterConfigCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var cluster api.CreateClusterConfigOpts

	err := json.NewDecoder(r.Body).Decode(&cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}
