package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type IpAddressInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*IpAddressList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*IpAddress, error)
	Create(ctx context.Context, ipaddress *IpAddress) (*IpAddress, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type ipAddressClient struct {
	restClient rest.Interface
	ns         string
}

func (c *ipAddressClient) List(ctx context.Context, opts metav1.ListOptions) (*IpAddressList, error) {
	result := IpAddressList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("ipaddresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *ipAddressClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (*IpAddress, error) {
	result := IpAddress{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("ipaddresses").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *ipAddressClient) Create(ctx context.Context, ipaddress *IpAddress) (*IpAddress, error) {
	result := IpAddress{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("ipaddresses").
		Body(ipaddress).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *ipAddressClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("ipaddresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}
