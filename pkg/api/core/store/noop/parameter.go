package noop

import "github.com/oslokommune/okctl/pkg/api"

type parameterStore struct{}

func (s *parameterStore) SaveSecret(*api.SecretParameter) error {
	return nil
}

// NewParameterStore returns a no operation store
func NewParameterStore() api.ParameterStore {
	return &parameterStore{}
}
