package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/oslokommune/okctl/pkg/truncate"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/mishudark/errors"

	"github.com/sanity-io/litter"
)

// HTTPClient stores state for invoking API operations
type HTTPClient struct {
	BaseURL  string
	Client   *http.Client
	Progress io.Writer
	Debug    bool
}

// New returns a client that wraps the common API operations
func New(debug bool, progress io.Writer, serverURL string) *HTTPClient {
	return &HTTPClient{
		Progress: progress,
		BaseURL:  serverURL,
		Client:   &http.Client{},
		Debug:    debug,
	}
}

// DoPost sends a POST request to the given endpoint
func (c *HTTPClient) DoPost(endpoint string, body interface{}, into interface{}) error {
	return c.Do(http.MethodPost, endpoint, body, into)
}

// DoDelete sends a DELETE request to the given endpoint
func (c *HTTPClient) DoDelete(endpoint string, body interface{}) error {
	return c.Do(http.MethodDelete, endpoint, body, nil)
}

// Do performs the request
// nolint: funlen
func (c *HTTPClient) Do(method, endpoint string, body interface{}, into interface{}) error {
	if c.Debug {
		_, err := fmt.Fprintf(c.Progress, "client (method: %s, endpoint: %s) starting request: %s", method, endpoint, litter.Sdump(body))
		if err != nil {
			return fmt.Errorf("failed to write debug output: %w", err)
		}
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("%s: %w", pretty("failed to marshal data for", method, endpoint), err)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.BaseURL, endpoint), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("%s: %w", pretty("failed to create request for", method, endpoint), err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", pretty("request failed for", method, endpoint), err)
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: %w", pretty("failed to read response for", method, endpoint), err)
	}

	if resp.StatusCode >= 400 { // nolint: gomnd
		return deserializeErrorPayload(out)
	}

	defer func() {
		err = resp.Body.Close()
	}()

	const logLineMaxlength = 5000

	if into != nil {
		if c.Debug {
			truncatedOut := truncate.Bytes(out, logLineMaxlength)

			_, err = fmt.Fprintf(c.Progress, "client (method: %s, endpoint: %s) received data: %s", method, endpoint, truncatedOut)
			if err != nil {
				return fmt.Errorf("failed to write debug output: %w", err)
			}
		}

		err = json.Unmarshal(out, into)
		if err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	if c.Debug {
		truncatedOut := truncate.Bytes(out, logLineMaxlength)

		_, err = io.Copy(c.Progress, strings.NewReader(string(truncatedOut)))
		if err != nil {
			return fmt.Errorf("%s: %w", pretty("failed to write progress for", method, endpoint), err)
		}
	}

	return nil
}

func pretty(msg, method, endpoint string) string {
	return fmt.Sprintf("%s: %s, %s", msg, method, endpoint)
}

type serializedError struct {
	Detail map[string]interface{} `json:"detail,omitempty"`
	Type   string                 `json:"type"`
	Error  string                 `json:"error"`
	Code   errors.Kind            `json:"code"`
}

func (s serializedError) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Error, validation.Required),
	)
}

// deserializeErrorPayload converts an JSON error payload from the backend to a structured error
func deserializeErrorPayload(jsonContent []byte) error {
	data := serializedError{}

	err := json.Unmarshal(jsonContent, &data)
	if err != nil {
		return errors.E(
			fmt.Errorf("unmarshalling error from server side: %w: %s", err, string(jsonContent)),
			errors.Internal)
	}

	if err = data.Validate(); err != nil {
		return errors.E(
			fmt.Errorf("validating deserialized error with content %s: %w", string(jsonContent), err),
			errors.Unmarshal,
		)
	}

	return errors.E(errors.New(data.Error), data.Code, data.Detail)
}
