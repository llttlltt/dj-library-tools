# Architecture Protocols

## The Core Mantra
>
> **"The Provider is a shell; the Package is the authority."**

This means that `internal/providers` packages should only handle orchestration, and mapping execution context. All domain-specific intelligence (XML manipulation, API client logic, metadata reconciliation, color mapping, and hierarchical cue matching) is co-located in the provider package.

## Layered Architecture

The system follows a strict layered architecture with a one-way dependency flow:

`ui ──> services ──> providers ──> core`
(everyone may import `core`; `core` imports nothing under `internal/`)

1. **Core** (`internal/core/`): Pure domain. Models, query engine, shared errors, location parsing. Zero knowledge of providers, formats, infra, or presentation.
2. **Infra** (`internal/infra/`): Adapters to external processes/OS (ffmpeg, system calls).
3. **Providers** (`internal/providers/`): Library sources and native format authority (co-located). Registered via factory.
4. **Services** (`internal/services/`): UI-agnostic business logic and orchestration (resolver, sync engine, orchestrator facade).
5. **UI** (`internal/ui/`): Presentation layer. `cli` and `gui`.

## Track-Centric Service Design

The system follows a nested, resource-oriented service structure:

- **`Tracks()`**: The primary entry point for managed music data.
- **`Tracks().Groups()`**: Handles memberships (Add/Remove/Move).
- **`Groups()`**: Handles structural containers (Create/Delete Playlists and Folders).
- **`System()`**: Handles global maintenance (Save, Fix, Sync).

## Hard Boundaries

- **Models as Source of Truth**: Neutral models (`Track`, `ResourceGroup`) are the sole authorities on their queryable data. They must implement a `Value(key string) string` method to represent their properties, driven by the central registry in `internal/core/models/metadata.go`.
- **Universal Field Registry**: All queryable and displayable fields must be defined in `internal/core/models/metadata.go`. This registry links field names to types, accessors, and required provider capabilities. Beatgrids, cues, and ratings are deliberate domain capabilities shared across DJ tools, not Rekordbox-specific leaks.
- **Deep Discovery**: Provider methods for resource selection or modification (e.g., `UpdateGroup`, `CreateGroup`) must be recursive by default. Searching for a container by name should span the entire hierarchy unless explicitly restricted.
- **Generic Query Logic**: The `internal/core/query` package must remain 100% generic logic. It must not have knowledge of specific fields like "Artist" or "BPM."
- **Implementation Authority**: Implementation packages (`rekordbox`, `plex`, `m3u`) are the sole authorities on their data formats. They handle their own mapping to neutral models.
- **Hierarchical Path Resolution**: Complex metadata (Cues, Markers, Playlists) must be queried via a standardized Path Resolver.
- **Strict Agnosticism**: Core packages must NEVER import specific implementation packages or anything under `internal/ui`.
- **Provider Registry**: All providers must self-register via `init()` in their respective packages under `internal/providers/`. Registration must include static `ProviderCapabilities` and granular `ResourceInfo` (defining `CanWrite` and `SupportsQuery` per resource).
- **Static Discovery**: Provider capabilities (read/write status, supported resources) must be discoverable via the factory registry without requiring provider instantiation or valid configuration. This allows UIs to validate user input and adapt selectors proactively.
- **Persistence Responsibility**: The **UIs** (CLI/GUI) are the sole authorities on persistence. Provider methods perform modifications in-memory; the caller must explicitly call `Save()` to commit changes, typically orchestrated via the `orchestrator` service.
- **Safe-by-Default**: All mutating operations must be non-destructive by default. The `--apply` flag (CLI) or "Apply" button (GUI) is the universal gatekeeper for the `Save()` operation.
- **UI Decoupling**: Core, infra, providers, and services must not import UI libraries or write directly to Stdout/Stderr. All user feedback must be channeled through the `Feedback` interface.
- **Orchestrator Facade**: All UI interactions flow through the `internal/services/orchestrator`. It is the single seam between presentation and logic, owning statistics computation, sorting, and default table columns. It returns inert data only.
- **Workflow Engine**: Multi-step operations are orchestrated by `internal/services/workflow`. It is a higher-level consumer of the orchestrator, capable of parallel execution and dependency management.
- **Inert Results**: `Orchestrator.List` returns a `ListResult` containing pure data and metadata (like `DefaultColumns`). It never returns a mutable `Provider` handle to the UI.
- **Option Ownership**: The orchestrator defines its own option DTOs (e.g., `SyncOptions`, `FixOptions`). UIs construct these types, and the orchestrator maps them to internal provider types.
- **Source-based Configuration**: Configuration is decentralized into individual JSON files for Sources, Workflows, and Path Maps. The app resolves these artifacts from the filesystem by UUID.
- **Context Threading**: `context.Context` must flow from every UI call through the orchestrator into `resolver.ResolveSelection` and from there into all provider calls. Provider list operations must never use `context.Background()` internally; doing so silently drops cancellation signals from the caller.

## GUI State Management

The GUI follows a **Reactive Store** pattern using **Effect Atoms** (`@effect-atom/atom`). 

- **Single Source of Truth**: Global domain state (Sources, Workflows, Provider Metadata) is stored in centralized Atoms in `src/store/`.
- **Reactive Subscriptions**: Components must consume state via the `useAtom` hook. Redundant localized `useState` for global artifacts is an anti-pattern.
- **Side Effect Encapsulation**: Data fetching and mutations (Wails binding calls) are encapsulated in Effect generators within the store.
- **Managed Runtime**: The frontend uses a global `ManagedRuntime` and `AtomRegistry` (setup in `src/lib/runtime.ts`) to provide a high-integrity execution context for all side effects.
- **Stateless Views**: Views and components should be "stateless observers" that trigger Store effects and render Atom values. 
