package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeCreateExternalSecretsHelmChart(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateExternalSecretsHelmChartOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateAlbIngressControllerHelmChart(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateAlbIngressControllerHelmChartOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateAWSLoadBalancerControllerHelmChart(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateAWSLoadBalancerControllerHelmChartOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateArgoCD(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateArgoCDOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
