package k8s

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/mrwestbury/frontdoor-ingress/pkg/k8s/v1alpha1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Scanner struct {
	ns             string
	client         *kubernetes.Clientset
	fdclient       v1alpha1.FdIngressV1Alpha1Interface
	ingressStore   cache.Store
	frontdoorStore cache.Store
	ipaddressStore cache.Store
}

func NewScanner() *Scanner {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fdClientset, err := v1alpha1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	scanner := &Scanner{
		client:   clientset,
		fdclient: fdClientset,
	}

	fd, err := os.Open("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		log.Fatalf("failed to detect current namespace: %s\n", err)
	}
	defer fd.Close()

	data, err := io.ReadAll(fd)
	if err != nil {
		log.Fatalf("failed to read current namespace: %s\n", err)
	}
	scanner.ns = string(data)

	v1alpha1.AddToScheme(scheme.Scheme)

	return scanner
}

func (scanner *Scanner) Ingresses() []*v1alpha1.Ingress {
	ingresses := scanner.ingressStore.List()
	result := make([]*v1alpha1.Ingress, 0)
	for _, ingressG := range ingresses {
		ingress := ingressG.(*v1.Ingress)

		for _, rule := range ingress.Spec.Rules {
			for _, path := range rule.IngressRuleValue.HTTP.Paths {
				newIngress := &v1alpha1.Ingress{
					Namespace:   ingress.Namespace,
					Host:        rule.Host,
					Path:        path.Path,
					PathType:    string(*path.PathType),
					ServiceName: path.Backend.Service.Name,
					ServicePort: int(path.Backend.Service.Port.Number),
				}
				result = append(result, newIngress)
			}
		}
	}
	return result
}

func (scanner *Scanner) FrontdoorIds() []string {
	frontdoors := scanner.frontdoorStore.List()
	result := make([]string, len(frontdoors))
	for idx, frontdoor := range frontdoors {
		fd := frontdoor.(*v1alpha1.Frontdoor)
		result[idx] = fd.Spec.FrontdoorId
	}
	return result
}

func (scanner *Scanner) IpAddresses() []string {
	ipAddresses := scanner.ipaddressStore.List()
	result := make([]string, len(ipAddresses))
	for idx, ipAddress := range ipAddresses {
		ipAddr := ipAddress.(*v1alpha1.IpAddress)
		result[idx] = ipAddr.Spec.IpAddress
	}
	return result
}

func (scanner *Scanner) Start() {
	scanner.ingressStore = scanner.WatchIngresses()
	scanner.frontdoorStore = scanner.WatchFrontdoors()
	scanner.ipaddressStore = scanner.WatchIpAddresses()
}

func (scanner *Scanner) WatchFrontdoors() cache.Store {
	frontdoorStore, frontdoorController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return scanner.fdclient.Frontdoors(scanner.ns).List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return scanner.fdclient.Frontdoors(scanner.ns).Watch(lo)
			},
		},
		&v1alpha1.Frontdoor{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go frontdoorController.Run(wait.NeverStop)
	return frontdoorStore
}

func (scanner *Scanner) WatchIpAddresses() cache.Store {
	ipAddressStore, ipAddressController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				ctx := context.TODO()
				return scanner.fdclient.IpAddresses(scanner.ns).List(ctx, lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				ctx := context.TODO()
				return scanner.fdclient.IpAddresses(scanner.ns).Watch(ctx, lo)
			},
		},
		&v1alpha1.IpAddress{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go ipAddressController.Run(wait.NeverStop)
	return ipAddressStore
}

func (scanner *Scanner) WatchIngresses() cache.Store {
	ingressStore, ingressController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				ctx := context.TODO()
				return scanner.client.NetworkingV1().Ingresses("").List(ctx, lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				ctx := context.TODO()
				return scanner.client.NetworkingV1().Ingresses("").Watch(ctx, lo)
			},
		},
		&v1.Ingress{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go ingressController.Run(wait.NeverStop)
	return ingressStore
}
