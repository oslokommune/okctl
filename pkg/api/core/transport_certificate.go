package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeCreateCertificate(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateCertificateOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeDeleteCertificate(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.DeleteCertificateOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeDeleteCognitoCertificate(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.DeleteCognitoCertificateOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
