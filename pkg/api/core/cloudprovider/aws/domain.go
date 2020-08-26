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

type domain struct {
	provider v1alpha1.CloudProvider
}

func (d *domain) CreateDomain(opts api.CreateDomainOpts) (*api.Domain, error) {
	b := cfn.New(components.NewHostedZoneComposer(opts.FQDN, "A public hosted zone for creating ingresses with"))

	parts := strings.SplitN(opts.Domain, ".", 2)
	if len(parts) != 2 { // nolint: gomnd
		return nil, fmt.Errorf("failed to extract sub domain")
	}

	stackName := cfn.NewStackNamer().Domain(opts.Repository, opts.Environment, parts[0])

	template, err := b.Build()
	if err != nil {
		return nil, errors.E(err, "failed to build cloud formation template")
	}

	r := cfn.NewRunner(d.provider)

	err = r.CreateIfNotExists(stackName, template, nil, defaultTimeOut)
	if err != nil {
		return nil, errors.E(err, "failed to create cloud formation template")
	}

	p := &api.Domain{
		Repository:  opts.Repository,
		Environment: opts.Environment,
		FQDN:        opts.FQDN,
		Domain:      opts.Domain,
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
