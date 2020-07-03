// Package request provides some http client request helpers
package request

import (
	"bytes"
	"fmt"
	"net/http"
)

// Methods defines the available http methods
type Methods interface {
	Post(endpoint string, body []byte) (string, error)
	Delete(endpoint string, body []byte) (string, error)
}

type request struct {
	baseURL string
	client  *http.Client
}

// New returns a new request helper
func New(baseURL string) Methods {
	return &request{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Do performs the request
func (r *request) Do(method, endpoint string, body []byte) (string, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", r.baseURL, endpoint), bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}

	out := make([]byte, resp.ContentLength)

	_, err = resp.Body.Read(out)
	if err != nil {
		return "", err
	}

	defer func() {
		err = resp.Body.Close()
	}()

	return string(out), nil
}

// Post sends a POST request to the given endpoint
func (r *request) Post(endpoint string, body []byte) (string, error) {
	return r.Do(http.MethodPost, endpoint, body)
}

// Delete sends a DELETE request to the given endpoint
func (r *request) Delete(endpoint string, body []byte) (string, error) {
	return r.Do(http.MethodDelete, endpoint, body)
}
