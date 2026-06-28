# Architecture Protocols

## The Core Mantra
> **"The Provider is a shell; the Package is the authority."**

This means that `internal/provider` packages should only handle orchestration, CLI feedback (colors, progress bars), and mapping execution context. All domain-specific intelligence (XML manipulation, API client logic, metadata reconciliation, color mapping) must live in the core implementation package (e.g., `internal/rekordbox`).

## Track-Centric Service Design
The system follows a nested, resource-oriented service structure:
- **`Tracks()`**: The entry point for managed music data.
- **`Tracks().Groups()`**: Handles memberships (Add/Remove/Move). Membership is a property of track organization.
- **`Groups()`**: Handles structural containers (Create/Delete Playlists and Folders).
- **`System()`**: Handles global maintenance (Save, Fix, Sync).

## Hard Boundaries
- **Strict Agnosticism**: Core packages (`models`, `library`, `query`, `sync`, `utils`) must NEVER import specific implementation packages.
- **Implementation Authority**: Implementation packages (`rekordbox`, `plex`, `m3u`) are the sole authorities on their data formats. They must handle their own mapping to neutral models (e.g., `ToNeutralTrack`).
- **Provider Registry**: All providers must self-register via `init()` in their respective packages under `internal/provider/`.
- **Discovery-Driven CLI**: The CLI must interact through standardized services. Avoid type-assertions to specific providers where possible.
- **Execution Context**: Always pass `provider.ExecutionContext` to respect runtime flags like `DryRun`.
- **UI Decoupling**: Core packages must not import UI libraries (e.g., `mpb`, `color`). These belong in the Provider or CLI layer.
