# Architecture

`djlt` is structured as a Go monorepo.

## Directory Structure

- `cmd/djlt/`: CLI entry points. Uses a **Verb-Centric** architecture.
- `internal/`: UI-agnostic core logic.
    - `cli/`: The **Surgical 6** verb implementations (`list`, `sync`, `make`, `move`, `remove`, `config`). Each verb is a single file; no provider-specific logic lives here.
    - `engine/`: Universal search and analysis engine (aliased as `library` in some contexts). Abstracted via the `Library` and `WritableLibrary` interfaces. Works exclusively with neutral `models`.
    - `models/`: Central domain models (`Track`, `ResourceGroup`, `Resource`) that provide a provider-agnostic language for the entire monorepo.
    - `provider/`: Thin coordination layer that translates CLI context to domain services. Houses specialized implementations for Rekordbox, Plex, and M3U.
    - `rekordbox/`: Domain core for Pioneer Rekordbox. Authority on XML parsing, identity resolution, custom query matching (cues/colors), and metadata updates.
    - `plex/`: Domain core for Plex Media Server. Authority on API interaction and mapping to neutral models.
    - `m3u/`: Domain core for M3U playlists. Authority on parsing and formatting.
    - `query/`: Lexer/Parser and Evaluator for the selection syntax.
    - `sync/`: Multi-threaded orchestration for data movement and XML injection.
    - `media/`: Parallel FFmpeg transcoding. Abstracted via the `sys.Runner` interface.
    - `sys/`: System abstractions for I/O and command execution.
- `pkg/`: Publicly accessible packages.
    - `rekordbox/`: XML types and a **High-Fidelity Formatter**. This package ensures that any XML modified by `djlt` is indistinguishable from one exported by Rekordbox, preserving idiosyncratic attribute ordering and wrapping rules.

## Core Concepts

### The Surgical 6

The CLI exposes exactly six top-level verbs. Legacy commands (`add`, `remove`, `rename`, `stat`) have been absorbed into these verbs via flags:

| Verb | Alias | Absorbed | Key Flag |
| :--- | :---- | :------- | :------- |
| `list` | `ls` | `stat` | `--stats` |
| `sync` | — | `add` | `--append` |
| `make` | `mk`, `create` | — | `--from` |
| `move` | `mv` | `rename` | `--name` |
| `remove` | `rm` | `remove` | `--from` |
| `config` | — | — | — |

The `remove` verb distinguishes two semantically different operations via `--from`:
- **Resource Deletion** (`rm rb/playlists name:Inbox`): permanently removes the ResourceGroup from the library.
- **Membership Removal** (`rm rb/tracks title:X --from "rb/playlists name:Inbox"`): unlinks tracks from a playlist without deleting them.

### Service-Oriented Domain (Track-Centric)

`djlt` uses a hierarchical service architecture to manage library data. This structure places the **Track** at the center of the domain model, reflecting that a track exists independently of its container.

- **`Tracks()`**: The primary service for music data.
    - `List()`: Query tracks.
    - `Update()`: Modify metadata.
    - **`Groups()`**: Sub-service for track organization.
        - `Add()` / `Remove()` / `Move()`: Manage memberships as an attribute of track organization.
- **`Groups()`**: Structural management for containers themselves.
    - `Create()` / `Delete()`: Manage the existence of Playlists and Folders.
    - `Update()`: Rename or move a group within the folder tree.
- **`System()`**: Global library maintenance (Save, Fix, Sync).

### The Implementation Mantra

> **"The Provider is a shell; the Package is the authority."**

To prevent **Leaky Abstractions**, `internal/provider` packages are restricted to:
1. Handling CLI-specific output (colors, progress bars).
2. Mapping execution context (e.g., `DryRun` flags).
3. Coordinating between the domain core and the query engine.

All implementation-specific intelligence—such as how to parse a Rekordbox cue point, how to map a Plex rating, or how to inject metadata into an XML stream—must live in the corresponding implementation package (`rekordbox`, `plex`, `m3u`).

### High-Fidelity XML Formatting
Rekordbox is sensitive to the structure of its XML. The `TokenStreamFormatter` in `internal/rekordbox` implements several rules to match this:
- **Attribute Ordering**: Attributes for `TRACK`, `NODE`, `POSITION_MARK`, etc., are sorted according to a specific schema.
- **Smart Wrapping**: Attributes are wrapped onto new lines only when the total decoded line length exceeds 88 characters.
- **Self-Closing Tags**: Empty tags are automatically converted to self-closing format (`<TAG/>`).
- **Entity Encoding**: Specific character escaping rules for ampersands, quotes, and other special characters.

### Selection Engine
The selection engine uses a recursive descent parser and a universal evaluator. It supports:
- Boolean logic: `&&`, `||`, `!`
- Numeric operators: `>`, `<`, `..` (range), `:` (exact equality for numeric fields)
- Exact match: `=` (case-sensitive)
- Regex match: `::`
- Field mapping: `bpm`, `rating`, `hotcues`, `playlists`, etc.

### Test Boundaries and Interfaces
`djlt` uses explicit interfaces to decouple core logic from external dependencies:
- **`Library`**: Decouples implementation-specific storage from the query `Engine`, allowing for agnostic search across different data formats.
- **`Provider`**: Service-registry interface that exposes the hierarchical `Tracks()`, `Groups()`, and `System()` domains.
- **`Resource`**: A universal interface for any item in a library (Track, ResourceGroup), allowing for generic movement and listing logic.
- **`sys.FileSystem` & `sys.Runner`**: Abstract the OS environment (Filesystem, FFmpeg), enabling side-effect-free testing of media and sync operations.

### Cobra Command Factory
Each verb is created by a `newXxxCmd()` constructor function. Flag variables are captured in closures, so each command instance carries its own isolated state. `NewRootCmd()` wires all constructors together and is the sole entry point used by both the production binary (`var RootCmd = NewRootCmd()`) and tests (`root := NewRootCmd()`).

The four persistent flags (`--file`, `--dry-run`, `--verbose`, `--json`) are still bound to package-level vars shared across all verbs. Tests reset these four vars in `resetTestState()` before creating a new root command. No `pflag.Changed` traversal is needed because each `NewRootCmd()` call produces a fresh flag set.
