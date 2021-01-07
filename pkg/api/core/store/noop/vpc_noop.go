package noop

import "github.com/oslokommune/okctl/pkg/api"

type vpcStore struct{}

func (s *vpcStore) SaveVpc(*api.Vpc) error {
	return nil
}

func (s *vpcStore) DeleteVpc(api.ID) error {
	return nil
}

func (s *vpcStore) GetVpc(api.ID) (*api.Vpc, error) {
	return &api.Vpc{}, nil
}

// NewVpcStore returns a no operation store
func NewVpcStore() api.VpcStore {
	return &vpcStore{}
}
