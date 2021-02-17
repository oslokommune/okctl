package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/kube/manifests/storageclass"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type blockstorageService struct {
	spinner spinner.Spinner
	api     client.BlockstorageAPI
	store   client.BlockstorageStore
	report  client.BlockstorageReport
	kube    client.ManifestService
}

func (s *blockstorageService) DeleteBlockstorage(_ context.Context, id api.ID) error {
	err := s.spinner.Start("blockstorage")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	err = s.api.DeleteBlockstorageServiceAccount(id)
	if err != nil {
		return err
	}

	err = s.api.DeleteBlockstoragePolicy(id)
	if err != nil {
		return err
	}

	report, err := s.store.RemoveBlockstorage(id)
	if err != nil {
		return err
	}

	err = s.report.ReportDeleteBlockstorage(report)
	if err != nil {
		return err
	}

	return nil
}

func (s *blockstorageService) CreateBlockstorage(ctx context.Context, opts client.CreateBlockstorageOpts) (*client.Blockstorage, error) {
	err := s.spinner.Start("blockstorage")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	policy, err := s.api.CreateBlockstoragePolicy(api.CreateBlockstoragePolicy{
		ID: opts.ID,
	})
	if err != nil {
		return nil, err
	}

	sa, err := s.api.CreateBlockstorageServiceAccount(api.CreateBlockstorageServiceAccountOpts{
		CreateServiceAccountOpts: api.CreateServiceAccountOpts{
			ID:        opts.ID,
			PolicyArn: policy.PolicyARN,
		},
	})
	if err != nil {
		return nil, err
	}

	chart, err := s.api.CreateBlockstorageHelmChart(api.CreateBlockstorageHelmChartOpts{
		ID: opts.ID,
	})
	if err != nil {
		return nil, err
	}

	a := &client.Blockstorage{
		Policy:         policy,
		ServiceAccount: sa,
		Chart:          chart,
	}

	_, err = s.kube.CreateStorageClass(ctx, api.CreateStorageClassOpts{
		ID:          opts.ID,
		Name:        "ebs-sc",
		Parameters:  storageclass.NewEBSParameters(),
		Annotations: storageclass.DefaultStorageClassAnnotation(),
	})
	if err != nil {
		return nil, fmt.Errorf("creating default storage class: %w", err)
	}

	report, err := s.store.SaveBlockstorage(a)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreateBlockstorage(a, report)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// NewBlockstorageService returns an initialised service
func NewBlockstorageService(
	spinner spinner.Spinner,
	api client.BlockstorageAPI,
	store client.BlockstorageStore,
	report client.BlockstorageReport,
	kube client.ManifestService,
) client.BlockstorageService {
	return &blockstorageService{
		spinner: spinner,
		api:     api,
		store:   store,
		report:  report,
		kube:    kube,
	}
}
