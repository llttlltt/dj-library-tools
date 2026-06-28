# Architecture

`djlt` is structured as a Go monorepo.

## Directory Structure

- `cmd/djlt/`: CLI entry points. Uses a **Verb-Centric** architecture.
- `internal/`: UI-agnostic core logic.
    - `cli/`: The **Surgical 6** verb implementations (`list`, `sync`, `make`, `move`, `remove`, `config`). Each verb is a single file; no provider-specific logic lives here.
    - `engine/`: Universal search and analysis engine. Abstracted via the `Library` and `WritableLibrary` interfaces. Works exclusively with neutral `models`.
    - `models/`: Central domain models (`Track`, `ResourceGroup`, `Resource`) that provide a provider-agnostic language for the entire monorepo.
    - `provider/`: Capability-based plugins for library sources (Rekordbox, Plex).
    - `plex/`: API client and models for Plex Media Server.
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
- **`Library`**: Decouples the `Engine` from the concrete data structure, allowing for high-speed in-memory testing using neutral `models`.
- **`Provider`**: Capability-based interface that unifies different sources (Rekordbox, Plex) into a single queryable and writable interface.
- **`Resource`**: A universal interface for any item in a library (Track, ResourceGroup), allowing for generic movement and listing logic.
- **`sys.FileSystem` & `sys.Runner`**: Abstract the OS environment (Filesystem, FFmpeg), enabling side-effect-free testing of media and sync operations.

### Cobra Command Factory
Each verb is created by a `newXxxCmd()` constructor function. Flag variables are captured in closures, so each command instance carries its own isolated state. `NewRootCmd()` wires all constructors together and is the sole entry point used by both the production binary (`var RootCmd = NewRootCmd()`) and tests (`root := NewRootCmd()`).

The four persistent flags (`--file`, `--dry-run`, `--verbose`, `--json`) are still bound to package-level vars shared across all verbs. Tests reset these four vars in `resetTestState()` before creating a new root command. No `pflag.Changed` traversal is needed because each `NewRootCmd()` call produces a fresh flag set.
