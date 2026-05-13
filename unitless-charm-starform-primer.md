# Unitless Charm / Starform Primer

This is a working-session primer for a future clean session. It is not a
user-facing design document. Some hackathon demo artifacts now exist, but the
controller-side implementation is still a spike.

## Goal

Explore a Juju charm execution mode where charm logic is not run by a unit
agent on a machine or container. Instead, a controller-side worker watches model
events and evaluates charm-provided Starform scriptlets.

The key design constraint from the prior discussion: do not model this as a
`jujuc` replacement. Starform is better suited to collecting declared intent
from scriptlets, then letting Juju validate and apply that intent after script
execution succeeds.

## Repository Context

Juju repo:

- `/home/joseph/go/src/github.com/juju/juju`

Starform source:

- The old primer path `/home/joseph/projects/canonical/starform` was not
  present in this workspace.
- `go list -m github.com/canonical/starform@latest` succeeds with the local
  private-module setup and resolved
  `github.com/canonical/starform v0.0.0-20260428155809-8da636f0fff9`.
- The module has been downloaded into the Go module cache for API reference.

Before making Juju changes, read and follow:

- `AGENTS.md`
- `AGENTS.architecture-rules.md`
- `AGENTS.core-domain-rules.md`

Important Juju layering constraints:

- Put workers under `internal/worker`.
- Keep business workflows in `domain` services.
- Keep persistence behind domain state interfaces.
- Do not put business logic in API facades.
- Workers must be restartable, deterministic, cancellation-aware, and use
  existing worker lifecycle patterns such as `catacomb` or `tomb`.

## Starform Findings

Files to read first:

- `$GOMODCACHE/github.com/canonical/starform@v0.0.0-20260428155809-8da636f0fff9/README.md`
- `$GOMODCACHE/github.com/canonical/starform@v0.0.0-20260428155809-8da636f0fff9/starform/scriptset.go`
- `$GOMODCACHE/github.com/canonical/starform@v0.0.0-20260428155809-8da636f0fff9/starform/appobject.go`
- `$GOMODCACHE/github.com/canonical/starform@v0.0.0-20260428155809-8da636f0fff9/starform/eventobject.go`
- `$GOMODCACHE/github.com/canonical/starform@v0.0.0-20260428155809-8da636f0fff9/formtest/formtest_test.go`
- `$GOMODCACHE/github.com/canonical/starform@v0.0.0-20260428155809-8da636f0fff9/SECURITY.md`

Core model:

- Starform embeds Canonical Starlark and exposes event-based scriptlets.
- Scriptlets register observers during `init` using `app.observe(event, fn)`.
- `ScriptSet.LoadSources` compiles and initializes scripts.
- `ScriptSet.Handle(ctx, event)` invokes registered observers for an event.
- Observer return values are ignored.
- Scriptlet output is host-owned state mutation through host-provided app
  methods, typically appending typed intents into `EventObject.State`.
- The README explicitly frames output as desired target state or action, with no
  host changes applied until after script execution completes.

Important API details:

- `ScriptSetOptions` includes `App`, `Cache`, `Logger`, `RequiredSafety`,
  `MaxAllocs`, `MaxSteps`, and `Modules`.
- `EventObject` has `Name`, `State`, and `Attrs`.
- `Attrs` are exposed as Starlark attributes and should be treated as event
  input.
- `State` is developer-supplied Go state and is the natural place for an intent
  collector.
- Custom app methods are only available while handling an event, not while
  loading modules or during `init`.
- `observe` is only available during `init`.
- Starform event names must be lowercase identifiers, 3-20 chars, with no
  dashes or consecutive underscores.
- Raw Juju hook names like `config-changed` do not fit directly. Either use
  snake-case event names such as `config_changed`, or use a generic event such
  as `hook` with `event.kind = "config-changed"`.

Safety implications:

- Use Canonical Starlark safety flags: `MemSafe`, `CPUSafe`, `TimeSafe`,
  `IOSafe`.
- If `MemSafe` is required, `MaxAllocs` must be set.
- If `CPUSafe` is required, `MaxSteps` must be set.
- Every exposed builtin must be tested for the required safety properties.
- Bound script set size and per-file size.
- Do not expose secrets or sensitive values as plain printable Starlark values.
- Logs are a disclosure surface because scripts can call `print` and `debug`.

Dependency note:

- Starform module is `github.com/canonical/starform`.
- It currently depends on `github.com/canonical/starlark`.
- Juju currently has `go.starlark.net` only as an indirect dependency.
- Adding Starform to Juju will require an explicit dependency decision.
- Local module fetch setup has been validated:
  - `GOPRIVATE`, `GONOPROXY`, and `GONOSUMDB` include
    `github.com/canonical/starform,github.com/canonical/starlark`.
  - Git rewrites `https://github.com/` to `git@github.com:`.

## Juju Findings

Files to read first:

- `internal/worker/uniter/doc.go`
- `internal/worker/uniter/uniter.go`
- `internal/worker/uniter/remotestate/watcher.go`
- `internal/worker/uniter/resolver/loop.go`
- `internal/worker/uniter/operation/runhook.go`
- `internal/worker/uniter/runner/runner.go`
- `domain/unitstate/types.go`
- `domain/unitstate/service/commithook.go`
- `domain/unitstate/state/commithook.go`

Relevant existing model:

- The uniter worker watches unit/application/relation state and dispatches hook
  operations to charm code.
- The uniter structure is useful conceptually: watcher snapshot, resolver,
  operation factory, executor, retry/error state.
- The runner and `jujuc` machinery are not a good fit for unitless Starform
  charms because they assume an executing charm process and hook-tool context.
- `domain/unitstate.CommitHookChangesArg` is a useful precedent for
  post-hook transactional output.
- `CommitHookChanges` validates and applies hook changes only after successful
  hook completion.

Important mismatch:

- `CommitHookChangesArg` is unit-scoped.
- Unitless charm logic may need application-scoped state and relation writes,
  or a deliberately defined synthetic/controller unit identity.
- This identity question affects status, relation data ownership, leadership,
  secrets, hook ordering, and lifecycle.

## Design Direction

Prefer this shape:

1. A controller-side worker watches relevant model events.
2. The worker loads a Starform `ScriptSet` for a charm revision.
3. For each event, the worker builds a read-only event snapshot as
   `EventObject.Attrs`.
4. The worker creates an event-specific intent collector as `EventObject.State`.
5. Starform handlers run with strict safety limits.
6. If execution fails, discard the collector and record/report the error.
7. If execution succeeds, validate collected intents.
8. Apply valid intents through domain services, ideally transactionally.

Avoid this shape:

- Exposing live domain services directly to Starlark.
- Exposing RPC-like `juju.*` calls that mutate state during script execution.
- Recreating `jujuc` as in-process methods.
- Letting scriptlets perform I/O or controller-side side effects directly.

Initial intent API names for the hackathon:

- `juju.status_set(status, message="")`
- `juju.state_set(key, value)`
- `juju.state_delete(key)`
- `juju.relation_set(endpoint, data, scope="application")`

Only `status_set` and `state_set` are used by the first scriptlet stub. Relation
data will be added when the worker-side intent collection and application path
exists. Ports are intentionally out of scope for the hackathon.

The important property is that methods append typed intents to a collector.
Juju validates and applies them after execution.

## Decisions And Artifacts So Far

Initial scriptlet charm:

- Location: `scriptlet/`.
- Charm metadata: `scriptlet/metadata.yaml`.
- Unitless marker: `scriptlet/scriptlet.yaml`.
- Marker shape:
  - `unitless: true`
  - `runtime: starform`
  - `app: juju`
  - `sources: [scriptlets/hooks.star]`
- Starform source location: `scriptlet/scriptlets/hooks.star`.
- Host app object name for Starlark: `juju`.
- Demo config key: `refresh-token`.

Initial event vocabulary:

- `config_changed`
- `relation_created`
- `relation_joined`
- `relation_changed`
- `relation_departed`
- `relation_broken`

These are snake-case Starform event names. They correspond to Juju hook names
but avoid dashes because Starform event names must be valid lowercase
identifiers.

Initial scriptlet handlers:

- All handlers are stubs for now.
- Each stub appends trivial intents:
  - `juju.status_set("active", message = <event message>)`
  - `juju.state_set("last_event", event.name)`
  - `juju.state_set("last_message", <event message>)`
- Real relation/config behavior will be filled in after the worker and intent
  applier exist.

Initial relation interface:

- Interface name: `scriptlet-model-info`.
- Scriptlet charm endpoint: `provides: model-info`.
- Conventional consumer charm endpoint: `requires: model-info`.
- For the demo, relation data is application-scoped.
- The first concrete relation settings expected by the consumer are:
  - `controller-uuid`
  - `model-name`

Conventional test charm:

- Location: `testcharms/charms/model-info-consumer/`.
- It implements the `scriptlet-model-info` relation interface.
- It handles the relation hooks conventionally with bash hooks.
- `model-info-relation-changed` reads remote application settings
  `controller-uuid` and `model-name` and stores them under
  `/var/lib/model-info-consumer`.
- Action `report-model-info` prints to stdout and records an action result:
  `I am related to model <model-name> in controller <controller-uuid>.`

## Hackathon Plan

For the hackathon, bias towards a narrow end-to-end demo over a durable
production design. Automated tests are not required for the spike; prefer a
small manual demo that proves deployment, relation, scriptlet execution, intent
application, and status visibility.

Assume these simplifying constraints unless they become blockers:

- Treat a unitless charm as an application with no real units.
- Prefer application-scoped status, state, and relation data.
- Avoid secrets, leadership, storage, actions, ports, and unit-scoped relation
  data.
- Use one controller-side worker path and accept crude HA behaviour for the
  demo if needed.
- Use the tiny event vocabulary listed in "Decisions And Artifacts So Far".
- Keep script errors visible through logs and/or application status.
- Use manual status and relation checks instead of a full integration suite.

Workstreams:

1. Add a scriptlet execution worker.
   - Put the worker under `internal/worker`.
   - Load Starform sources for scriptlet-based charms.
   - Watch or poll for the small set of events needed by the demo.
   - Build event attrs, run handlers with safety limits, collect intents, and
     discard intents on script failure.
   - For the hackathon, use the simplest scheduling mechanism that survives a
     basic controller restart only if this is easy.

2. Define a scriptlet charm package shape.
   - Initial shape exists in `scriptlet/`.
   - `scriptlet.yaml` is the current hackathon marker.
   - Starform files currently live under `scriptlets/`.
   - The demo charm currently has stub handlers only.

3. Add a normal charm for the other side of the relation.
   - Initial charm exists in `testcharms/charms/model-info-consumer/`.
   - It consumes the `scriptlet-model-info` interface.
   - Keep it conventional so demo weirdness is isolated to the scriptlet side.

4. Add intent collection and application.
   - Start with a tiny script-facing API, such as status, charm state, and
     application relation data.
   - Implement typed intent structs in Go.
   - Apply intents through domain services where practical.
   - If a proper domain path is too large for the hackathon, isolate the
     shortcut behind a narrow adapter so it is easy to replace later.

5. Make scriptlet charms deployable.
   - Allow a scriptlet charm application to exist without normal unit-agent
     provisioning.
   - Prevent the deploy path from waiting for machines, containers, or workload
     agents that will never exist.
   - Ensure removal cleans up enough state for repeated demos.

6. Make `juju status` show something useful.
   - Remove or bypass assumptions that every deployed application must have
     visible workload units.
   - Show the application, relation, and scriptlet execution result.
   - Surface script errors as application status or another obvious status
     field.

7. Wire the manual demo.
   - Deploy the normal charm and scriptlet charm.
   - Relate them.
   - Change config or relation data.
   - Show the scriptlet worker running, applying intents, and updating
     `juju status`.

Concrete demo shape:

- Make the scriptlet charm a model-info provider.
- When related to the conventional test charm, write application relation data
  such as:
  - `ready=true`
  - `controller-uuid=<uuid>`
  - `model-name=<name>`
  - `model-type=iaas|caas`
  - `cloud=<cloud>`
  - `region=<region>`
  - `generation=<monotonic counter or content hash>`
- Treat `ready=true` as the obvious "scriptlet executed successfully" signal.
- Treat `generation` as the watched freshness value and bump it whenever the
  scriptlet republishes relation data.
- The conventional consumer charm handles `relation-changed`, reads the remote
  application bag, stores `controller-uuid` and `model-name`, and sets status
  once both values are present.
- A scriptlet charm config key such as `refresh-token` can be used as an easy
  manual trigger: changing it should republish model information with a new
  generation.
- The first consumer action only requires `controller-uuid` and `model-name`.

Existing uniter behaviour that makes this work:

- The remote-state watcher watches relation membership through
  `unit.WatchRelations`.
- For each relation, it starts `WatchRelationUnits`.
- The relation-unit watcher records both unit data versions and application
  data bag versions in `RelationSnapshot.ApplicationMembers`.
- The relation resolver turns an application data bag version change into
  `relation-changed` with `RemoteUnit: ""` and
  `RemoteApplication: <app-name>`.
- This means a normal charm should react to the scriptlet charm's
  application-scoped relation data without any special hook machinery on the
  normal charm side.

Things to defer:

- Full HA leasing and fencing.
- Durable event queues with precise retry semantics.
- Complete charm upgrade semantics.
- Unit-scoped relation data compatibility.
- Secret support.
- Production-quality API/facade changes.
- Stress, race, and integration tests.

## Remaining Open Questions

- Production identity remains open: application, synthetic unit, controller
  unit, or a new concept. The hackathon decision is "application with no real
  units" unless blocked by deploy/status assumptions.
- How many workers run for one unitless application in HA controller setups,
  and what lease prevents duplicate execution?
- How should event ordering and retries be persisted across controller restarts?
- How are script errors surfaced: application status, model status, controller
  log, or a new error surface?
- How are charm upgrades handled, including reloading scripts and preserving
  state?
- Is `scriptlet.yaml` the right long-term marker, or only a hackathon marker?
- What schema/domain support is needed to make application-scoped charm state
  first-class instead of a demo shortcut?

## Suggested Next Session

Start by re-reading the Starform files listed above, then the Juju uniter and
unitstate files listed above. The private module fetch works, so the next
implementation session can add the Starform dependency when the worker code is
ready for it.

For the next implementation step, define a small host app object and collector:

- Load `scriptlet/scriptlets/hooks.star`.
- Support the `config_changed` event first.
- Expose `juju.status_set` and `juju.state_set`.
- Collect intents in `EventObject.State`.
- On script failure, discard the collector and report the error.
- On script success, log or apply the trivial status/state intents.

For the hackathon path, skip automated tests and validate through the manual
demo flow described above.
