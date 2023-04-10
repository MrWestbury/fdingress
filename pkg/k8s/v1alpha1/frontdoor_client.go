package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type FrontdoorInterface interface {
	List(opts metav1.ListOptions) (*FrontdoorList, error)
	Get(name string, options metav1.GetOptions) (*Frontdoor, error)
	Create(*Frontdoor) (*Frontdoor, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type frontdoorClient struct {
	restClient rest.Interface
	ns         string
}

func (c *frontdoorClient) List(opts metav1.ListOptions) (*FrontdoorList, error) {
	result := FrontdoorList{}
	ctx := context.TODO()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("frontdoors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *frontdoorClient) Get(name string, opts metav1.GetOptions) (*Frontdoor, error) {
	result := Frontdoor{}
	ctx := context.TODO()
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("frontdoors").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *frontdoorClient) Create(frontdoor *Frontdoor) (*Frontdoor, error) {
	result := Frontdoor{}
	ctx := context.TODO()
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("frontdoors").
		Body(frontdoor).
		Do(ctx).
		Into(&result)

	return &result, err
}

func (c *frontdoorClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	ctx := context.TODO()
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("frontdoors").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}
