// nolint stylecheck
package v1alpha5

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const (
	// ClusterConfigKind maps up the resource Kind defined by eksctl
	ClusterConfigKind = "ClusterConfig"
	// ClusterConfigAPIVersion maps up the API Version we currently use towards eksctl
	ClusterConfigAPIVersion = "eksctl.io/v1alpha5"
)

// ClusterConfig is a recreation of:
// https://github.com/weaveworks/eksctl/blob/master/pkg/apis/eksctl.io/v1alpha5/types.go
// where we have extract the parts that we are interested in for managing a eksctl cluster
type ClusterConfig struct {
	metav1.TypeMeta `json:",inline"`

	Metadata        ClusterMeta      `json:"metadata"`
	IAM             ClusterIAM       `json:"iam"`
	VPC             *ClusterVPC      `json:"vpc,omitempty"`
	FargateProfiles []FargateProfile `json:"fargateProfiles,omitempty"`
	NodeGroups      []NodeGroup      `json:"nodeGroups,omitempty"`
	Status          *ClusterStatus   `json:"status,omitempty"`
}

// ClusterStatus hold read-only attributes of a cluster
type ClusterStatus struct {
	Endpoint                 string `json:"endpoint,omitempty"`
	CertificateAuthorityData []byte `json:"certificateAuthorityData,omitempty"`
	ARN                      string `json:"arn,omitempty"`
	StackName                string `json:"stackName,omitempty"`
}

// ClusterMeta comes from eksctl and maps up what we need
type ClusterMeta struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

func (c *ClusterMeta) String() string {
	return fmt.Sprintf("%s.%s.eksctl.io", c.Name, c.Region)
}

// ClusterIAM comes from eksctl and maps up what we need
type ClusterIAM struct {
	ServiceRolePermissionsBoundary             string                      `json:"serviceRolePermissionsBoundary,omitempty"`
	FargatePodExecutionRolePermissionsBoundary string                      `json:"fargatePodExecutionRolePermissionsBoundary,omitempty"`
	WithOIDC                                   bool                        `json:"withOIDC"`
	ServiceAccounts                            []*ClusterIAMServiceAccount `json:"serviceAccounts,omitempty"`
}

// ClusterIAMMeta holds information we can use to create ObjectMeta for service
// accounts
type ClusterIAMMeta struct {
	Name      string            `json:"name,omitempty"`
	Namespace string            `json:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// ClusterIAMServiceAccount comes from eksctl and maps up what we need
type ClusterIAMServiceAccount struct {
	ClusterIAMMeta      `json:"metadata,omitempty"`
	AttachPolicyARNs    []string `json:"attachPolicyARNs"`
	PermissionsBoundary string   `json:"permissionsBoundary"`
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

// FargateProfile comes from eksctl and maps up what we need
type FargateProfile struct {
	Name      string                   `json:"name"`
	Selectors []FargateProfileSelector `json:"selectors"`
}

// FargateProfileSelector comes from eksctl and maps up what we need
type FargateProfileSelector struct {
	Namespace string `json:"namespace"`
}

// NodeGroup comes from eksctl and maps up what we need
type NodeGroup struct {
	Name              string            `json:"name"`
	InstanceType      string            `json:"instanceType"`
	Labels            map[string]string `json:"labels"`
	Tags              map[string]string `json:"tags"`
	PrivateNetworking bool              `json:"privateNetworking"`
	IAM               NodeGroupIAM      `json:"iam"`

	ScalingConfig `json:",inline"`
}

// NodeGroupIAM comes from eksctl and maps up what we need
type NodeGroupIAM struct {
	InstanceRolePermissionsBoundary string `json:"instanceRolePermissionsBoundary"`
}

// ScalingConfig comes from eksctl and maps up what we need
type ScalingConfig struct {
	DesiredCapacity int `json:"desiredCapacity"`
	MinSize         int `json:"minSize"`
	MaxSize         int `json:"maxSize"`
}

// YAML returns a serializes version of the config
func (c *ClusterConfig) YAML() ([]byte, error) {
	return yaml.Marshal(c)
}
