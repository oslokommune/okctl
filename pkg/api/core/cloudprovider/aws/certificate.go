package aws

import (
	"fmt"
	"strings"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
)

type certificate struct {
	provider v1alpha1.CloudProvider
}

func (c *certificate) CreateCertificate(opts api.CreateCertificateOpts) (*api.Certificate, error) {
	b := cfn.New(components.NewPublicCertificateComposer(opts.FQDN, opts.HostedZoneID))

	parts := strings.SplitN(opts.Domain, ".", 2)
	if len(parts) != 2 { // nolint: gomnd
		return nil, fmt.Errorf("failed to extract sub domain")
	}

	stackName := cfn.NewStackNamer().Certificate(opts.Repository, opts.Environment, parts[0])

	template, err := b.Build()
	if err != nil {
		return nil, errors.E(err, "failed to build cloud formation template")
	}

	r := cfn.NewRunner(c.provider)

	err = r.CreateIfNotExists(stackName, template, nil, defaultTimeOut)
	if err != nil {
		return nil, errors.E(err, "failed to create cloud formation template")
	}

	p := &api.Certificate{
		Repository:             opts.Repository,
		Environment:            opts.Environment,
		FQDN:                   opts.FQDN,
		Domain:                 opts.Domain,
		HostedZoneID:           opts.HostedZoneID,
		StackName:              stackName,
		CloudFormationTemplate: template,
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		"Certificate": cfn.String(&p.CertificateARN),
	})
	if err != nil {
		return nil, errors.E(err, "failed to process outputs")
	}

	return p, nil
}

// NewCertificateCloudProvider returns an initialised cloud provider
func NewCertificateCloudProvider(provider v1alpha1.CloudProvider) api.CertificateCloudProvider {
	return &certificate{
		provider: provider,
	}
}
