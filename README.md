[![Go Report Card](https://goreportcard.com/badge/github.com/postfinance/kubewire)](https://goreportcard.com/report/github.com/postfinance/kubewire)
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
In order to compile the latest version from source, do a
```
go get -u github.com/postfinance/kubewire
```

Precompiled binaries are available on [Github Releases](https://github.com/postfinance/kubewire/releases)

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
Element                                                                 A                                        B
ScanStart                                                               2018-06-12 14:19:14.152560709 +0200 CEST 2018-06-14 10:22:18.083728367 +0200 CEST m=+0.028297121
ScanEnd                                                                 2018-06-12 14:19:42.870490496 +0200 CEST 2018-06-14 10:22:46.602422832 +0200 CEST m=+28.546991607
ResourceObjects." v1 namespaces  appl-shouldnotbehere"                  does not exist                           exists
ResourceObjects." v1 secrets kube-system shouldnotbehere-token-rwmcl"   does not exist                           exists
ResourceObjects." v1 serviceaccounts kube-system shouldnotbehere"       does not exist                           exists
```

#### Other functions
Kubewire supports the following commands:

```
$ kubewire -h
...
  diff            Compare snapshots with another or a live cluster
  help            Help about any command
  resourceobjects List API resource objects
  resources       List API resources
  serverinfo      Prints server info
  snapshot        Take a snapshot of cluster resources and objects
```

so you can use it to list or export an inventory of API resources and their objects.
The supported export formats are json and yaml.

Example listings:
```
$ kubewire resources
GroupVersion              Kind              Name               Namespaced  Verbs
v1                        Binding           bindings           true        [create]
v1                        ComponentStatus   componentstatuses  false       [get list]
v1                        ConfigMap         configmaps         true        [create delete deletecollection get list patch update watch]
apps/v1beta1              Deployment        deployments        true        [create delete deletecollection get list patch update watch]
crd.projectcalico.org/v1  BGPConfiguration  bgpconfigurations  false       [delete deletecollection get list patch create update watch]
...

$ kubewire resourceobjects
GroupVersion  Resource           Namespace    Name
v1            componentstatuses               controller-manager
v1            componentstatuses               etcd-0
v1            componentstatuses               etcd-1
v1            componentstatuses               etcd-2
v1            componentstatuses               scheduler
v1            configmaps         kube-system  calico-config
v1            secrets            default      default-token-wsq94
v1            secrets            kube-public  default-token-c5qs4
apps/v1       daemonsets         kube-system  calico-node
apps/v1       daemonsets         kube-system  ip-masq-agent
...
```

### Kubeconfig
kubewire detects if it is running in a Kubernetes cluster and uses the service account
of the Pod if available. If this is not the case, it looks in the default kubectl
paths for a kubeconfig. Both cases can be overriden by setting the 'kubeconfig' flag.

### RBAC Rules
kubewire needs permission to list all resource objects in a Kubernetes cluster.
It does not require to get the objects itself.

An example ClusterRole and ClusterRoleBinding is provided in the [rbac.yaml](deployment/rbac.yaml) file,
which assumes that kubewire runs in a Pod as the service account `kubewire` in the
`kube-system` namespace.

## Requirements
This utility should work with any Kubernetes 1.7+ compatible cluster.

## Next

- [ ] Scan namespaced resources with global impact e.g. PodSecurityPolicy usages
- [ ] Add example reports
- [ ] Review ReportDiff format and make it more usable and readable
- [ ] Add more tests
