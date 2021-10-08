package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/oslokommune/okctl-metrics-service/pkg/endpoints/metrics"
	okmetrics "github.com/oslokommune/okctl/pkg/metrics"
)

const requiredUserAgent = "okctl"

// Publish posts an event to an okctl metrics backend
func (c client) Publish(_ context.Context, event metrics.Event) error {
	err := event.Validate()
	if err != nil {
		return fmt.Errorf("validating metrics event: %w", err)
	}

	rawEvent, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshalling event as JSON: %w", err)
	}

	err = postEvent(c.apiURL, rawEvent)
	if err != nil {
		return fmt.Errorf("POSTing event: %w", err)
	}

	return nil
}

// NewMetricsClient initializes a metrics REST API client with a specific API URL
func NewMetricsClient(apiURL url.URL) okmetrics.Client {
	return &client{apiURL: apiURL}
}

// NewDefaultMetricsClient initializes a metrics REST API client with the default API URL
func NewDefaultMetricsClient() okmetrics.Client {
	rawURL := "https://metrics.kjoremiljo.oslo.systems"

	apiURL, err := url.Parse(rawURL)
	if err != nil {
		panic(fmt.Sprintf("apiURL '%s' did not compile", rawURL))
	}

	return NewMetricsClient(*apiURL)
}
