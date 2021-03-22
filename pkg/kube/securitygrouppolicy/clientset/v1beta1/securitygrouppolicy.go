// Package v1beta1 provides the allowed operations on the CRD
package v1beta1

import (
	"context"

	v1beta1types "github.com/oslokommune/okctl/pkg/kube/securitygrouppolicy/api/types/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const (
	awsSecurityGroupPolicyCRDResourceName = "securitygrouppolicies"
)

// SecurityGroupPolicyInterface enumerates the allowed operations
// For implementing more of these, see:
// - https://github.com/kubernetes/client-go/blob/master/kubernetes/typed/core/v1/pod.go
type SecurityGroupPolicyInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1types.SecurityGroupPolicyList, error)
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1types.SecurityGroupPolicy, error)
	Create(ctx context.Context, secret *v1beta1types.SecurityGroupPolicy) (*v1beta1types.SecurityGroupPolicy, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

type securityGroupPolicyClient struct {
	restClient *rest.RESTClient
	ns         string
}

// Delete a security group policy
func (e *securityGroupPolicyClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return e.restClient.Delete().
		Namespace(e.ns).
		Resource(awsSecurityGroupPolicyCRDResourceName).
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// List the security group policies
func (e *securityGroupPolicyClient) List(ctx context.Context, opts metav1.ListOptions) (*v1beta1types.SecurityGroupPolicyList, error) {
	result := v1beta1types.SecurityGroupPolicyList{}
	err := e.restClient.
		Get().
		Namespace(e.ns).
		Resource(awsSecurityGroupPolicyCRDResourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

// Get a security group policy
func (e *securityGroupPolicyClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1beta1types.SecurityGroupPolicy, error) {
	result := v1beta1types.SecurityGroupPolicy{}
	err := e.restClient.
		Get().
		Namespace(e.ns).
		Resource(awsSecurityGroupPolicyCRDResourceName).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

// Create a security group policy
func (e *securityGroupPolicyClient) Create(ctx context.Context, policy *v1beta1types.SecurityGroupPolicy) (*v1beta1types.SecurityGroupPolicy, error) {
	result := v1beta1types.SecurityGroupPolicy{}
	err := e.restClient.
		Post().
		Namespace(e.ns).
		Resource(awsSecurityGroupPolicyCRDResourceName).
		Body(policy).
		Do(ctx).
		Into(&result)

	return &result, err
}

// Watch a security group policy
func (e *securityGroupPolicyClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true

	return e.restClient.
		Get().
		Namespace(e.ns).
		Resource(awsSecurityGroupPolicyCRDResourceName).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}
