package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type parameterService struct {
	spinner spinner.Spinner
	api     client.ParameterAPI
	store   client.ParameterStore
	report  client.ParameterReport
}

func (s *parameterService) DeleteSecret(ctx context.Context, provider v1alpha1.CloudProvider, name string) error {
	err := s.spinner.Start("parameter")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	fmt.Println(name)

	_, err = provider.SSM().DeleteParameter(&ssm.DeleteParameterInput{
		Name: aws.String(name),
	})

	if err != nil {
		if strings.Contains(err.Error(), "ParameterNotFound:") {
			return nil
		}

		return err
	}

	return nil
}

func (s *parameterService) CreateSecret(_ context.Context, opts api.CreateSecretOpts) (*api.SecretParameter, error) {
	err := s.spinner.Start("parameter")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	secret, err := s.api.CreateSecret(opts)
	if err != nil {
		return nil, err
	}

	report, err := s.store.SaveSecret(secret)
	if err != nil {
		return nil, err
	}

	err = s.report.SaveSecret(secret, report)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// NewParameterService returns an initialised service
func NewParameterService(spinner spinner.Spinner, api client.ParameterAPI, store client.ParameterStore, report client.ParameterReport) client.ParameterService {
	return &parameterService{
		spinner: spinner,
		api:     api,
		store:   store,
		report:  report,
	}
}
