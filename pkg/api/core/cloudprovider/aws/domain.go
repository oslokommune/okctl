package aws

import (
	"github.com/gosimple/slug"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
	"github.com/oslokommune/okctl/pkg/cfn/components/hostedzone"
)

type domain struct {
	provider v1alpha1.CloudProvider
}

func (d *domain) CreateHostedZone(opts api.CreateHostedZoneOpts) (*api.HostedZone, error) {
	b := cfn.New(components.NewHostedZoneComposer(opts.FQDN, "A public hosted zone for creating ingresses with"))

	stackName := cfn.NewStackNamer().Domain(opts.ID.Repository, opts.ID.Environment, slug.Make(opts.Domain))

	template, err := b.Build()
	if err != nil {
		return nil, errors.E(err, "failed to build cloud formation template")
	}

	template, err = hostedzone.PatchYAML(template)
	if err != nil {
		return nil, errors.E(err, "failed to patch template body")
	}

	r := cfn.NewRunner(d.provider)

	err = r.CreateIfNotExists(stackName, template, nil, defaultTimeOut)
	if err != nil {
		return nil, errors.E(err, "failed to create cloud formation template")
	}

	p := &api.HostedZone{
		ID:                     opts.ID,
		Managed:                true,
		FQDN:                   opts.FQDN,
		Domain:                 opts.Domain,
		HostedZoneID:           "",
		NameServers:            nil,
		StackName:              stackName,
		CloudFormationTemplate: template,
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		"PublicHostedZone": cfn.String(&p.HostedZoneID),
		"NameServers":      cfn.StringSlice(&p.NameServers),
	})
	if err != nil {
		return nil, errors.E(err, "failed to process outputs")
	}

	return p, nil
}

// NewDomainCloudProvider returns an initialised cloud provider for domains
func NewDomainCloudProvider(provider v1alpha1.CloudProvider) api.DomainCloudProvider {
	return &domain{
		provider: provider,
	}
}
