kubernetes-cluster-manager (KCM)
================================

[![Build Status](https://travis-ci.com/martinohmann/kubernetes-cluster-manager.svg?branch=master)](https://travis-ci.com/martinohmann/kubernetes-cluster-manager)
[![codecov](https://codecov.io/gh/martinohmann/kubernetes-cluster-manager/branch/master/graph/badge.svg)](https://codecov.io/gh/martinohmann/kubernetes-cluster-manager)
[![Go Report Card](https://goreportcard.com/badge/github.com/martinohmann/kubernetes-cluster-manager?style=flat)](https://goreportcard.com/report/github.com/martinohmann/kubernetes-cluster-manager)
[![GoDoc](https://godoc.org/github.com/martinohmann/kubernetes-cluster-manager?status.svg)](https://godoc.org/github.com/martinohmann/kubernetes-cluster-manager)

Inspired by [Zalando's CLM](https://github.com/zalando-incubator/cluster-lifecycle-manager). The Kubernetes Cluster Manager project was started because CLM is tightly coupled to AWS Cloudformation for managing the cluster infrastructure. KCM tries to provide an interface for using different infrastructure manager and manifest renderers. It also tries to provide visibility about changes by providing diffs for things like manifest changes.

**Use with caution. This project is currently alpha quality and APIs are likely to change until the first stable release.**

Features:
- Make output of infrastructure manager available to manifest renderer
- Show diffs of changes in infrastructure output values and manifests
- Render manifests with helm
- Minikube integration for local testing
- Dry run, apply and destroy changes (infrastructure + kubernetes manifests)

Currently supported infrastructure managers:
- [Terraform](https://github.com/hashicorp/terraform)
- [Minikube](https://github.com/kubernetes/minikube) for local testing

Currently supported manifest renderers:
- [Helm](https://github.com/helm/helm)

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

The documentation is still work in progress. For now, refer to [godoc](https://godoc.org/github.com/martinohmann/kubernetes-cluster-manager) and the command line help:

```sh
$ kcm help
```

Provision infrastructure using terraform and render manifests via helm:

```sh
$ kcm provision \
  --manager terraform \
  --renderer helm \
  --working-dir /path/to/terraform/repo \
  --helm-chart /path/to/cluster/helm/chart \
  --dry-run
```

As the bare minimum `kcm` expects the infrastructure manager to create a kubernetes cluster and to either return the path to a generated `kubeconfig` in its output, or `server` and `token` values needed for establishing a connection to the kubernetes api-server. Detailed examples will follow.

License
-------

The source code of kubernetes-cluster-manager is released under the MIT License. See the bundled
LICENSE file for details.
