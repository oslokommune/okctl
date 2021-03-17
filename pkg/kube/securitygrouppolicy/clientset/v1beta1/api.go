// Package v1beta1 implements the client for the AWS EKS CRD SecurityGroupPolicy
package v1beta1

import (
	"github.com/oslokommune/okctl/pkg/kube/securitygrouppolicy/api/types/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// Interface defines client operations
type Interface interface {
	SecurityGroupPolicy(namespace string) SecurityGroupPolicyInterface
}

// Client contains the required state for interacting
// with the k8s api server
type Client struct {
	client *rest.RESTClient
}

// NewForConfig creates a REST client for interacting with the
// SecurityGroupPolicy CRD
func NewForConfig(config *rest.Config) (*Client, error) {
	err := v1beta1.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, err
	}

	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1beta1.GroupName, Version: v1beta1.GroupVersion}
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	restClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: restClient,
	}, nil
}

// SecurityGroupPolicy returns an initialised client for the given namespace
func (v *Client) SecurityGroupPolicy(namespace string) SecurityGroupPolicyInterface {
	return &securityGroupPolicyClient{
		restClient: v.client,
		ns:         namespace,
	}
}
