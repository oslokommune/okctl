package noop

import "github.com/oslokommune/okctl/pkg/api"

type domainStore struct{}

func (s *domainStore) SaveDomain(_ *api.Domain) error {
	return nil
}

// NewDomainStore returns a no operation store
func NewDomainStore() api.DomainStore {
	return &domainStore{}
}
