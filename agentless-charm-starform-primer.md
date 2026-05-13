# Agentless Charm / Starform Primer

This is a working-session primer for a future clean session. It is not a
user-facing design document and does not imply an implementation has started.

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

- `/home/joseph/projects/canonical/juju`

Starform clone:

- `/home/joseph/projects/canonical/starform`

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

- `/home/joseph/projects/canonical/starform/README.md`
- `/home/joseph/projects/canonical/starform/starform/scriptset.go`
- `/home/joseph/projects/canonical/starform/starform/appobject.go`
- `/home/joseph/projects/canonical/starform/starform/eventobject.go`
- `/home/joseph/projects/canonical/starform/formtest/formtest_test.go`
- `/home/joseph/projects/canonical/starform/SECURITY.md`

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
- The runner and `jujuc` machinery are not a good fit for agentless Starform
  charms because they assume an executing charm process and hook-tool context.
- `domain/unitstate.CommitHookChangesArg` is a useful precedent for
  post-hook transactional output.
- `CommitHookChanges` validates and applies hook changes only after successful
  hook completion.

Important mismatch:

- `CommitHookChangesArg` is unit-scoped.
- Agentless charm logic may need application-scoped state and relation writes,
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

Possible intent API examples, not decisions:

- `juju.status_set(status, message="")`
- `juju.state_set(key, value)`
- `juju.state_delete(key)`
- `juju.relation_set(endpoint, data, scope="application")`
- `juju.open_port(endpoint, protocol, port)`
- `juju.close_port(endpoint, protocol, port)`

The names above are illustrative. The important property is that methods append
typed intents to a collector. Juju validates and applies them after execution.

## Hackathon Plan

For the hackathon, bias towards a narrow end-to-end demo over a durable
production design. Automated tests are not required for the spike; prefer a
small manual demo that proves deployment, relation, scriptlet execution, intent
application, and status visibility.

Assume these simplifying constraints unless they become blockers:

- Treat an agentless charm as an application with no real units.
- Prefer application-scoped status, state, and relation data.
- Avoid secrets, leadership, storage, actions, ports, and unit-scoped relation
  data.
- Use one controller-side worker path and accept crude HA behaviour for the
  demo if needed.
- Use a tiny event vocabulary, probably config changed and relation changed.
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
   - Add the smallest metadata marker needed to identify a charm as
     scriptlet-based.
   - Decide where Starform files live in the charm archive.
   - Package one demo charm that responds to config and relation input.

3. Add a normal charm for the other side of the relation.
   - Use a simple relation interface that can exchange application data with
     the scriptlet charm.
   - Keep the charm conventional so any demo weirdness is isolated to the
     scriptlet side.

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
- When related to a normal charm, write application relation data such as:
  - `ready=true`
  - `model-uuid=<uuid>`
  - `model-name=<name>`
  - `model-type=iaas|caas`
  - `cloud=<cloud>`
  - `region=<region>`
  - `generation=<monotonic counter or content hash>`
- Treat `ready=true` as the obvious "scriptlet executed successfully" signal.
- Treat `generation` as the watched freshness value and bump it whenever the
  scriptlet republishes relation data.
- The normal charm should handle `relation-changed`, read the remote
  application bag, and set status to show the received generation.
- A scriptlet charm config key such as `refresh-token` can be used as an easy
  manual trigger: changing it should republish model information with a new
  generation.

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

## Open Questions

- What is the entity identity for an agentless charm: application, synthetic
  unit, controller unit, or a new concept?
- How many workers run for one agentless application in HA controller setups,
  and what lease prevents duplicate execution?
- Which MVP events should exist first: config, update-status, relation events,
  lifecycle events, or a single generic hook event?
- How should event ordering and retries be persisted across controller restarts?
- Where are Starform script sources stored in a charm package, and how is a
  charm marked as agentless/scriptlet-based?
- Should relation data be application-only, or should there be a unit-scoped
  concept for compatibility?
- How are status and charm state represented without a unit agent?
- Which secret operations, if any, are safe for an MVP?
- How are script errors surfaced: application status, model status, controller
  log, or a new error surface?
- How are charm upgrades handled, including reloading scripts and preserving
  state?

## Suggested Next Session

Start by re-reading the Starform files listed above, then the Juju uniter and
unitstate files listed above. Do not begin by wiring dependencies into Juju.

For a production-oriented minimal spike, define a small host app object and
collector in isolation:

- One event, probably `config_changed` or generic `hook`.
- One or two event attributes.
- One or two intent methods, probably status and charm state.
- A test showing that failed script execution discards intents.
- A test showing successful execution validates and converts intents into a
  domain-layer argument shape.

For the hackathon path, skip these tests and validate the same behaviour through
the manual demo flow described above.

Only after that should the worker shape be designed.
