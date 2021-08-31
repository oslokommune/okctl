package aws

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/cleaner"

	"github.com/gosimple/slug"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
	"github.com/oslokommune/okctl/pkg/cognito"
)

const (
	certificateTimeout = 15
)

type certificate struct {
	provider v1alpha1.CloudProvider
}

func (c *certificate) DeleteCognitoCertificate(opts api.DeleteCognitoCertificateOpts) error {
	d := cognito.NewCertDeleter(c.provider)

	err := d.DeleteAuthCert(cognito.DeleteAuthCertOpts{
		Domain: opts.Domain,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *certificate) DeleteCertificate(opts api.DeleteCertificateOpts) error {
	err := cleaner.New(c.provider).RemoveThingsUsingCertForDomain(opts.Domain)
	if err != nil {
		return fmt.Errorf(constant.RemoveCertificateUsagesError, err)
	}

	return cfn.NewRunner(c.provider).Delete(cfn.NewStackNamer().Certificate(opts.ID.ClusterName, slug.Make(opts.Domain)))
}

func (c *certificate) CreateCertificate(opts api.CreateCertificateOpts) (*api.Certificate, error) {
	b := cfn.New(components.NewPublicCertificateComposer(opts.Domain, opts.HostedZoneID))

	stackName := cfn.NewStackNamer().Certificate(opts.ID.ClusterName, slug.Make(opts.Domain))

	template, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf(constant.BuildCloudFormationTemplateError, err)
	}

	r := cfn.NewRunner(c.provider)

	err = r.CreateIfNotExists(opts.ID.ClusterName, stackName, template, nil, certificateTimeout)
	if err != nil {
		return nil, fmt.Errorf(constant.ApplyCloudformationTemplateError, err)
	}

	p := &api.Certificate{
		ID:                     opts.ID,
		FQDN:                   opts.FQDN,
		Domain:                 opts.Domain,
		HostedZoneID:           opts.HostedZoneID,
		StackName:              stackName,
		CloudFormationTemplate: template,
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		"PublicCertificate": cfn.String(&p.CertificateARN),
	})
	if err != nil {
		return nil, fmt.Errorf(constant.ProcessOutputsError, err)
	}

	return p, nil
}

// NewCertificateCloudProvider returns an initialised cloud provider
func NewCertificateCloudProvider(provider v1alpha1.CloudProvider) api.CertificateCloudProvider {
	return &certificate{
		provider: provider,
	}
}
