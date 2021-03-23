// Package v1beta1 implements the AWS EKS CRD SecurityGroupPolicy
// - https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html
//go:generate controller-gen object paths=$GOFILE
package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// SecurityGroupPolicySpec defines the spec content
//+kubebuilder:object:generate:=true
type SecurityGroupPolicySpec struct {
	PodSelector    SecurityGroupPolicyPodSelector    `json:"podSelector"`
	SecurityGroups SecurityGroupPolicySecurityGroups `json:"securityGroups"`
}

// SecurityGroupPolicyPodSelector adds the content of pod selector
//+kubebuilder:object:generate:=true
type SecurityGroupPolicyPodSelector struct {
	MatchLabels map[string]string `json:"matchLabels"`
}

// SecurityGroupPolicySecurityGroups adds the content of security group ids
//+kubebuilder:object:generate:=true
type SecurityGroupPolicySecurityGroups struct {
	GroupIDs []string `json:"groupIds"`
}

// SecurityGroupPolicy is the root of the CRD
//+kubebuilder:object:root:=true
//+kubebuilder:object:generate:=true
type SecurityGroupPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SecurityGroupPolicySpec `json:"spec"`
}

// SecurityGroupPolicyList is a list of the CRDs
//+kubebuilder:object:root:=true
//+kubebuilder:object:generate:=true
type SecurityGroupPolicyList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Items []SecurityGroupPolicy `json:"items"`
}
