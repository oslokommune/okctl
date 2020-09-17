package core

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/oslokommune/okctl/pkg/api"
)

func decodeIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var id api.ID

	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		return nil, err
	}

	return id, nil
}
