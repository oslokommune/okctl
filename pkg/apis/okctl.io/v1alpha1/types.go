package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	RegionEuWest1 = "eu-west-1"

	OkPrincipalARNPattern = "arn:aws:iam::%s:saml-provider/keycloak"
	OkRoleARNPattern      = "arn:aws:iam::%s:role/oslokommune/iamadmin-SAML"

	OkSamlURL = "https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws"

	ClusterConfigKind       = "ClusterConfig"
	ClusterConfigAPIVersion = "eksctl.io/v1alpha5"
)

func SupportedRegions() []string {
	return []string{
		RegionEuWest1,
	}
}

func IsSupportedRegion(region string) bool {
	for _, r := range SupportedRegions() {
		if region == r {
			return true
		}
	}

	return false
}

func SupportedAvailabilityZones(region string) ([]string, error) {
	switch region {
	case RegionEuWest1:
		return []string{
			"eu-west-1a",
			"eu-west-1b",
			"eu-west-1c",
		}, nil
	default:
		return nil, fmt.Errorf("region: %s is not supported", region)
	}
}

func PrincipalARN(awsAccountID string) string {
	return fmt.Sprintf(OkPrincipalARNPattern, awsAccountID)
}

func RoleARN(awsAccountID string) string {
	return fmt.Sprintf(OkRoleARNPattern, awsAccountID)
}

func PermissionsBoundaryARN(awsAccountID string) string {
	return fmt.Sprintf("arn:aws:iam::%s:policy/oslokommune/oslokommune-boundary", awsAccountID)
}

// ClusterConfig is a recreation of:
// https://github.com/weaveworks/eksctl/blob/master/pkg/apis/eksctl.io/v1alpha5/types.go#L451
// where we have extract the parts that we are interested in for managing a eksctl cluster
type ClusterConfig struct {
	metav1.TypeMeta `json:",inline"`

	Metadata   ClusterMeta `json:"metadata"`
	IAM        ClusterIAM  `json:"iam"`
	VPC        ClusterVPC  `json:"vpc"`
	NodeGroups []NodeGroup `json:"nodeGroups"`
}

type ClusterMeta struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type ClusterIAM struct {
	ServiceRolePermissionsBoundary             string `json:"serviceRolePermissionsBoundary"`
	FargatePodExecutionRolePermissionsBoundary string `json:"fargatePodExecutionRolePermissionsBoundary"`
	WithOIDC                                   bool   `json:"withOIDC"`
}

type ClusterVPC struct {
	ID               string           `json:"id"`
	CIDR             string           `json:"cidr"`
	ClusterEndpoints ClusterEndpoints `json:"clusterEndpoints"`
	Subnets          ClusterSubnets   `json:"subnets"`
}

type ClusterEndpoints struct {
	PrivateAccess bool `json:"privateAccess"`
	PublicAccess  bool `json:"publicAccess"`
}

type ClusterSubnets struct {
	Private map[string]ClusterNetwork `json:"private"`
	Public  map[string]ClusterNetwork `json:"public"`
}

type ClusterNetwork struct {
	ID   string `json:"id"`
	CIDR string `json:"cidr"`
}

type NodeGroup struct {
}

func ClusterConfigTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       ClusterConfigKind,
		APIVersion: ClusterConfigAPIVersion,
	}
}

func NewClusterConfig() *ClusterConfig {
	return &ClusterConfig{
		TypeMeta: ClusterConfigTypeMeta(),
		IAM: ClusterIAM{
			WithOIDC: true,
		},
		VPC: ClusterVPC{
			ClusterEndpoints: ClusterEndpoints{
				PrivateAccess: true,
				PublicAccess:  true,
			},
			Subnets: ClusterSubnets{
				Private: map[string]ClusterNetwork{},
				Public:  map[string]ClusterNetwork{},
			},
		},
	}
}

func (c *ClusterConfig) YAML() ([]byte, error) {
	return yaml.Marshal(c)
}
