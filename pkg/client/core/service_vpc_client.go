package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type vpcService struct {
	spinner spinner.Spinner
	api     client.VPCAPI
	store   client.VPCStore
	report  client.VPCReport
	state   client.VPCState
}

func (s *vpcService) GetVPC(ctx context.Context, id api.ID) (*api.Vpc, error) {
	return s.store.GetVpc(id)
}

func (s *vpcService) CreateVpc(_ context.Context, opts api.CreateVpcOpts) (*api.Vpc, error) {
	err := s.spinner.Start("vpc")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	vpc, err := s.api.CreateVpc(opts)
	if err != nil {
		return nil, err
	}

	r1, err := s.store.SaveVpc(vpc)
	if err != nil {
		return nil, err
	}

	r2, err := s.state.SaveVpc(vpc)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreateVPC(vpc, []*store.Report{r1, r2})
	if err != nil {
		return nil, err
	}

	return vpc, nil
}

func (s *vpcService) DeleteVpc(_ context.Context, opts api.DeleteVpcOpts) error {
	err := s.spinner.Start("vpc")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = s.api.DeleteVpc(opts)
	if err != nil {
		return err
	}

	_, err = s.store.DeleteVpc(opts.ID)
	if err != nil {
		return err
	}

	_, err = s.state.DeleteVpc(opts.ID)
	if err != nil {
		return err
	}

	return nil
}

// NewVPCService returns an initialised VPC service
func NewVPCService(spinner spinner.Spinner, api client.VPCAPI, store client.VPCStore, report client.VPCReport, state client.VPCState) client.VPCService {
	return &vpcService{
		spinner: spinner,
		api:     api,
		store:   store,
		report:  report,
		state:   state,
	}
}
