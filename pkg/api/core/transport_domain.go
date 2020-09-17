package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeCreateHostedZone(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateHostedZoneOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeDeleteHostedZone(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.DeleteHostedZoneOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
