[![Go Report Card](https://goreport.com/badge/github.com/postfinance/kubewire)](https://goreport.com/badge/github.com/postfinance/kubewire)
[![Build Status](https://travis-ci.org/postfinance/kubewire.svg?branch=master)](https://travis-ci.org/postfinance/kubewire)

# kubewire
kubewire is a Kubernetes integrity checker which acts as a tripwire for global
Kubernetes resources or namespaced resources which could impact the
whole cluster.

*Status*: Alpha, anything can change at any time

## Use case
Kubernetes cluster administrators have great power. This means that
a mistake they make could cause the cluster to become unhealthy or insecure and,
as such, could impact any or all tenants sharing the cluster. Kubewire does not 
prevent mistakes but it is intended to notice modifications.

Common sources for such modifications are:
* `kubectl create` on objects which define a wrong namespace
* Wrong kubeconfig or a not defined namespace
* Running tools which create object in different namespaces e.q. Helms Tiller is deployed to kube-system by default

Kubewire is not focused on hidden malicious acts and also does not keep any
object backed up. So it's best used together with an automated deployment/configuration
tool which ensures that all global objects have the state you wish. Kubewire
just ensures that no additional objects are created unintentionally.

## Installation
In order to get the latest version, do
```
go get -u github.com/postfinance/kubewire
```

You should hava a working Go installation.

Precompiled binaries will be provided soon for the common platforms.

## Usage
By default, all non-namespaced resources will be scanned. In addition to that,
the following namespaces are considered to hava a global effect, so the namespaced
resources of them will also be scanned:
- default
- kube-system
- kube-public

This list can be customized with the `--namespaces` flag.

```
$ kubewire snapshot > baseline.yaml

$ ./thisdoessomemagic

$ kubewire diff --baseline=baseline.yaml
Element												                    A				                	B
ScanStart										                	    2018-06-12 14:19:14.152560709 +0200 CEST2018-06-14 10:22:18.083728367 +0200 CEST m=+0.028297121
ScanEnd												                    2018-06-12 14:19:42.870490496 +0200 CEST2018-06-14 10:22:46.602422832 +0200 CEST m=+28.546991607
ResourceObjects." v1 namespaces  appl-shouldnotbehere"					does not exist			        	exists
ResourceObjects." v1 secrets kube-system shouldnotbehere-token-rwmcl"	does not exist			        	exists
ResourceObjects." v1 serviceaccounts kube-system shouldnotbehere"		does not exist			        	exists
```

### Kubeconfig
kubewire detects if it is running in a Kubernetes cluster and uses the service account
of the Pod if available. If this is not the case, it looks in the default kubectl
paths for a kubeconfig. Both cases can be overriden by setting the 'kubeconfig' flag.

### RBAC Rules
RBAC Rules will be provided soon

## Requirements
This utility should work with any recent Kubernetes compatible cluster.

It is known to work with:
- Kubernetes 1.9+

It is known not to work:
- OpenShift <=3.5

## Next

- [ ] Scan namespaced resources with global impact e.g. PodSecurityPolicy usages
- [ ] Add example reports
- [ ] Provide RBAC yaml
- [ ] Review ReportDiff format and make it more usable and readable
- [ ] Add more tests
