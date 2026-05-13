# Scriptlet Worker — Hackathon Implementation Plan

## Overview

Building a controller-side "scriptlet" worker for Juju that enables agentless
charm execution via Starform scriptlets. Charm logic runs on the controller
rather than in unit agents. The worker watches for scriptlet applications and
dispatches events to per-application child workers.

## What's Been Done

### 1. Worker Package (`internal/worker/scriptlet/`)

| File | Purpose |
|------|---------|
| `doc.go` | Package documentation |
| `worker.go` | Main worker with catacomb + `worker.Runner` (objectstore pattern). `ScriptletService` interface. Per-app `applicationRunner` child workers. |
| `manifold.go` | `ManifoldConfig`, `Config`, `Manifold()`, `GetScriptletService()` helper |
| `event.go` | `EventKind` constants: install, start, stop, config_changed, relation_joined/changed/departed/broken, update_status |
| `intent.go` | `IntentCollector` (thread-safe) with StatusSet, StateSet, StateDelete, RelationSet, OpenPort, ClosePort |

### 2. Domain Service (`domain/scriptlet/`)

| File | Purpose |
|------|---------|
| `doc.go` | Package documentation |
| `service/service.go` | `State` interface, `WatcherFactory`, `Service`, `WatchableService` with `WatchScriptletApplications()`. Also declares `ApplicationService` interface for future config delegation. |
| `state/state.go` | `State` struct with placeholder `GetScriptletApplicationNames()` and `NamespaceForWatchScriptletApplications()` |

### 3. Wiring

| File | Change |
|------|--------|
| `domain/services/model.go` | Added `Scriptlet()` method + imports for `scriptletservice` and `scriptletstate` |
| `internal/services/interface.go` | Added `Scriptlet() *scriptletservice.WatchableService` to `ModelDomainServices` interface |
| `cmd/jujud-controller/agent/model/manifolds.go` | Registered `scriptletWorkerName` in `commonManifolds` with `ifNotMigrating` gate |

### 4. Runner Config

- Uses `worker.Runner` with `ShouldRestart: internalworker.ShouldRunnerRestart`
  (matches objectstore pattern)
- `IsFatal: false`, `RestartDelay: 10s`
- Runner is a child of the catacomb (passed in `Init`)

## Architecture

```
Model Manifolds
  └── scriptlet worker (catacomb + worker.Runner)
        ├── watches ScriptletService.WatchScriptletApplications()
        └── per-app child workers (applicationRunner)
              └── TODO: watch config/relation/lifecycle, dispatch to Starform
```

## Key Design Decisions

- **No deployer/uniter involvement**: Scriptlet apps bypass normal unit agent
  infrastructure
- **worker.Runner pattern**: Per-app workers auto-restart on failure
  (objectstore style)
- **Intent collection**: Scriptlets declare intent (status, state, relations,
  ports); applied transactionally after successful execution
- **Domain service separation**: `domain/scriptlet/` is independent from
  `domain/application/`; uses its own state for filtering scriptlet apps

## What's Left (TODO)

1. **State implementation**: Real DB query in
   `GetScriptletApplicationNames()` to filter apps by charm type/format
2. **Charm format extension**: Add scriptlet indicator to charm metadata so
   apps can be identified
3. **Per-app event dispatching**: Wire config watcher, relation watcher,
   lifecycle watcher into `applicationRunner`
4. **Starform integration**: Actually execute scriptlets and collect intents
5. **Intent application**: After scriptlet success, apply collected intents
   via domain services
6. **Deployer bypass**: Ensure deployer doesn't trigger for scriptlet
   applications

## Build Notes

- All packages pass `go vet` ✓
- Cannot do full `go build` in this environment (dqlite/sqlite3 build
  constraints)
- Import ordering verified with `gci`
