// Package v1alpha1 defines the currently active version of the okctl api
package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	// RegionEuWest1 defines the AWS region
	RegionEuWest1 = "eu-west-1"

	// OkPrincipalARNPattern defines what the Oslo kommune principal ARN for KeyCloak looks like
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/intro-structure.html#intro-structure-principal
	OkPrincipalARNPattern = "arn:aws:iam::%s:saml-provider/keycloak"
	// OkRoleARNPattern defines the standard role that we should select, i.e., the one with most privileges
	// This might be made configurable in the future, but for now it only makes sense to do it this way
	// as the other roles are pretty much useless in their current state
	OkRoleARNPattern = "arn:aws:iam::%s:role/oslokommune/iamadmin-SAML"

	// OkSamlURL is the starting point for authenticating via KeyCloak towards AWS
	OkSamlURL = "https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws"

	// ClusterConfigKind maps up the resource Kind defined by eksctl
	ClusterConfigKind = "ClusterConfig"
	// ClusterConfigAPIVersion maps up the API Version we currently use towards eksctl
	ClusterConfigAPIVersion = "eksctl.io/v1alpha5"
)

// SupportedRegions returns the supported regions on AWS
func SupportedRegions() []string {
	return []string{
		RegionEuWest1,
	}
}

// IsSupportedRegion determines if a given region is supported
func IsSupportedRegion(region string) bool {
	for _, r := range SupportedRegions() {
		if region == r {
			return true
		}
	}

	return false
}

// SupportedAvailabilityZones returns the availability zones for a
// region, given we know of it
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

// PrincipalARN returns the Ok principal ARN with account number set
func PrincipalARN(awsAccountID string) string {
	return fmt.Sprintf(OkPrincipalARNPattern, awsAccountID)
}

// RoleARN return the Ok role ARN with account number set
func RoleARN(awsAccountID string) string {
	return fmt.Sprintf(OkRoleARNPattern, awsAccountID)
}

// PermissionsBoundaryARN return the Ok permissions boundary ARN
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

// ClusterMeta comes from eksctl and maps up what we need
type ClusterMeta struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

// ClusterIAM comes from eksctl and maps up what we need
type ClusterIAM struct {
	ServiceRolePermissionsBoundary             string `json:"serviceRolePermissionsBoundary"`
	FargatePodExecutionRolePermissionsBoundary string `json:"fargatePodExecutionRolePermissionsBoundary"`
	WithOIDC                                   bool   `json:"withOIDC"`
}

// ClusterVPC comes from eksctl and maps up what we need
type ClusterVPC struct {
	ID               string           `json:"id"`
	CIDR             string           `json:"cidr"`
	ClusterEndpoints ClusterEndpoints `json:"clusterEndpoints"`
	Subnets          ClusterSubnets   `json:"subnets"`
}

// ClusterEndpoints comes from eksctl and maps up what we need
type ClusterEndpoints struct {
	PrivateAccess bool `json:"privateAccess"`
	PublicAccess  bool `json:"publicAccess"`
}

// ClusterSubnets comes from eksctl and maps up what we need
type ClusterSubnets struct {
	Private map[string]ClusterNetwork `json:"private"`
	Public  map[string]ClusterNetwork `json:"public"`
}

// ClusterNetwork comes from eksctl and maps up what we need
type ClusterNetwork struct {
	ID   string `json:"id"`
	CIDR string `json:"cidr"`
}

// NodeGroup comes from eksctl and maps up what we need
type NodeGroup struct {
}

// ClusterConfigTypeMeta returns the defaults
func ClusterConfigTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       ClusterConfigKind,
		APIVersion: ClusterConfigAPIVersion,
	}
}

// NewClusterConfig fills in all known default
// values
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

// YAML returns a serializes version of the config
func (c *ClusterConfig) YAML() ([]byte, error) {
	return yaml.Marshal(c)
}
