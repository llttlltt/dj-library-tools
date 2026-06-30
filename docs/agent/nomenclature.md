# Domain & Nomenclature

## Core Terminology

- **Provider**: A music library source (Rekordbox, Plex, M3U).
- **Resource**: A type of data served by a provider (e.g., tracks, groups).
- **ResourceGroup**: A logical container for items.
  - **Playlist**: A group that contains tracks.
  - **Folder**: A container for other groups (hierarchical).
- **Location**: A composite identifier (`provider/resource`).
- **Query**: Criteria used to filter resources.
- **Selection**: A resolved, inert set of resources (tracks/groups) ready for an operation. It carries data only; it cannot trigger mutations.
- **Orchestrator**: The UI-agnostic facade (`internal/services/orchestrator`) that every UI calls to perform an operation — the single seam between presentation and logic. It owns resolution, sort-field validation, statistics, sorting, and persistence gating (on `Apply`), invokes provider methods, and returns inert results.
- **ListResult**: The inert result struct from `Orchestrator.List` — resource data plus presentation metadata (e.g. `DefaultColumns`). It never carries a mutable `Provider` handle.
- **Feedback**: The interface through which all user-facing output flows (`OnStatus`, `OnPreview`, `OnSuccess`, `OnWarning`, `OnTable`, `OnProgress`). The CLI implements it for the terminal; a future GUI implements it for its widgets. Core/services/providers never print directly.
- **Metadata Aspect**: Optional track attributes (Beatgrids, BPM, Rating) reconciled during an operation.
- **Join**: Agnostic identity matching between tracks from different providers.
- **Path Querying**: Hierarchical traversal and statistical analysis of deep metadata (Cues, Markers, Playlists) using the `Collection.Index/Property-Stat` convention.

## Domain Standards

- **Unified Rating Scale**: Use 0-255 globally to normalize between rating systems.
- **Resource Identity**:
  - **Tracks**: Unique persistent integer or UUID provided by the source.
  - **Groups (Hierarchical)**: The **Canonical Path** (e.g., `Folder/Subfolder/Playlist`) is the unique ID. Implementation packages must support finding resources by both this Path ID and their base Name for legacy compatibility.
  - **Identity Authority**: The provider is the sole authority for generating and resolving IDs. Core services should treat IDs as opaque strings.
