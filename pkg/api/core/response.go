package core

import (
	"encoding/json"
	"net/http"
)

// Created wraps a resource that has just been created
type Created struct {
	Data interface{}
}

// StatusCode returns the http status code
func (c Created) StatusCode() int {
	return http.StatusCreated
}

// MarshalJSON encodes the inner data to json
func (c Created) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Data)
}

// Success indicates that all went well
type Success struct {
	Success bool `json:"success"`
}

// Ok returns a created ok
func Ok() Created {
	return Created{
		Data: Success{
			Success: true,
		},
	}
}
