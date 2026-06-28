# Architecture Protocols

## The Core Mantra
> **"The Provider is a shell; the Package is the authority."**

This means that `internal/provider` packages should only handle orchestration, CLI feedback, and mapping execution context. All domain-specific intelligence (XML manipulation, API client logic, metadata reconciliation, color mapping, and hierarchical cue matching) must live in the core implementation package (e.g., `internal/rekordbox`).

## Track-Centric Service Design
The system follows a nested, resource-oriented service structure:
- **`Tracks()`**: The primary entry point for managed music data.
- **`Tracks().Groups()`**: Handles memberships (Add/Remove/Move).
- **`Groups()`**: Handles structural containers (Create/Delete Playlists and Folders).
- **`System()`**: Handles global maintenance (Save, Fix, Sync).

## Hard Boundaries
- **Models as Source of Truth**: Neutral models (`Track`, `ResourceGroup`) are the sole authorities on their queryable data. They must implement a `Value(key string) string` method to represent their properties, driven by the central registry in `internal/models/metadata.go`.
- **Universal Field Registry**: All queryable and displayable fields must be defined in `internal/models/metadata.go`. This registry links field names to types, accessors, and required provider capabilities.
- **Generic Query Logic**: The `internal/query` package must remain 100% generic logic. It must not have knowledge of specific fields like "Artist" or "BPM."
- **Implementation Authority**: Implementation packages (`rekordbox`, `plex`, `m3u`) are the sole authorities on their data formats. They must handle their own mapping to neutral models and any complex custom matching (via `CustomMatch`).
- **Strict Agnosticism**: Core packages (`models`, `library`, `query`, `sync`, `utils`) must NEVER import specific implementation packages.
- **Provider Registry**: All providers must self-register via `init()` in their respective packages under `internal/provider/`.
- **Discovery-Driven CLI**: The CLI must interact through standardized services. Avoid type-assertions to specific providers where possible.
- **Execution Context**: Always pass `provider.ExecutionContext`. The system is **Safe-by-Default**; operations must only be persisted if `ctx.Apply` is explicitly true.
- **UI Decoupling**: Core packages must not import UI libraries or write directly to Stdout. All user feedback must be channeled through the `Feedback` interface in the `ExecutionContext`.
- **Provider-Driven UI**: The CLI should remain a thin, dynamic wrapper. It must use provider services (like `TableHeaders()`) to determine its presentation logic rather than hardcoding provider-specific behavior.
