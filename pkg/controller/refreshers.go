package controller

import (
	"fmt"
	"path"

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

		node.ResourceState = reconciler.IdentityManagerResourceState{
			HostedZoneID: hostedZone.ID,
			Domain:       hostedZone.Domain,
		}
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

// CreateArgocdStateRefresher creates a function that gathers required runtime data for a ArgoCD resource
func CreateArgocdStateRefresher(hostedZoneFetcher HostedZoneFetcher) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		node.ResourceState = reconciler.ArgocdResourceState{
			HostedZone: hostedZoneFetcher(),
			Repository: nil,
			UserPoolID: "",
			AuthDomain: "",
		}
	}
}
