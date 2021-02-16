package controller

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/github"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
	"github.com/spf13/afero"
)

func getVpcState(fs *afero.Afero, outputDir string) api.Vpc {
	vpc := api.Vpc{}

	baseDir := path.Join(outputDir, "vpc")

	_, err := store.NewFileSystem(baseDir, fs).
		GetStruct(config.DefaultVpcOutputs, &vpc, store.FromJSON()).
		Do()
	if err != nil {
		panic(fmt.Errorf("error reading from vpc state file: %w", err))
	}

	return vpc
}

// StringFetcher defines a function which can be used to delay fetching of strings
type StringFetcher func() string

// HostedZoneFetcher defines a function which can be used to delay fetching of a hosted zone
type HostedZoneFetcher func() *state.HostedZone

// IdentityPoolFetcher defines a function which can be used to delay fetching of IdentityPool data
type IdentityPoolFetcher func() state.IdentityPool

// CreateClusterStateRefresher creates a function that gathers required runtime data for a cluster resource
func CreateClusterStateRefresher(fs *afero.Afero, outputDir string, cidrFn StringFetcher) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		vpc := getVpcState(fs, outputDir)

		vpc.Cidr = cidrFn()

		node.ResourceState = reconciler.ClusterResourceState{VPC: vpc}
	}
}

// CreateALBIngressControllerRefresher creates a function that gathers required runtime data for a ALB Ingress
// Controller resource
func CreateALBIngressControllerRefresher(fs *afero.Afero, outputDir string) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		vpc := getVpcState(fs, outputDir)

		node.ResourceState = reconciler.AlbIngressControllerResourceState{VpcID: vpc.VpcID}
	}
}

// CreateAWSLoadBalancerControllerRefresher creates a function that gathers required runtime data for AWS
// load balancer controller
func CreateAWSLoadBalancerControllerRefresher(fs *afero.Afero, outputDir string) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		vpc := getVpcState(fs, outputDir)

		node.ResourceState = reconciler.AWSLoadBalancerControllerResourceState{
			VpcID: vpc.VpcID,
		}
	}
}

// CreateExternalDNSStateRefresher creates a function that gathers required runtime data for a External DNS resource
func CreateExternalDNSStateRefresher(primaryHostedZoneFetcher HostedZoneFetcher) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		hostedZone := primaryHostedZoneFetcher()

		node.ResourceState = reconciler.ExternalDNSResourceState{
			HostedZoneID: hostedZone.ID,
			Domain:       hostedZone.Domain,
		}
	}
}

// CreateIdentityManagerRefresher creates a function that gathers required runtime data for a Identity Manager resource
func CreateIdentityManagerRefresher(primaryHostedZoneFetcher HostedZoneFetcher) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		hostedZone := primaryHostedZoneFetcher()

		node.ResourceState = reconciler.IdentityManagerResourceState{HostedZoneID: hostedZone.ID}
	}
}

// CreateGithubStateRefresher creates a function that gathers required runtime data for a Github resource
func CreateGithubStateRefresher(ghGetter reconciler.GithubGetter, ghSetter reconciler.GithubSetter) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		node.ResourceState = reconciler.GithubResourceState{
			Getter: ghGetter,
			Saver:  ghSetter,
		}
	}
}

type argoCDRefresherOptions struct {
	IACRepositoryName         string
	IACRepositoryOrganization string
	ClusterID                 api.ID
	hostedZoneFetcher         HostedZoneFetcher
	idpFetcher                IdentityPoolFetcher
	ghGetter                  reconciler.GithubGetter
}

// CreateArgocdStateRefresher creates a function that gathers required runtime data for a ArgoCD resource
func CreateArgocdStateRefresher(options argoCDRefresherOptions) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		idp := options.idpFetcher()

		originalRepo := options.ghGetter().Repositories[fmt.Sprintf("%s/%s", options.IACRepositoryOrganization, options.IACRepositoryName)]

		node.ResourceState = reconciler.ArgocdResourceState{
			HostedZone: options.hostedZoneFetcher(),
			Repository: stateToClientRepo(originalRepo, options.ClusterID),
			UserPoolID: idp.UserPoolID,
			AuthDomain: idp.AuthDomain,
		}
	}
}

// CreateNameserverDelegationStateRefresher creates a function that gathers required runtime data for a nameserver delegation
// request
func CreateNameserverDelegationStateRefresher(fetcher HostedZoneFetcher) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		zone := fetcher()

		node.ResourceState = reconciler.NameserverHandlerReconcilerResourceState{
			PrimaryHostedZoneFQDN: zone.FQDN,
			Nameservers:           zone.NameServers,
		}
	}
}

func stateToClientRepo(original state.GithubRepository, id api.ID) *client.GithubRepository {
	deployKey := stateToClientDeployKey(original.DeployKey, id, github.DefaultOrg, original.Name)

	return &client.GithubRepository{
		ID:           id,
		Organisation: github.DefaultOrg,
		Repository:   original.Name,
		FullName:     original.FullName,
		GitURL:       original.GitURL,
		DeployKey:    &deployKey,
	}
}

func stateToClientDeployKey(original state.DeployKey, id api.ID, organization, repository string) client.GithubDeployKey {
	return client.GithubDeployKey{
		ID:           id,
		Organisation: organization,
		Repository:   repository,
		Identifier:   original.ID,
		Title:        original.Title,
		PublicKey:    original.PublicKey,
		PrivateKeySecret: &client.GithubSecret{
			Name:    original.PrivateKeySecret.Name,
			Path:    original.PrivateKeySecret.Path,
			Version: original.PrivateKeySecret.Version,
		},
	}
}
