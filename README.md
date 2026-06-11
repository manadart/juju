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

# Kubernetes controller HA spike

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

# Review: Kubernetes Controller HA Spike (`manadart/k8s-ha-spike`)

**Reviewer:** Sinan (with multi-agent review fleet: 6 review lenses, adversarial
verification of every finding, 3 independent alternative designs, judge synthesis —
108 agents; 42 findings confirmed, 2 refuted)
**Branch:** 9 commits by Joseph Phillips on `main` @ `484adafe3e` (4.1-beta2), +965/−90 across 27 files
**Date:** 2026-06-10

---

## Verdict

**Correct direction, with two sub-mechanisms that must be replaced — not polished —
before productization.** As a spike, it succeeds: it proves the controller workload
can scale past one pod with stable identity and a working dqlite cluster, and it
usefully *disproves* (by counterexample) that cluster membership can be derived
statically from a pod's own ordinal. Recommend proceeding on this direction with
the modifications below.

## What the spike does

Pod ordinal = controller identity. A new headless Service (`<stack>-dqlite`,
`PublishNotReadyAddresses`) gives each controller pod a stable DNS peer address;
a new `dqlite-bind-address` agent-config key lifts the CAAS loopback-only dqlite
restriction. A shell script embedded in the pod command seds the bootstrap agent.conf
template per ordinal and writes `controller.conf` (`db-bind-addresses`), feeding the
**existing, unmodified** IAAS dbaccessor join/rebind machinery. The CAAS provisioner's
status-only path gains `EnsureControllerScale` to drive StatefulSet replicas
(scale-up only). Bootstrap seeds an application password and secret env so scaled
pods can complete unit introduction; dbaccessor writes a controller-node password
row at dqlite init so scaled agents can log in to their local API.

## Build & run (all verified live)

- `go build ./...` clean; all touched packages pass `-race` unit tests.
- Bootstrapped microk8s, `juju scale-application -m controller controller 3`:
  **3-voter dqlite cluster formed** on stable headless DNS, identical
  `cluster.yaml` on all nodes.
- **Killed the bootstrap pod (`controller-0`): the Juju API never blipped**; pod
  recreated in ~16 s, rejoined as voter; DB writes succeeded post-failover.
- Caveats found while running: scaled controller agents (ordinal > 0) crash-loop
  their `api-caller`-dependent workers — `apiserver/facades/agent/agent/agent.go:299`
  hardcodes `NotSupported("in HA")` for controller-agent tags > 0 (an upstream
  `TODO(ha)`, not spike code, but the spike is incomplete without it).
  `modeloperator` crash (`jujud model` unregistered) is pre-existing main breakage,
  unrelated to the spike.

## What the spike proves

1. Headless service + `ServiceName` switch is the canonical K8s pattern for quorum
   systems, and it works here (previously per-pod DNS didn't exist at all).
2. dqlite binds/advertises a DNS name on CAAS; identity survives pod reschedule
   (go-dqlite pins the address in `info.yaml` on the PVC); **TLS works unchanged**.
3. The IAAS dbaccessor cluster state machine is substrate-portable — scaled K8s
   nodes join with zero changes to the join logic. The seam choice is right.
4. The status-only provisioner path is the right hook for controller scaling, and
   unit introduction works for scaled controller pods with seeded credentials.

## Direction-level problems (must be redesigned, not patched)

1. **Static, truncated `controller.conf` arms a split-brain.** Each pod writes
   members `0..own-ordinal` once per container start; pod 0 permanently sees a
   single-member cluster. dbaccessor treats this file as authoritative full
   membership: an extant pod 0 restarting without quorum (rolling restart,
   partition) hits the `serverCount==1` backstop → `SetClusterToLocalNode` →
   forced single-node raft reconfiguration while pods 1..N still hold the 3-node
   config. Two divergent controller DBs behind one Service. Also: pod-0 PVC loss
   re-runs `bootstrap-state` into a rival cluster. Membership must be dynamic,
   symmetric, and notified (DB/charm/watch-driven) — this is a prerequisite of the
   design, not hardening.
2. **Live-observed config ownership conflict:** the controller charm's IAAS
   `db-bind-addresses` writer overwrote pod-1's `controller.conf` with the API
   service VIP for every entry (useless as dqlite peer addresses) while pods 0/2
   kept the init-script content — three pods, three different cluster views in a
   10-minute-old deployment. Two uncoordinated writers of the same contract file;
   one owner must be chosen.
3. **One shared fleet credential.** The sed template copies controller-0's API
   password, old password, and TLS server key to every pod; the app password
   reuses the same unit password held in a Secret/ConfigMap. Compromise of any
   pod = compromise of the fleet. Per-node credential provisioning has to be part
   of the identity mechanism itself.

## Significant but fixable (selected from 42 confirmed findings)

- **Scale-down wedges the controller** (verified live): CLI accepts
  `scale-application controller 2`, provisioner persists scaling state, then
  crash-loops on `NotSupported`; status stuck at `3/2`. Validate at the
  domain/facade layer *before* persisting; worker should treat it as a quiet no-op.
- **No topology validation**: scale 2 silently yields a non-replicating spare
  (false HA); even counts accepted; any model writer can scale the controller.
- **No adoption path for existing controllers**: `StatefulSet.spec.serviceName`
  is immutable — pre-spike controllers can't be migrated without STS recreate, and
  ungated scaling of an old bootstrap-frozen template would mint duplicate
  controller-0s. Needs a capability gate + migration plan.
- **No join-failure retry**: dbaccessor swallows `errNotReady` after a 60 s
  Ready timeout; nothing on K8s re-triggers reconciliation (no `/reload` notify
  channel exists on CAAS) — convergence currently relies on ordering luck.
- **Shell-script config generation** (HOSTNAME parsing + sed over YAML) couples
  three layers' serialization formats with zero error handling; must become Go
  (e.g. `jujud controller-init` or the existing `containeragent init --controller`).
- dbaccessor password write runs unconditionally (IAAS too) through a
  `NewService(nil, …)`-armed domain service; gate to CAAS and construct properly.
- Missing for real HA (scoped, deferred legitimately): member removal/scale-down,
  per-node API address publication (apiremotecaller/object-store replication gets
  no peer addresses), pod-0 bootstrap-role decoupling, rolling-upgrade quorum
  handover, dqlite-aware readiness, enable-ha vs scale-application UX decision.

## Engineering-practice notes (judged as a spike: good)

Convention-literate work — error-handling idioms match each file, the new
`EnsureControllerScale` mirrors the existing `EnsureScale`/`tryAgain` idiom, the
removed status-only scaling-state reset is correctly reasoned (its 2023 rationale
no longer holds), and the commit series is clean and reviewable. Tests cover the
happy paths; the shell script is only golden-string asserted. The `agent.SetValue`
nil-map fix is a genuine pre-existing bug — extract and land separately (with a
test that actually exercises the nil branch, which the current one doesn't).

## Recommended path

- **Phase 0 (land now, HA-independent):** headless dqlite Service + ServiceName
  switch (label-marked auxiliary services instead of name matching);
  `dqlite-bind-address` key + NodeManager support; `SetValue` nil-map fix;
  `controllerNamespace()` fallback fix.
- **Phase 1 (productize scale-up-only HA, ~3–5 eng-weeks, single repo):** replace
  the shell script with a `jujud controller-init` Go subcommand (downward-API pod
  name, typed agent.Config writer, **fresh per-node password**, bootstrap gated on
  ordinal-0 *and* cluster-existence); replace static `controller.conf` with a
  CAAS-only worker that watches `controller_node` + desired scale and writes the
  **full symmetric membership** on every node with a programmatic notify into the
  existing dbaccessor reconciler; CAAS-gate the password write; reject scale-down
  and invalid topologies at the service layer; pin `OrderedReady`/update strategy
  and capability-gate scaling; publish per-pod API addresses. Resolve the upstream
  `agent.go:299` ordinal-0-only TODO — without it scaled controller agents run
  degraded.
- **Phase 2 (membership lifecycle & UX):** scale-down via preStop handover +
  member removal + PVC cleanup; pod-0 role decoupling; credential rotation;
  upgrade/adoption migration; decide `enable-ha` vs blessing `scale-application`.
  Re-evaluate the **charm-driven** membership writer here (JU179 alignment) — the
  Phase 1 symmetric-conf machinery is forward-compatible with it either way.

## Alternatives considered (ranked by judge panel)

1. **Ordinal HA, Go-native** (medium effort) — keep spike topology; `controller-init`
   subcommand + symmetric DB-watching cluster-config writer. *Recommended Phase 1.*
2. **Charm-driven membership** (large) — IAAS parity via dbcluster peer relation;
   right long-term owner candidate but inverts failure domains (charm needs a live
   API to repair a dead cluster) and couples two repos. *Re-evaluate in Phase 2.*
3. **K8s-native watch/operator** (large) — EndpointSlice-driven membership,
   preStop handover; most capable, but injects a K8s API dependency into the DB
   clustering path. Cherry-pick its lifecycle ideas for scale-down.
4. Spike-as-is — not landable (split-brain, shared credential, two config writers).

## Reproduction / artifacts

- Worktree `/data/dev/juju-k8s-ha-spike` (branch `k8s-ha-spike`, remote `manadart`).
- Live controller `spike-ha` on microk8s left running for inspection
  (`juju destroy-controller spike-ha --destroy-all-models` to clean up; local
  registry holds `localhost:32000/juju/{jujud-operator:4.1-beta2,charm-base}`).
- Bootstrap gotcha on main: `caas-image-repo=localhost:32000/juju` fails
  (Go `url.Parse` scheme quirk) — use the JSON form with explicit
  `"serveraddress":"http://localhost:32000"`.

