package metrics

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

const metricsPath = "/v1/metrics/events"

func postEvent(apiURL url.URL, userAgent string, payload []byte) error {
	apiURL.Path = metricsPath

	request, err := http.NewRequest(http.MethodPost, apiURL.String(), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("preparing request: %w", err)
	}

	request.Header.Set("User-Agent", userAgent)
	request.Header.Set("Content-Type", "application/json")

	httpClient := http.Client{}

	response, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("posting metrics event: %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code %d: %s", response.StatusCode, response.Status)
	}

	return nil
}
