kubernetes-cluster-manager (`kcm`)
==================================

[![Build Status](https://travis-ci.com/martinohmann/kubernetes-cluster-manager.svg?branch=master)](https://travis-ci.com/martinohmann/kubernetes-cluster-manager)
[![codecov](https://codecov.io/gh/martinohmann/kubernetes-cluster-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/martinohmann/kubernetes-cluster-manager)
[![Go Report Card](https://goreportcard.com/badge/github.com/martinohmann/kubernetes-cluster-manager?style=flat)](https://goreportcard.com/report/github.com/martinohmann/kubernetes-cluster-manager)
[![GoDoc](https://godoc.org/github.com/martinohmann/kubernetes-cluster-manager?status.svg)](https://godoc.org/github.com/martinohmann/kubernetes-cluster-manager)

Inspired by [Zalando's
CLM](https://github.com/zalando-incubator/cluster-lifecycle-manager). The
Kubernetes Cluster Manager project was started because CLM is tightly coupled
to AWS Cloudformation for managing the cluster infrastructure. KCM tries to
provide an interface for using different infrastructure provisioners and
manifest renderers. It also tries to provide visibility about changes by
providing diffs for things like manifest changes.

**Use with caution. This project is currently alpha quality and APIs are likely
to change until the first stable release.**

Features:
- Make output of infrastructure provisioners available to manifest renderer
- Show diffs of changes in infrastructure output values and manifests
- Render manifests via helm or plain go templates
- Minikube integration for local testing
- Dry run, apply and destroy changes (infrastructure + kubernetes manifests)

Currently supported infrastructure provisioners:
- `null` (default)
- [`terraform`](https://github.com/hashicorp/terraform)
- [`minikube`](https://github.com/kubernetes/minikube) for local testing

Currently supported manifest renderers:
- [`gotemplate`](https://golang.org/pkg/text/template/) with
  [sprig](https://github.com/Masterminds/sprig) function library (default)
- [`helm`](https://github.com/helm/helm)
- `null`

Design
------

The current design assumes that cluster configuration is managed in a git
repository. `kcm` will write changes to yaml files which then should be
commited by the CI tool of your choice. `kcm` can manage multiple clusters
simply by having separate directories in the repository for the configuration
of each cluster. The git approach is a trade-off to avoid introducing a
persistence layer for `kcm` in the early development phase. However, there
might be other options in the future. Check out the [Roadmap](#roadmap).

Installation
------------

```sh
$ git clone https://github.com/martinohmann/kubernetes-cluster-manager
$ cd kubernetes-cluster-manager
$ make install
```

This will install the `kcm` binary to `$GOPATH/bin/kcm`.

Usage
-----

The documentation is still work in progress. For now, refer to
[godoc](https://godoc.org/github.com/martinohmann/kubernetes-cluster-manager),
the examples [`_examples`](_examples/) directory and the command line help:

```sh
$ kcm help
```

Quick examples
--------------

### Provision infrastructure using terraform and render manifests via helm

```sh
$ kcm provision \
  --provisioner terraform \
  --renderer helm \
  --working-dir /path/to/terraform/repo \
  --templates-dir /path/to/helm/charts \
  --dry-run
```

As the bare minimum `kcm` expects the infrastructure provisioner to create a
kubernetes cluster and to either return the path to a generated `kubeconfig` in
its output, or `server` and `token` values needed for establishing a connection
to the kubernetes api-server. Alternatively you can manually provide kubernetes
credentials via the `--cluster-*` flags. Detailed examples will follow.

### Using a config file and skipping manifest rendering/deployment

```sh
$ kcm provision --config config.yaml --skip-manifests
```

### Working with manifests

The `kcm manifests` command will only render and apply/delete manifests and
will skip any infrastructure changes:

```sh
$ kcm manifests apply \
  --renderer gotemplate \
  --templates-dir templates/ \
  --manifests-dir manifest/ \
  --cluster-kubeconfig ~/.kube/config \
  --cluster-context eks-dev
```

Apply all manifests, even if unchanged:

```sh
$ kcm manifests apply --config config.yaml --all-manifests
```

Delete manifests:

```sh
$ kcm manifests delete --config config.yaml
```

### Destroying a cluster

```sh
$ kcm destroy --dry-run
```

Roadmap
-------

The following features are currently planned, but I'm also happy about other
contributions. PRs welcome!

* Add node pool manager (e.g. for managing [spotinst
  elastigroups](https://api.spotinst.com/introducing-elastigroup/))
* Triggering of rolling updates of node pools
* Support for more provisioners (e.g. Cloudformation, kubeadm)
* Replace shell-execs with native go-libraries where possible (and sensible)
* Add support for other persistence layers and configuration sources besides
  the git approach mentioned in the design section

License
-------

The source code of kubernetes-cluster-manager is released under the MIT
License. See the bundled LICENSE file for details.
