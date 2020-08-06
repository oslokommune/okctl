// Package clusterconfig knows how to construct a clusterconfiguration
package clusterconfig

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PermissionsBoundary sets the ARN of the permissions boundary
type PermissionsBoundary interface {
	PermissionsBoundary(arn string) Region
}

// Region sets the AWS region
type Region interface {
	Region(region string) Vpc
}

// Vpc sets the vpc id and cidr
type Vpc interface {
	Vpc(id, cidr string) Subnets
}

// Subnets sets the public and private subnets
type Subnets interface {
	Subnets(public, private []api.VpcSubnet) Build
}

// Build creates a cluster config using the given args
type Build interface {
	Build() *api.ClusterConfig
}

type args struct {
	clusterName            string
	region                 string
	permissionsBoundaryARN string
	vpcID                  string
	vpcCidr                string
	publicSubnets          []api.VpcSubnet
	privateSubnets         []api.VpcSubnet
}

// New initialises the creation of a new cluster config
func New(clusterName string) PermissionsBoundary {
	return &args{
		clusterName: clusterName,
	}
}

// PermissionsBoundary sets the AWS IAM permissions boundary
func (a *args) PermissionsBoundary(arn string) Region {
	a.permissionsBoundaryARN = arn

	return a
}

// Region sets the AWS region
func (a *args) Region(region string) Vpc {
	a.region = region

	return a
}

// Vpc sets the vpc id and cidr
func (a *args) Vpc(id, cidr string) Subnets {
	a.vpcID = id
	a.vpcCidr = cidr

	return a
}

// Subnets sets the private and public subnets
func (a *args) Subnets(public, private []api.VpcSubnet) Build {
	a.publicSubnets = public
	a.privateSubnets = private

	return a
}

// TypeMeta returns the defaults
func TypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       api.ClusterConfigKind,
		APIVersion: api.ClusterConfigAPIVersion,
	}
}

// New creates a cluster config
// nolint: funlen
func (a *args) Build() *api.ClusterConfig {
	cfg := &api.ClusterConfig{
		TypeMeta: TypeMeta(),
		Metadata: api.ClusterMeta{
			Name:   a.clusterName,
			Region: a.region,
		},
		IAM: api.ClusterIAM{
			ServiceRolePermissionsBoundary:             a.permissionsBoundaryARN,
			FargatePodExecutionRolePermissionsBoundary: a.permissionsBoundaryARN,
			WithOIDC: true,
		},
		FargateProfiles: []api.FargateProfile{
			{
				Name: "fp-default",
				Selectors: []api.FargateProfileSelector{
					{Namespace: "default"},
					{Namespace: "kube-system"},
				},
			},
		},
		NodeGroups: []api.NodeGroup{
			{
				Name:         "ng-generic",
				InstanceType: "m5.large",
				ScalingConfig: api.ScalingConfig{
					DesiredCapacity: 2, //nolint: gomnd
					MinSize:         1,
					MaxSize:         10, //nolint: gomnd
				},
				Labels: map[string]string{
					"pool": "ng-generic",
				},
				Tags: map[string]string{
					"k8s.io/cluster-autoscaler/enabled":                        "true",
					fmt.Sprintf("k8s.io/cluster-autoscaler/%s", a.clusterName): "owned",
				},
				PrivateNetworking: true,
				IAM: api.NodeGroupIAM{
					InstanceRolePermissionsBoundary: a.permissionsBoundaryARN,
				},
			},
		},
		VPC: api.ClusterVPC{
			ClusterEndpoints: api.ClusterEndpoints{
				PrivateAccess: true,
				PublicAccess:  true,
			},
			Subnets: api.ClusterSubnets{
				Private: map[string]api.ClusterNetwork{},
				Public:  map[string]api.ClusterNetwork{},
			},
		},
	}

	for _, p := range a.publicSubnets {
		cfg.VPC.Subnets.Public[p.AvailabilityZone] = api.ClusterNetwork{
			ID:   p.ID,
			CIDR: p.Cidr,
		}
	}

	for _, p := range a.privateSubnets {
		cfg.VPC.Subnets.Private[p.AvailabilityZone] = api.ClusterNetwork{
			ID:   p.ID,
			CIDR: p.Cidr,
		}
	}

	return cfg
}
