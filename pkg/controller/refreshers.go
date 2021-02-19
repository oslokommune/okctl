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
			PrimaryHostedZoneID: hostedZone.ID,
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

// CreateArgocdStateRefresher creates a function that gathers required runtime data for a ArgoCD resource
func CreateArgocdStateRefresher(idpFetcher IdentityPoolFetcher, hzFetcher HostedZoneFetcher) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		idp := idpFetcher()

		node.ResourceState = reconciler.ArgocdResourceState{
			HostedZone: hzFetcher(),
			UserPoolID: idp.UserPoolID,
			AuthDomain: idp.AuthDomain,
		}
	}
}

// CreateKubePromStackRefresher creates a function that gathers required runtime data for a KubePromStack resource
func CreateKubePromStackRefresher(idpFetcher IdentityPoolFetcher, hzFetecher HostedZoneFetcher) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		idp := idpFetcher()

		node.ResourceState = reconciler.KubePromStackState{
			HostedZone: hzFetecher(),
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
