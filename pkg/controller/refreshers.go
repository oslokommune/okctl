package controller

import (
	"errors"
	"fmt"

	"github.com/asdine/storm/v3"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/controller/reconciler"
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// CreateClusterStateRefresher creates a function that gathers required runtime data for a cluster resource
func CreateClusterStateRefresher(id api.ID, vpc client.VPCState) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		v, err := vpc.GetVpc(cfn.NewStackNamer().Vpc(id.ClusterName))
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Errorf("getting vpc state: %w", err))
		}

		if v == nil {
			node.ResourceState = reconciler.ClusterResourceState{}

			return
		}

		node.ResourceState = reconciler.ClusterResourceState{VPC: *v}
	}
}

// CreateAWSLoadBalancerControllerRefresher creates a function that gathers required runtime data for AWS
// load balancer controller
func CreateAWSLoadBalancerControllerRefresher(id api.ID, vpc client.VPCState) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		v, err := vpc.GetVpc(cfn.NewStackNamer().Vpc(id.ClusterName))
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Errorf("getting vpc state: %w", err))
		}

		if v == nil {
			node.ResourceState = reconciler.AWSLoadBalancerControllerResourceState{}

			return
		}

		node.ResourceState = reconciler.AWSLoadBalancerControllerResourceState{
			VpcID: v.VpcID,
		}
	}
}

// CreateExternalDNSStateRefresher creates a function that gathers required runtime data for a External DNS resource
func CreateExternalDNSStateRefresher(domainState client.DomainState) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		d, err := domainState.GetPrimaryHostedZone()
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Errorf("getting primary hosted zone: %w", err))
		}

		if d == nil {
			node.ResourceState = reconciler.ExternalDNSResourceState{}

			return
		}

		node.ResourceState = reconciler.ExternalDNSResourceState{
			PrimaryHostedZoneID: d.HostedZoneID,
		}
	}
}

// CreateIdentityManagerRefresher creates a function that gathers required runtime data for a Identity Manager resource
func CreateIdentityManagerRefresher(domainState client.DomainState) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		d, err := domainState.GetPrimaryHostedZone()
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Errorf("getting primary hosted zone: %w", err))
		}

		if d == nil {
			node.ResourceState = reconciler.IdentityManagerResourceState{}

			return
		}

		node.ResourceState = reconciler.IdentityManagerResourceState{
			HostedZoneID: d.HostedZoneID,
		}
	}
}

// CreateArgocdStateRefresher creates a function that gathers required runtime data for a ArgoCD resource
func CreateArgocdStateRefresher(id api.ID, domainState client.DomainState, managerState client.IdentityManagerState) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		d, err := domainState.GetPrimaryHostedZone()
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Errorf("getting primary hosted zone: %w", err))
		}

		node.ResourceState = reconciler.ArgocdResourceState{}

		if d == nil {
			return
		}

		i, err := managerState.GetIdentityPool(cfn.NewStackNamer().IdentityPool(id.ClusterName))
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Errorf("getting identity pool: %w", err))
		}

		if i == nil {
			return
		}

		node.ResourceState = reconciler.ArgocdResourceState{
			HostedZone: d,
			UserPoolID: i.UserPoolID,
			AuthDomain: i.AuthDomain,
		}
	}
}

// CreateKubePromStackRefresher creates a function that gathers required runtime data for a KubePromStack resource
func CreateKubePromStackRefresher(id api.ID, domainState client.DomainState, managerState client.IdentityManagerState) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		d, err := domainState.GetPrimaryHostedZone()
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Errorf("getting primary hosted zone: %w", err))
		}

		node.ResourceState = reconciler.KubePromStackState{}

		if d == nil {
			return
		}

		i, err := managerState.GetIdentityPool(cfn.NewStackNamer().IdentityPool(id.ClusterName))
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Errorf("getting identity pool: %w", err))
		}

		if i == nil {
			return
		}

		node.ResourceState = reconciler.KubePromStackState{
			HostedZone: d,
			UserPoolID: i.UserPoolID,
			AuthDomain: i.AuthDomain,
		}
	}
}

// CreateUsersRefresher creates a function that gathers required runtime data for a Users resource
func CreateUsersRefresher(id api.ID, managerState client.IdentityManagerState) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		i, err := managerState.GetIdentityPool(cfn.NewStackNamer().IdentityPool(id.ClusterName))
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Errorf("getting identity pool: %w", err))
		}

		if i == nil {
			node.ResourceState = reconciler.UsersState{}

			return
		}

		node.ResourceState = reconciler.UsersState{
			UserPoolID: i.UserPoolID,
		}
	}
}

// CreateNameserverDelegationStateRefresher creates a function that gathers required runtime data for a nameserver delegation
// request
func CreateNameserverDelegationStateRefresher(domainState client.DomainState) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		d, err := domainState.GetPrimaryHostedZone()
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Errorf("getting primary hosted zone: %w", err))
		}

		if d == nil {
			node.ResourceState = reconciler.NameserverHandlerReconcilerResourceState{}

			return
		}

		node.ResourceState = reconciler.NameserverHandlerReconcilerResourceState{
			PrimaryHostedZoneFQDN: d.FQDN,
			Nameservers:           d.NameServers,
		}
	}
}

// CreatePostgresDatabasesRefresher creates a function that gathers required runtime data
func CreatePostgresDatabasesRefresher(id api.ID, vpc client.VPCState) resourcetree.StateRefreshFn {
	return func(node *resourcetree.ResourceNode) {
		v, err := vpc.GetVpc(cfn.NewStackNamer().Vpc(id.ClusterName))
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			panic(fmt.Errorf("getting vpc state: %w", err))
		}

		if v == nil {
			node.ResourceState = reconciler.PostgresState{}

			return
		}

		ids := make([]string, len(v.DatabaseSubnets))
		cidrs := make([]string, len(v.DatabaseSubnets))

		for i, s := range v.DatabaseSubnets {
			ids[i] = s.ID
			cidrs[i] = s.Cidr
		}

		node.ResourceState = reconciler.PostgresState{
			VpcID:             v.VpcID,
			DBSubnetGroupName: v.DatabaseSubnetsGroupName,
			DBSubnetIDs:       ids,
			DBSubnetCIDRs:     cidrs,
		}
	}
}
