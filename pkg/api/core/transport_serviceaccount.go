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
