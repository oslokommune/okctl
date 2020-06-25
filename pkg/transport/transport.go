package transport

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	kit "github.com/go-kit/kit/transport/http"
	"sigs.k8s.io/yaml"
)

type Texter interface {
	Text() []byte
}

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

	data := []byte(spew.Sdump(response))
	if t, ok := response.(Texter); ok {
		data = t.Text()
	}

	_, err := io.Copy(w, bytes.NewReader(data))
	if err != nil {
		return err
	}

	return nil
}
