package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

type parameter struct {
	provider v1alpha1.CloudProvider
}

func (p *parameter) DeleteSecret(opts api.DeleteSecretOpts) error {
	path := p.getSecretPath(opts.ID.ClusterName, opts.Name)

	_, err := p.provider.SSM().DeleteParameter(&ssm.DeleteParameterInput{
		Name: aws.String(path),
	})
	if err != nil {
		if strings.Contains(err.Error(), "ParameterNotFound:") {
			return nil
		}

		return err
	}

	return nil
}

func (p *parameter) CreateSecret(opts api.CreateSecretOpts) (*api.SecretParameter, error) {
	path := p.getSecretPath(opts.ID.ClusterName, opts.Name)

	got, err := p.provider.SSM().PutParameter(&ssm.PutParameterInput{
		DataType:    aws.String("text"),
		Description: aws.String(fmt.Sprintf("Secret: %s, created by okctl", opts.Name)),
		Name:        aws.String(path),
		Overwrite:   aws.Bool(true),
		Type:        aws.String(ssm.ParameterTypeSecureString),
		Value:       aws.String(opts.Secret),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create parameter: %w", err)
	}

	return &api.SecretParameter{
		Parameter: api.Parameter{
			ID:      opts.ID,
			Name:    opts.Name,
			Version: *got.Version,
			Path:    path,
		},
	}, nil
}

func (p *parameter) getSecretPath(clusterName, secretName string) string {
	return fmt.Sprintf("/okctl/%s/%s", clusterName, secretName)
}

// NewParameterCloudProvider returns an initialised cloud provider
func NewParameterCloudProvider(provider v1alpha1.CloudProvider) api.ParameterCloudProvider {
	return &parameter{
		provider: provider,
	}
}
