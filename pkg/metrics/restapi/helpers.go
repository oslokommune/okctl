package restapi

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

func postEvent(apiURL url.URL, payload []byte) error {
	request, err := http.NewRequest(http.MethodPost, apiURL.String(), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("preparing request: %w", err)
	}

	request.Header.Set("User-Agent", requiredUserAgent)
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
