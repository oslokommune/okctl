package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeCreateAutoscalerServiceAccount(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateAutoscalerServiceAccountOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateBlockstorageServiceAccount(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateBlockstorageServiceAccountOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
