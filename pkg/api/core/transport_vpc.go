package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeVpcCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateVpcOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeVpcDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.DeleteVpcOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
