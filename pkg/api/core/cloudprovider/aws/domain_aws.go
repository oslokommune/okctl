package aws

import (
	"fmt"

	"github.com/gosimple/slug"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
	"github.com/oslokommune/okctl/pkg/cfn/components/hostedzone"
	"github.com/oslokommune/okctl/pkg/route53"
)

type domain struct {
	provider v1alpha1.CloudProvider
}

func (d *domain) DeleteHostedZone(opts api.DeleteHostedZoneOpts) error {
	err := cfn.NewRunner(d.provider).Delete(cfn.NewStackNamer().AliasRecordSet(opts.ID.ClusterName+"-auth", slug.Make(opts.Domain)))
	if err != nil {
		return err
	}

	_, err = route53.New(d.provider).DeleteHostedZoneRecordSets(opts.HostedZoneID)
	if err != nil {
		return err
	}

	err = cfn.NewRunner(d.provider).Delete(cfn.NewStackNamer().Domain(opts.ID.ClusterName, slug.Make(opts.Domain)))
	if err != nil {
		return err
	}

	return nil
}

func (d *domain) CreateHostedZone(opts api.CreateHostedZoneOpts) (*api.HostedZone, error) {
	// Create hosted zone
	b := cfn.New(components.NewHostedZoneComposer(opts.FQDN, "A public hosted zone for creating ingresses with"))

	stackName := cfn.NewStackNamer().Domain(opts.ID.ClusterName, slug.Make(opts.Domain))

	template, err := b.Build()
	if err != nil {
		return nil, errors.E(err, "failed to build cloud formation template")
	}

	template, err = hostedzone.PatchYAML(template)
	if err != nil {
		return nil, errors.E(err, "failed to patch template body")
	}

	r := cfn.NewRunner(d.provider)

	err = r.CreateIfNotExists(opts.ID.ClusterName, stackName, template, nil, defaultTimeOut)
	if err != nil {
		return nil, errors.E(err, "failed to create cloud formation template")
	}

	p := &api.HostedZone{
		ID:                     opts.ID,
		Managed:                true,
		FQDN:                   opts.FQDN,
		Domain:                 opts.Domain,
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

	// Adjust TTL of NS record, as it can't be set when created
	NSTTL := opts.NSTTL
	if NSTTL == 0 {
		NSTTL = 900
	}

	err = route53.New(d.provider).SetNSRecordTTL(p.HostedZoneID, NSTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to set NS record TTL: %w", err)
	}

	return p, nil
}

// NewDomainCloudProvider returns an initialised cloud provider for domains
func NewDomainCloudProvider(provider v1alpha1.CloudProvider) api.DomainCloudProvider {
	return &domain{
		provider: provider,
	}
}
