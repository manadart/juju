<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/.sphinx/_static/logos/juju-logo-dark.png?raw=true">
  <source media="(prefers-color-scheme: light)" srcset="docs/.sphinx/_static/logos/juju-logo.png?raw=true">
  <img alt="Juju logo next to the text Canonical Juju" src="docs/.sphinx/_static/logos/juju-logo.png?raw=true" width="30%">
</picture>

Juju is an open source application orchestration engine that enables any
application operation (deployment, integration, lifecycle management) on any
infrastructure (Kubernetes or otherwise) at any scale (development or
production) in the same easy way (typically, one line of code), through special
operators called ‘charms’.

[![juju](https://snapcraft.io/juju/badge.svg)](https://snapcraft.io/juju)
[![snap](https://github.com/juju/juju/actions/workflows/snap.yml/badge.svg)](https://github.com/juju/juju/actions/workflows/snap.yml)
[![build](https://github.com/juju/juju/actions/workflows/build.yml/badge.svg)](https://github.com/juju/juju/actions/workflows/build.yml)

- [Give it a try!](https://documentation.ubuntu.com/juju/latest/tutorial/)
- Read the [docs](https://documentation.ubuntu.com/juju/).
- Read our [Code of conduct](https://ubuntu.com/community/code-of-conduct) and join our [chat](https://matrix.to/#/#charmhub-juju:ubuntu.com) and [forum](https://discourse.charmhub.io/) or [open an issue](https://github.com/juju/juju/issues).
- Read our [CONTRIBUTING guide](./CONTRIBUTING.md) and contribute!

## Kubernetes controller HA spike

This branch contains exploratory work toward high availability for Juju
controllers running on Kubernetes. The spike is intended to illustrate the
shape of a possible implementation rather than present a final design.

The prototype makes the controller workload scale past one pod by deriving a
stable controller identity from each StatefulSet ordinal. It also introduces a
separate headless service for Dqlite so each controller pod has a unique,
routable peer address instead of sharing the normal controller API service
address.

Controller startup then generates per-pod agent configuration from the bootstrap
controller template. That configuration carries the pod-specific controller
identity and the Dqlite peer addresses needed to form a cluster. Dqlite startup
can use those configured bind addresses on Kubernetes rather than being limited
to loopback.

Bootstrap also seeds the initial Dqlite cluster configuration and the
credentials required by the Kubernetes controller path. When additional
controller pods register as Dqlite nodes, the worker records a matching
controller-node password so those agents can authenticate to their local API.

The overall model is: give every controller pod stable identity, provide Dqlite
with unique peer addresses, propagate those addresses into controller
configuration, and ensure every scaled controller agent has matching credentials
in the controller database.
