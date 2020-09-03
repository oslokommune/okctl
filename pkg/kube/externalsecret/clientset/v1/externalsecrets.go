// nolint
package v1

import (
	"context"

	v12 "github.com/oslokommune/okctl/pkg/kube/externalsecret/api/types/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// For implementing more of these, see:
// - https://github.com/kubernetes/client-go/blob/master/kubernetes/typed/core/v1/pod.go
type ExternalSecretInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v12.ExternalSecretList, error)
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v12.ExternalSecret, error)
	Create(ctx context.Context, secret *v12.ExternalSecret) (*v12.ExternalSecret, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
}

type externalSecretClient struct {
	restClient *rest.RESTClient
	ns         string
}

func (e *externalSecretClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return e.restClient.Delete().
		Namespace(e.ns).
		Resource("externalsecrets").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (e *externalSecretClient) List(ctx context.Context, opts metav1.ListOptions) (*v12.ExternalSecretList, error) {
	result := v12.ExternalSecretList{}
	err := e.restClient.
		Get().
		Namespace(e.ns).
		Resource("externalSecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (e *externalSecretClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v12.ExternalSecret, error) {
	result := v12.ExternalSecret{}
	err := e.restClient.
		Get().
		Namespace(e.ns).
		Resource("externalSecrets").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (e *externalSecretClient) Create(ctx context.Context, secret *v12.ExternalSecret) (*v12.ExternalSecret, error) {
	result := v12.ExternalSecret{}
	err := e.restClient.
		Post().
		Namespace(e.ns).
		Resource("externalSecrets").
		Body(secret).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (e *externalSecretClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return e.restClient.
		Get().
		Namespace(e.ns).
		Resource("externalSecrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}
