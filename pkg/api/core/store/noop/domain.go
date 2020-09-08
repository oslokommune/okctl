package noop

import "github.com/oslokommune/okctl/pkg/api"

type domainStore struct{}

func (s *domainStore) SaveHostedZone(_ *api.HostedZone) error {
	return nil
}

// NewDomainStore returns a no operation store
func NewDomainStore() api.DomainStore {
	return &domainStore{}
}
