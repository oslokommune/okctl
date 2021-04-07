package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

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

func decodeCreateAutoscalerHelmChart(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateAutoscalerHelmChartOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateBlockstorageHelmChart(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateBlockstorageHelmChartOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateKubePrometheusStackHelmChart(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateKubePrometheusStackOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreateLokiHelmChart(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreateLokiHelmChartOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

func decodeCreatePromtailHelmChart(_ context.Context, r *http.Request) (interface{}, error) {
	var opts api.CreatePromtailHelmChartOpts

	err := json.NewDecoder(r.Body).Decode(&opts)
	if err != nil {
		return nil, err
	}

	return opts, nil
}

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
