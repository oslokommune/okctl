// Package v1beta1 implements the types for security group policy
// provided by AWS EKS CRD
// nolint: gochecknoglobals
package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// GroupName for the CRD
	GroupName = "vpcresources.k8s.aws"
	// GroupVersion for the CRD
	GroupVersion = "v1beta1"
)

// SchemeGroupVersion ...
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

var (
	// SchemeBuilder ...
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme ...
	AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&SecurityGroupPolicy{},
		&SecurityGroupPolicyList{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}
