// Package clusterconfig knows how to construct a clusterconfiguration
package clusterconfig

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Args contains the input arguments for creating a valid
// cluster configuration
type Args struct {
	ClusterName            string
	PermissionsBoundaryARN string
	PrivateSubnets         []api.VpcSubnet
	PublicSubnets          []api.VpcSubnet
	Region                 string
	VpcCidr                string
	VpcID                  string
}

// New initialises the creation of a new cluster config
func New(a *Args) (*v1alpha1.ClusterConfig, error) {
	err := a.validate()
	if err != nil {
		return nil, err
	}

	return a.build(), nil
}

func (a *Args) validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ClusterName, validation.Required),
		validation.Field(&a.PermissionsBoundaryARN, validation.Required),
		validation.Field(&a.PrivateSubnets, validation.Required),
		validation.Field(&a.PublicSubnets, validation.Required),
		validation.Field(&a.Region, validation.Required),
		validation.Field(&a.VpcCidr, validation.Required),
		validation.Field(&a.VpcID, validation.Required),
	)
}

// New creates a cluster config
// nolint: funlen
func (a *Args) build() *v1alpha1.ClusterConfig {
	cfg := &v1alpha1.ClusterConfig{
		TypeMeta: TypeMeta(),
		Metadata: v1alpha1.ClusterMeta{
			Name:   a.ClusterName,
			Region: a.Region,
		},
		IAM: v1alpha1.ClusterIAM{
			ServiceRolePermissionsBoundary:             a.PermissionsBoundaryARN,
			FargatePodExecutionRolePermissionsBoundary: a.PermissionsBoundaryARN,
			WithOIDC: true,
		},
		FargateProfiles: []v1alpha1.FargateProfile{
			{
				Name: "fp-default",
				Selectors: []v1alpha1.FargateProfileSelector{
					{Namespace: "default"},
					{Namespace: "kube-system"},
				},
			},
		},
		NodeGroups: []v1alpha1.NodeGroup{
			{
				Name:         "ng-generic",
				InstanceType: "m5.large",
				ScalingConfig: v1alpha1.ScalingConfig{
					DesiredCapacity: 2, //nolint: gomnd
					MinSize:         1,
					MaxSize:         10, //nolint: gomnd
				},
				Labels: map[string]string{
					"pool": "ng-generic",
				},
				Tags: map[string]string{
					"k8s.io/cluster-autoscaler/enabled":                        "true",
					fmt.Sprintf("k8s.io/cluster-autoscaler/%s", a.ClusterName): "owned",
				},
				PrivateNetworking: true,
				IAM: v1alpha1.NodeGroupIAM{
					InstanceRolePermissionsBoundary: a.PermissionsBoundaryARN,
				},
			},
		},
		VPC: &v1alpha1.ClusterVPC{
			ID:   a.VpcID,
			CIDR: a.VpcCidr,
			ClusterEndpoints: v1alpha1.ClusterEndpoints{
				PrivateAccess: true,
				PublicAccess:  true,
			},
			Subnets: v1alpha1.ClusterSubnets{
				Private: map[string]v1alpha1.ClusterNetwork{},
				Public:  map[string]v1alpha1.ClusterNetwork{},
			},
		},
	}

	for _, p := range a.PublicSubnets {
		cfg.VPC.Subnets.Public[p.AvailabilityZone] = v1alpha1.ClusterNetwork{
			ID:   p.ID,
			CIDR: p.Cidr,
		}
	}

	for _, p := range a.PrivateSubnets {
		cfg.VPC.Subnets.Private[p.AvailabilityZone] = v1alpha1.ClusterNetwork{
			ID:   p.ID,
			CIDR: p.Cidr,
		}
	}

	return cfg
}

// TypeMeta returns the defaults
func TypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       v1alpha1.ClusterConfigKind,
		APIVersion: v1alpha1.ClusterConfigAPIVersion,
	}
}

// ServiceAccountArgs contains the arguments for creating a valid
// service account
type ServiceAccountArgs struct {
	ClusterName            string
	Labels                 map[string]string
	Name                   string
	Namespace              string
	PermissionsBoundaryArn string
	PolicyArn              string
	Region                 string
}

// NewServiceAccount returns an initialised cluster config for creating a service account
// with an associated IAM managed policy
func NewServiceAccount(a *ServiceAccountArgs) (*v1alpha1.ClusterConfig, error) {
	err := a.validate()
	if err != nil {
		return nil, err
	}

	return a.build(), nil
}

func (a *ServiceAccountArgs) validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ClusterName, validation.Required),
		validation.Field(&a.Labels, validation.Required),
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Namespace, validation.Required),
		validation.Field(&a.PermissionsBoundaryArn, validation.Required),
		validation.Field(&a.PolicyArn, validation.Required),
		validation.Field(&a.Region, validation.Required),
	)
}

func (a *ServiceAccountArgs) build() *v1alpha1.ClusterConfig {
	return &v1alpha1.ClusterConfig{
		TypeMeta: TypeMeta(),
		Metadata: v1alpha1.ClusterMeta{
			Name:   a.ClusterName,
			Region: a.Region,
		},
		IAM: v1alpha1.ClusterIAM{
			WithOIDC: true,
			ServiceAccounts: []*v1alpha1.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: v1alpha1.ClusterIAMMeta{
						Name:      a.Name,
						Namespace: a.Namespace,
						Labels:    a.Labels,
					},
					AttachPolicyARNs: []string{
						a.PolicyArn,
					},
					PermissionsBoundary: a.PermissionsBoundaryArn,
				},
			},
		},
	}
}

// NewExternalSecretsServiceAccount returns an initialised configuration for
// creating an external secrets service account
func NewExternalSecretsServiceAccount(clusterName, region, policyArn, permissionsBoundaryArn string) (*v1alpha1.ClusterConfig, error) {
	return NewServiceAccount(&ServiceAccountArgs{
		ClusterName: clusterName,
		Labels: map[string]string{
			"aws-usage": "cluster-ops",
		},
		Name:                   "external-secrets",
		Namespace:              "kube-system",
		PermissionsBoundaryArn: permissionsBoundaryArn,
		PolicyArn:              policyArn,
		Region:                 region,
	})
}

// NewAlbIngressControllerServiceAccount returns an initialised configuration
// for creating an aws-alb-ingress-controller service account
func NewAlbIngressControllerServiceAccount(clusterName, region, policyArn, permissionsBoundaryArn string) (*v1alpha1.ClusterConfig, error) {
	return NewServiceAccount(&ServiceAccountArgs{
		ClusterName: clusterName,
		Labels: map[string]string{
			"aws-usage": "cluster-ops",
		},
		Name:                   "alb-ingress-controller",
		Namespace:              "kube-system",
		PermissionsBoundaryArn: permissionsBoundaryArn,
		PolicyArn:              policyArn,
		Region:                 region,
	})
}

// NewExternalDNSServiceAccount returns an initialised configuration
// for creating an external-dns service account
func NewExternalDNSServiceAccount(clusterName, region, policyArn, permissionsBoundaryArn string) (*v1alpha1.ClusterConfig, error) {
	return NewServiceAccount(&ServiceAccountArgs{
		ClusterName: clusterName,
		Labels: map[string]string{
			"aws-usage": "cluster-ops",
		},
		Name:                   "external-dns",
		Namespace:              "kube-system",
		PermissionsBoundaryArn: permissionsBoundaryArn,
		PolicyArn:              policyArn,
		Region:                 region,
	})
}
