package filesystem

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type domainStore struct {
	paths Paths
	fs    *afero.Afero
}

// Domain contains the outputs we will store
type Domain struct {
	Repository   string
	Environment  string
	FQDN         string
	Domain       string
	HostedZoneID string
	NameServers  []string
	StackName    string
}

func (d *domainStore) SaveDomain(domain *api.Domain) error {
	p := Domain{
		Repository:   domain.Repository,
		Environment:  domain.Environment,
		FQDN:         domain.FQDN,
		Domain:       domain.Domain,
		HostedZoneID: domain.HostedZoneID,
		NameServers:  domain.NameServers,
		StackName:    domain.StackName,
	}

	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	err = d.fs.MkdirAll(path.Join(d.paths.BaseDir, domain.Domain), 0o744)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err = d.fs.WriteFile(path.Join(d.paths.BaseDir, domain.Domain, d.paths.OutputFile), data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write outputs: %w", err)
	}

	err = d.fs.WriteFile(path.Join(d.paths.BaseDir, domain.Domain, d.paths.CloudFormationFile), domain.CloudFormationTemplate, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write cloud formation template: %w", err)
	}

	return nil
}

// NewDomainStore returns an initialised domain store
func NewDomainStore(paths Paths, fs *afero.Afero) api.DomainStore {
	return &domainStore{
		paths: paths,
		fs:    fs,
	}
}
