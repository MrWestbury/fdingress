package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type FdIngressV1Alpha1Interface interface {
	Frontdoors(namespace string) FrontdoorInterface
	IpAddresses(namespace string) IpAddressInterface
}

type FdIngressV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*FdIngressV1Alpha1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &FdIngressV1Alpha1Client{restClient: client}, nil
}

func (c *FdIngressV1Alpha1Client) Frontdoors(namespace string) FrontdoorInterface {
	return &frontdoorClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *FdIngressV1Alpha1Client) IpAddresses(namespace string) IpAddressInterface {
	return &ipAddressClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
