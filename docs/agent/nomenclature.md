# Domain & Nomenclature

## Core Terminology

- **Provider**: A music library source (Rekordbox, Plex, M3U).
- **Resource**: A type of data served by a provider (e.g., tracks, playlists, folders).
- **ResourceInfo**: Static metadata about a resource (Name, CanWrite, SupportsQuery) used for UI validation.
- **ResourceGroup**: A logical container for items.
  - **Playlist**: A group that contains tracks.
  - **Folder**: A container for other groups (hierarchical).
- **Location**: A composite identifier (`provider/resource` or `source_uuid/resource`). The CLI uses provider prefixes by convention; the GUI uses explicit Source UUIDs.
- **Source**: A user-named, configured instance of a provider (e.g., "Main Library" → Rekordbox, "Plex Home"). Each Source has a unique UUID and connection details (file path, host, token).
- **Workflow**: A named, user-defined collection of Steps. Workflows are the primary unit of library automation in the GUI.
- **Step**: An atomic operation within a Workflow (e.g., sync, fix, edit). Supports parallel execution and sequential dependencies via the `after` field. One Step has one source Endpoint and can fan out to multiple target Endpoints.
- **Endpoint**: The `{source_id, resource, query}` triple identifying one side of a Step. The `resource` must be validated against the provider's static registry.
- **Path Map**: A declared path-translation relationship between two Sources (e.g., translating Rekordbox local paths to Plex server paths).
- **Query**: Criteria used to filter resources.
- **Selection**: A resolved, inert set of resources (tracks/groups) ready for an operation. It carries data only; it cannot trigger mutations.
- **Orchestrator**: The UI-agnostic facade (`internal/services/orchestrator`) that every UI calls to perform an operation — the single seam between presentation and logic. It owns resolution, sort-field validation, statistics, sorting, and persistence gating (on `Apply`), invokes provider methods, and returns inert results.
- **ListResult**: The inert result struct from `Orchestrator.List` — resource data plus presentation metadata (e.g. `DefaultColumns`). It never carries a mutable `Provider` handle.
- **Feedback**: The interface through which all user-facing output flows (`OnStatus`, `OnPreview`, `OnSuccess`, `OnWarning`, `OnTable`, `OnProgress`). The CLI implements it for the terminal; the GUI implements it for its widgets. Core/services/providers never print directly.
- **Metadata Aspect**: Optional track attributes (Beatgrids, BPM, Rating) reconciled during an operation.
- **Join**: Agnostic identity matching between tracks from different providers.
- **Path Querying**: Hierarchical traversal and statistical analysis of deep metadata (Cues, Markers, Playlists) using the `Collection.Index/Property-Stat` convention. Supports both track-level collections and group-content traversal (`tracks/title`).

## Domain Standards

- **Unified Rating Scale**: Use 0-255 globally to normalize between rating systems.
- **Resource Identity**:
  - **Tracks**: Unique persistent integer or UUID provided by the source.
  - **Groups (Hierarchical)**: The **Canonical Path** (e.g., `Folder/Subfolder/Playlist`) is the unique ID. Implementation packages must support finding resources by both this Path ID and their base Name for legacy compatibility.
  - **Identity Authority**: The provider is the sole authority for generating and resolving IDs. Core services should treat IDs as opaque strings.
- **Source Resolution**: The resolver is the authority for mapping location prefixes to concrete provider configurations. It handles both hardcoded provider names (by convention) and explicit Source UUIDs.
- **Pure Domain**: The `internal/core` packages are restricted to pure logic. They must never import I/O-capable standard library packages (`os`, `net`, `syscall`). Any environmental state (like file existence) must be hydrated by the provider layer.
