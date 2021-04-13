package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	kit "github.com/go-kit/kit/transport/http"
)

// decodeStructRequest can be expanded in the future to act on a given type
// of `Accept` header, e.g., for marshalling from yaml or other formats.
func decodeStructRequest(v interface{}) kit.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			return nil, fmt.Errorf("decoding request as json: %w", err)
		}

		return v, nil
	}
}
