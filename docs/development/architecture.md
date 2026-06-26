# Architecture

`djlt` is structured as a Go monorepo.

## Directory Structure

- `cmd/djlt/`: CLI entry points. Uses a **Verb-Centric** architecture.
- `internal/`: UI-agnostic core logic.
    - `engine/`: Query engine and library orchestration. Abstracted via the `Library` interface.
    - `provider/`: Unified interface for library sources (Rekordbox, Plex).
    - `plex/`: Parallel prober and API client for Plex.
    - `query/`: Lexer/Parser for selection syntax.
    - `sync/`: Orchestration for data movement and XML injection.
    - `media/`: Parallel FFmpeg transcoding. Abstracted via the `sys.Runner` interface.
    - `sys/`: System abstractions for I/O and command execution.
- `pkg/`: Publicly accessible packages.
    - `rekordbox/`: XML types and a **High-Fidelity Formatter**. This package ensures that any XML modified by `djlt` is indistinguishable from one exported by Rekordbox, preserving idiosyncratic attribute ordering and wrapping rules.

## Core Concepts

### High-Fidelity XML Formatting
Rekordbox is sensitive to the structure of its XML. The `TokenStreamFormatter` in `pkg/rekordbox` implements several rules to match this:
- **Attribute Ordering**: Attributes for `TRACK`, `NODE`, `POSITION_MARK`, etc., are sorted according to a specific schema.
- **Smart Wrapping**: Attributes are wrapped onto new lines only when the total decoded line length exceeds 88 characters.
- **Self-Closing Tags**: Empty tags are automatically converted to self-closing format (`<TAG/>`).
- **Entity Encoding**: Specific character escaping rules for ampersands, quotes, and other special characters.

### Selection Engine
The selection engine uses a recursive descent parser. It supports:
- Boolean logic: `&&`, `||`, `!`
- Numeric operators: `>`, `<`, `..` (range), `:` (exact equality for numeric fields)
- Field mapping: `bpm`, `rating`, `hotcues`, `playlistcount`, etc.

### Test Boundaries and Interfaces
`djlt` uses explicit interfaces to decouple core logic from external dependencies:
- **`Library`**: Decouples the `Engine` from the concrete Rekordbox XML structure, allowing for high-speed in-memory testing.
- **`Provider`**: Unifies different sources (Rekordbox, Plex) into a single queryable interface.
- **`sys.FileSystem` & `sys.Runner`**: Abstract the OS environment (Filesystem, FFmpeg), enabling side-effect-free testing of media and sync operations.
