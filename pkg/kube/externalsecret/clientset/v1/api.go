package v1

import (
	v1 "github.com/oslokommune/okctl/pkg/kube/externalsecret/clientset/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type Interface interface {
	ExternalSecrets(namespace string) v1.ExternalSecretInterface
}

type Client struct {
	client *rest.RESTClient
}

func NewForConfig(config *rest.Config) (*Client, error) {
	err := AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, err
	}
	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
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

func (v *Client) ExternalSecrets(namespace string) v1.ExternalSecretInterface {
	return &v1.externalSecretClient{
		restClient: v.client,
		ns:         namespace,
	}
}
