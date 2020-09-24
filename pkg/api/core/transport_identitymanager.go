package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeCreateIdentityPool(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateIdentityPoolOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateIdentityPoolClient(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateIdentityPoolClientOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
