package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type vpcService struct {
	api    client.VPCAPI
	store  client.VPCStore
	report client.VPCReport
}

func (s *vpcService) CreateVpc(_ context.Context, opts api.CreateVpcOpts) (*api.Vpc, error) {
	vpc, err := s.api.CreateVpc(opts)
	if err != nil {
		return nil, err
	}

	report, err := s.store.SaveVpc(vpc)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreateVPC(vpc, report)
	if err != nil {
		return nil, err
	}

	return vpc, nil
}

func (s *vpcService) DeleteVpc(_ context.Context, opts api.DeleteVpcOpts) error {
	err := s.api.DeleteVpc(opts)
	if err != nil {
		return err
	}

	_, err = s.store.DeleteVpc(opts.ID)
	if err != nil {
		return err
	}

	return nil
}

// NewVPCService returns an initialised VPC service
func NewVPCService(api client.VPCAPI, store client.VPCStore, report client.VPCReport) client.VPCService {
	return &vpcService{
		api:    api,
		store:  store,
		report: report,
	}
}
