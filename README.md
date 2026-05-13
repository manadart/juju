![SCRIPTLETS](./scriptlets.gif)

# Unitless Charms: The Controller Takes The Reins

This branch is about a new charm execution mode: no unit agent, no workload
machine, no little hook process standing in the dust waiting for orders.

The controller watches the model. The controller loads the charm's Starform
scriptlets. The controller runs them under tight limits. The controller collects
their declared intent. Then Juju decides what becomes law.

That is the whole proposition.

## The Shape Of The Thing

Traditional charms execute hook code through a unit agent. This spike cuts that
path down to its bones for a new class of charm:

- The charm is deployed as an application.
- It has no real units.
- Its logic is Starform scriptlets.
- Scriptlets observe model events.
- Scriptlets do not mutate Juju directly.
- Scriptlets emit typed intent.
- Juju validates and applies that intent after script execution succeeds.

No in-process `jujuc` clone. No live domain services exposed to Starlark. No
scriptlet side effects sneaking out of the saloon.

The scriptlet speaks. Juju judges.

## Starform

Starform gives us event-based Starlark scriptlets. A script registers handlers
during `init`:

```python
def init():
    juju.observe("config_changed", on_config_changed)
```

When the controller sees a relevant model event, it builds an event snapshot,
runs the registered Starform handlers, and collects the intents they append.
Handler return values do not matter. The collector matters.

For this spike, the host object is named `juju`.

Initial intent calls:

```python
juju.status_set("active", message = "config changed")
juju.state_set("last_event", event.name)
juju.state_set("last_message", "config changed")
```

Relation publishing will follow the same rule: collect first, apply later.

## Scriptlet Charm

The first scriptlet charm lives here:

```text
scriptlet/
```

Current files:

- `scriptlet/metadata.yaml`
- `scriptlet/config.yaml`
- `scriptlet/scriptlet.yaml`
- `scriptlet/scriptlets/hooks.star`

The hackathon marker is `scriptlet/scriptlet.yaml`:

```yaml
unitless: true
runtime: starform
app: juju
sources:
  - scriptlets/hooks.star
```

Current event names are Starform-safe snake case:

- `config_changed`
- `relation_created`
- `relation_joined`
- `relation_changed`
- `relation_departed`
- `relation_broken`

The handlers are stubs for now. They prove the shape: observe events, emit
trivial intent, keep moving.

## Relation Interface

The demo relation interface is:

```text
scriptlet-model-info
```

The scriptlet charm provides:

```yaml
provides:
  model-info:
    interface: scriptlet-model-info
```

The first meaningful relation data will be application-scoped and will include:

- `controller-uuid`
- `model-name`

More fields are waiting in the saddle:

- `ready`
- `model-type`
- `cloud`
- `region`
- `generation`

The important part is not the list. The important part is ownership: the
scriptlet charm declares what relation data should exist, and Juju applies it
from the controller side.

## Conventional Consumer Charm

The other side of the demo is a normal bash charm:

```text
testcharms/charms/model-info-consumer/
```

It requires the same `scriptlet-model-info` interface. On
`model-info-relation-changed`, it reads the remote application bag:

- `controller-uuid`
- `model-name`

It stores those values under:

```text
/var/lib/model-info-consumer/
```

It also includes an action:

```text
report-model-info
```

The action writes this to stdout:

```text
I am related to model <model-name> in controller <controller-uuid>.
```

That sentence is the smoke test. If it prints with real values, the relation
data rode all the way across.

## Controller Worker

The worker still needs to be built. It belongs under:

```text
internal/worker
```

Its job:

1. Find scriptlet-based applications.
2. Load Starform sources from the charm archive.
3. Watch or poll the small event set needed for the demo.
4. Build read-only event attrs.
5. Create an intent collector as `EventObject.State`.
6. Run Starform with memory and CPU limits.
7. Discard intents on script failure.
8. Validate and apply intents on success.
9. Surface errors through logs and application status.

For the hackathon, the worker can be blunt. It needs to prove the trail from
deploy to relation data to conventional charm action. Durable queues, perfect
HA fencing, charm upgrade semantics, secrets, storage, actions, ports, and
unit-scoped relation compatibility can wait outside.

## Dependency

Starform fetch has been validated:

```text
github.com/canonical/starform v0.0.0-20260428155809-8da636f0fff9
```

Local private module setup is expected:

```text
GOPRIVATE=github.com/canonical/starform,github.com/canonical/starlark
```

Git is configured to fetch GitHub modules through SSH.

## Demo Target

The intended demo is narrow and ruthless:

1. Deploy the unitless charm with `juju deploy-unitless ./scriptlet`.
2. Deploy `model-info-consumer`.
3. Relate them over `scriptlet-model-info`.
4. Trigger scriptlet execution with config or relation events.
5. Have the controller-side worker publish application relation data.
6. Watch the conventional charm receive it.
7. Run `report-model-info`.
8. See the model name and controller UUID printed from the normal charm.

When that works, the controller has crossed the street, stared down the old
unit-agent assumption, and taken the hook loop for itself.

## Not Yet Settled

The hackathon path is application-scoped and deliberately severe. The
production questions remain:

- Whether the long-term identity is application, synthetic unit, controller
  unit, or a new entity.
- How HA controllers prevent duplicate execution.
- How event ordering and retries survive restarts.
- How application-scoped charm state becomes first-class domain state.
- How charm upgrades reload scriptlets and preserve state.
- Whether `scriptlet.yaml` survives as the real package marker.

Those fights come later.

For now, the controller rides.
