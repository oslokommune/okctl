//Package transport implements some basic functionality for transport encoding
package transport

import (
	"bytes"
	"context"
	"io"
	"net/http"

	kit "github.com/go-kit/kit/transport/http"
	"github.com/sanity-io/litter"
	"sigs.k8s.io/yaml"
)

// Texter defines the interface types must implement to
// control the text output of their response
type Texter interface {
	Text() []byte
}

// EncodeYAMLResponse writes a YAML serialised response to the receiver
func EncodeYAMLResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if headerer, ok := response.(kit.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}

	code := http.StatusOK

	if sc, ok := response.(kit.StatusCoder); ok {
		code = sc.StatusCode()
	}

	w.WriteHeader(code)

	if code == http.StatusNoContent {
		return nil
	}

	data, err := yaml.Marshal(response)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, bytes.NewReader(data))
	if err != nil {
		return err
	}

	return nil
}

// EncodeTextResponse writes a plaintext response to the receiver
func EncodeTextResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if headerer, ok := response.(kit.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}

	code := http.StatusOK

	if sc, ok := response.(kit.StatusCoder); ok {
		code = sc.StatusCode()
	}

	w.WriteHeader(code)

	if code == http.StatusNoContent {
		return nil
	}

	data := []byte(litter.Sdump(response))
	if t, ok := response.(Texter); ok {
		data = t.Text()
	}

	_, err := io.Copy(w, bytes.NewReader(data))
	if err != nil {
		return err
	}

	return nil
}
