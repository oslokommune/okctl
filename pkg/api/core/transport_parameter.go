package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeCreateSecret(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateSecretOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeDeleteSecret(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.DeleteSecretOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
