# FD Ingress Controller

A kubernetes ingress controller designed to filter traffic from Azure frontdoor.

## Deploy

The helm chart deploys the controller and 2 CRDs. 

`helm upgrade fdingress ./helm/frontdoor-ingress/ --namespace ingress --install`

## Allow a frontdoor ID

The following is an example of the resource to allow a given frontdoor ID. It must be deployed to the same namespace as the controller

```yaml
---
apiVersion: fdingress.com/v1alpha1
kind: Frontdoor
metadata:
  name: test-frontdoor
  namespace: ingress
spec:
  frontdoorId: 11111111-2222-3333-4444-5555-666666666666
```

## Allow an IP address

The following allows an given IP address access to the route. It must be deployed to the same namespace as the controller.

```yaml
---
apiVersion: fdingress.com/v1alpha1
kind: IpAddress
metadata:
  name: test-ipaddress
  namespace: ingress
spec:
  ipAddress: 10.1.2.3
```