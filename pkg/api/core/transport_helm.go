package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeCreateHelmRelease(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateHelmReleaseOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeDeleteHelmRelease(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.DeleteHelmReleaseOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeGetHelmRelease(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.GetHelmReleaseOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
