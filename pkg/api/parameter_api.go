package api

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Parameter contains the state for a parameter
type Parameter struct {
	ID      ID
	Name    string
	Path    string
	Version int64
	Content string
}

// SecretParameter contains the state for a secret parameter
type SecretParameter struct {
	Parameter
}

// AnonymizeResponse removes sensitive data from the logs
func (p *SecretParameter) AnonymizeResponse(response interface{}) interface{} {
	r, _ := response.(*Parameter)
	rCopy := *r
	rCopy.Content = "XXXXXXXX"

	return &rCopy
}

// CreateSecretOpts contains the input required for creating a secret parameter
type CreateSecretOpts struct {
	ID     ID
	Name   string
	Secret string
}

// DeleteSecretOpts contains the input required for deleting a secret parameter
type DeleteSecretOpts struct {
	Name string
}

// Validate the inputs
func (o DeleteSecretOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Name, validation.Required),
	)
}

// AnonymizeRequest removes sensitive data from the logs
func (o CreateSecretOpts) AnonymizeRequest(request interface{}) interface{} {
	r, _ := request.(CreateSecretOpts)
	rCopy := r
	rCopy.Secret = "XXXXXXXXX"

	return rCopy
}

// Validate the inputs
func (o CreateSecretOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.ID, validation.Required),
		validation.Field(&o.Name, validation.Required),
		validation.Field(&o.Secret, validation.Required),
	)
}

// ParameterService defines the service layer operations
type ParameterService interface {
	CreateSecret(ctx context.Context, opts CreateSecretOpts) (*SecretParameter, error)
	DeleteSecret(ctx context.Context, opts DeleteSecretOpts) error
}

// ParameterCloudProvider defines the cloud layer operations
type ParameterCloudProvider interface {
	CreateSecret(opts CreateSecretOpts) (*SecretParameter, error)
	DeleteSecret(opts DeleteSecretOpts) error
}
