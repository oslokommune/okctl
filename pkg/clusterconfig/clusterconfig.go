// Package clusterconfig knows how to construct a clusterconfiguration
package clusterconfig

import (
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/version"

	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	vpcCniVersion         = "1.11.0-eksbuild.1"
	defaultNodeVolumeSize = 80
)

// Args contains the input arguments for creating a valid
// cluster configuration
type Args struct {
	ClusterName            string
	PermissionsBoundaryARN string
	PrivateSubnets         []api.VpcSubnet
	PublicSubnets          []api.VpcSubnet
	Region                 string
	Version                string
	VpcCidr                string
	VpcID                  string
}

// New initialises the creation of a new cluster config
func New(a *Args) (*v1alpha5.ClusterConfig, error) {
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
		validation.Field(&a.Version, validation.Required),
		validation.Field(&a.VpcCidr, validation.Required),
		validation.Field(&a.VpcID, validation.Required),
	)
}

// New creates a cluster config
// nolint: funlen
func (a *Args) build() *v1alpha5.ClusterConfig {
	v := version.GetVersionInfo()

	nodeGroups := createNodeGroups(a)

	cfg := &v1alpha5.ClusterConfig{
		TypeMeta: TypeMeta(),
		Metadata: v1alpha5.ClusterMeta{
			Name:    a.ClusterName,
			Region:  a.Region,
			Version: a.Version,
			Tags: map[string]string{
				v1alpha1.OkctlVersionTag:     v.Version,
				v1alpha1.OkctlCommitTag:      v.ShortCommit,
				v1alpha1.OkctlManagedTag:     "true",
				v1alpha1.OkctlClusterNameTag: a.ClusterName,
			},
		},
		IAM: v1alpha5.ClusterIAM{
			ServiceRolePermissionsBoundary:             a.PermissionsBoundaryARN,
			FargatePodExecutionRolePermissionsBoundary: a.PermissionsBoundaryARN,
			WithOIDC: true,
		},
		Addons: []*v1alpha5.Addon{
			{
				Name: "vpc-cni",
				AttachPolicyARNs: []string{
					"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
				},
				Version:             vpcCniVersion,
				PermissionsBoundary: a.PermissionsBoundaryARN,
			},
		},
		VPC: &v1alpha5.ClusterVPC{
			ID:   a.VpcID,
			CIDR: a.VpcCidr,
			ClusterEndpoints: v1alpha5.ClusterEndpoints{
				PrivateAccess: true,
				PublicAccess:  true,
			},
			Subnets: v1alpha5.ClusterSubnets{
				Private: map[string]v1alpha5.ClusterNetwork{},
				Public:  map[string]v1alpha5.ClusterNetwork{},
			},
		},
		FargateProfiles: []v1alpha5.FargateProfile{
			{
				Name: "fp-default",
				Selectors: []v1alpha5.FargateProfileSelector{
					{Namespace: "default"},
					{Namespace: "kube-system"},
					{Namespace: "argocd"},
				},
			},
		},
		NodeGroups: nodeGroups,
		CloudWatch: &v1alpha5.ClusterCloudWatch{
			ClusterLogging: &v1alpha5.ClusterCloudWatchLogging{
				EnableTypes: v1alpha5.AllCloudWatchLogging(),
			},
		},
	}

	for _, p := range a.PublicSubnets {
		cfg.VPC.Subnets.Public[p.AvailabilityZone] = v1alpha5.ClusterNetwork{
			ID:   p.ID,
			CIDR: p.Cidr,
		}
	}

	for _, p := range a.PrivateSubnets {
		cfg.VPC.Subnets.Private[p.AvailabilityZone] = v1alpha5.ClusterNetwork{
			ID:   p.ID,
			CIDR: p.Cidr,
		}
	}

	return cfg
}

func createNodeGroups(a *Args) []v1alpha5.NodeGroup {
	availabilityZoneIds := []string{"a", "b", "c"}
	nodeGroups := make([]v1alpha5.NodeGroup, 0, len(availabilityZoneIds))
	versionWithDash := DotToDash(constant.DefaultEKSKubernetesVersion)

	for _, azID := range availabilityZoneIds {
		az := a.Region + azID

		nodeGroups = append(nodeGroups, v1alpha5.NodeGroup{
			// Name example: ng-generic-1-20-1a
			Name:         fmt.Sprintf("ng-generic-%s-1%s", versionWithDash, azID),
			InstanceType: "m5.large",
			ScalingConfig: v1alpha5.ScalingConfig{
				DesiredCapacity: 0,  //nolint: gomnd
				MinSize:         0,  //nolint: gomnd
				MaxSize:         10, //nolint: gomnd
			},
			Labels: map[string]string{
				"pool": fmt.Sprintf("ng-generic-%s", az),
			},
			Tags: map[string]string{
				"k8s.io/cluster-autoscaler/enabled":                        "true",
				fmt.Sprintf("k8s.io/cluster-autoscaler/%s", a.ClusterName): "owned",
			},
			PrivateNetworking: true,
			AvailabilityZones: []string{az},
			VolumeSize:        defaultNodeVolumeSize,
			VolumeEncrypted:   true,
		})
	}

	return nodeGroups
}

// DotToDash replaces dots with dashes
func DotToDash(str string) string {
	return strings.ReplaceAll(str, ".", "-")
}

// TypeMeta returns the defaults
func TypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       v1alpha5.ClusterConfigKind,
		APIVersion: v1alpha5.ClusterConfigAPIVersion,
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
	PolicyArns             []string
	Region                 string
}

// NewServiceAccount returns an initialised cluster config for creating a service account
// with an associated IAM managed policy
func NewServiceAccount(a *ServiceAccountArgs) (*v1alpha5.ClusterConfig, error) {
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
		validation.Field(&a.PolicyArns, validation.Required),
		validation.Field(&a.Region, validation.Required),
	)
}

func (a *ServiceAccountArgs) build() *v1alpha5.ClusterConfig {
	v := version.GetVersionInfo()

	return &v1alpha5.ClusterConfig{
		TypeMeta: TypeMeta(),
		Metadata: v1alpha5.ClusterMeta{
			Name:   a.ClusterName,
			Region: a.Region,
			Tags: map[string]string{
				v1alpha1.OkctlVersionTag:     v.Version,
				v1alpha1.OkctlCommitTag:      v.ShortCommit,
				v1alpha1.OkctlManagedTag:     "true",
				v1alpha1.OkctlClusterNameTag: a.ClusterName,
			},
		},
		IAM: v1alpha5.ClusterIAM{
			WithOIDC: true,
			ServiceAccounts: []*v1alpha5.ClusterIAMServiceAccount{
				{
					ClusterIAMMeta: v1alpha5.ClusterIAMMeta{
						Name:      a.Name,
						Namespace: a.Namespace,
						Labels:    a.Labels,
					},
					AttachPolicyARNs:    a.PolicyArns,
					PermissionsBoundary: a.PermissionsBoundaryArn,
				},
			},
		},
	}
}

// NewExternalSecretsServiceAccount returns an initialised configuration for
// creating an external secrets service account
func NewExternalSecretsServiceAccount(clusterName, region, permissionsBoundaryArn string, policyArns []string) (*v1alpha5.ClusterConfig, error) {
	return NewServiceAccount(&ServiceAccountArgs{
		ClusterName: clusterName,
		Labels: map[string]string{
			"aws-usage": "cluster-ops",
		},
		Name:                   "external-secrets",
		Namespace:              "kube-system",
		PermissionsBoundaryArn: permissionsBoundaryArn,
		PolicyArns:             policyArns,
		Region:                 region,
	})
}

// NewAlbIngressControllerServiceAccount returns an initialised configuration
// for creating an aws-alb-ingress-controller service account
func NewAlbIngressControllerServiceAccount(clusterName, region, permissionsBoundaryArn string, policyArns []string) (*v1alpha5.ClusterConfig, error) {
	return NewServiceAccount(&ServiceAccountArgs{
		ClusterName: clusterName,
		Labels: map[string]string{
			"aws-usage": "cluster-ops",
		},
		Name:                   "alb-ingress-controller",
		Namespace:              "kube-system",
		PermissionsBoundaryArn: permissionsBoundaryArn,
		PolicyArns:             policyArns,
		Region:                 region,
	})
}

// NewAWSLoadBalancerControllerServiceAccount returns an initialised configuration
// for creating an aws-load-balancer-controller service account
func NewAWSLoadBalancerControllerServiceAccount(clusterName, region, permissionsBoundaryArn string, policyArns []string) (*v1alpha5.ClusterConfig, error) {
	return NewServiceAccount(&ServiceAccountArgs{
		ClusterName: clusterName,
		Labels: map[string]string{
			"aws-usage": "cluster-ops",
		},
		Name:                   "aws-load-balancer-controller",
		Namespace:              "kube-system",
		PermissionsBoundaryArn: permissionsBoundaryArn,
		PolicyArns:             policyArns,
		Region:                 region,
	})
}

// NewExternalDNSServiceAccount returns an initialised configuration
// for creating an external-dns service account
func NewExternalDNSServiceAccount(clusterName, region, permissionsBoundaryArn string, policyArns []string) (*v1alpha5.ClusterConfig, error) {
	return NewServiceAccount(&ServiceAccountArgs{
		ClusterName: clusterName,
		Labels: map[string]string{
			"aws-usage": "cluster-ops",
		},
		Name:                   "external-dns",
		Namespace:              "kube-system",
		PermissionsBoundaryArn: permissionsBoundaryArn,
		PolicyArns:             policyArns,
		Region:                 region,
	})
}

// NewAutoscalerServiceAccount returns an initialised configuration
// for creating a cluster autoscaler service account
func NewAutoscalerServiceAccount(clusterName, region, permissionsBoundaryArn string, policyArns []string) (*v1alpha5.ClusterConfig, error) {
	return NewServiceAccount(&ServiceAccountArgs{
		ClusterName: clusterName,
		Labels: map[string]string{
			"aws-usage": "cluster-ops",
		},
		Name:                   "autoscaler",
		Namespace:              "kube-system",
		PermissionsBoundaryArn: permissionsBoundaryArn,
		PolicyArns:             policyArns,
		Region:                 region,
	})
}

// NewBlockstorageServiceAccount returns an initialised configuration
// for creating a cluster Blockstorage service account
func NewBlockstorageServiceAccount(clusterName, region, permissionsBoundaryArn string, policyArns []string) (*v1alpha5.ClusterConfig, error) {
	return NewServiceAccount(&ServiceAccountArgs{
		ClusterName: clusterName,
		Labels: map[string]string{
			"aws-usage": "cluster-ops",
		},
		Name:                   "blockstorage",
		Namespace:              "kube-system",
		PermissionsBoundaryArn: permissionsBoundaryArn,
		PolicyArns:             policyArns,
		Region:                 region,
	})
}

// NewCloudwatchDatasourceServiceAccount returns an initialised configuration
// for creating a cluster CloudwatchDatasource service account
func NewCloudwatchDatasourceServiceAccount(clusterName, region, namespace, permissionsBoundaryArn string, policyArns []string) (*v1alpha5.ClusterConfig, error) {
	return NewServiceAccount(&ServiceAccountArgs{
		ClusterName: clusterName,
		Labels: map[string]string{
			"aws-usage": "cluster-ops",
		},
		Name:                   "cloudwatch-datasource",
		Namespace:              namespace,
		PermissionsBoundaryArn: permissionsBoundaryArn,
		PolicyArns:             policyArns,
		Region:                 region,
	})
}

// NewLokiServiceAccount returns an initialized configuration
// for creating a service account used by Loki to manipulate S3 and DynamoDB
func NewLokiServiceAccount(clusterName, region, namespace, permissionsBoundaryARN string, policyARNs []string) (*v1alpha5.ClusterConfig, error) {
	return NewServiceAccount(&ServiceAccountArgs{
		ClusterName: clusterName,
		Labels: map[string]string{
			"aws-usage": "cluster-ops",
		},
		Name:                   "loki",
		Namespace:              namespace,
		PermissionsBoundaryArn: permissionsBoundaryARN,
		PolicyArns:             policyARNs,
		Region:                 region,
	})
}
