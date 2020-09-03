// nolint
package v1

import (
	v1 "github.com/oslokommune/okctl/pkg/kube/externalsecret/api/types/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type Interface interface {
	ExternalSecrets(namespace string) ExternalSecretInterface
}

type Client struct {
	client *rest.RESTClient
}

func NewForConfig(config *rest.Config) (*Client, error) {
	err := v1.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, err
	}
	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1.GroupName, Version: v1.GroupVersion}
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

func (v *Client) ExternalSecrets(namespace string) ExternalSecretInterface {
	return &externalSecretClient{
		restClient: v.client,
		ns:         namespace,
	}
}
