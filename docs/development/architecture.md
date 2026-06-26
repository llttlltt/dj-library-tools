# Architecture

`djlt` is structured as a Go monorepo.

## Directory Structure

- `cmd/djlt/`: CLI entry points. Uses a **Verb-Centric** architecture.
- `internal/`: UI-agnostic core logic.
    - `engine/`: Universal search and analysis engine. Abstracted via the `Library` and `WritableLibrary` interfaces. Works exclusively with neutral `models`.
    - `models/`: Central domain models (`Track`, `Node`, `Resource`) that provide a provider-agnostic language for the entire monorepo.
    - `provider/`: Capability-based plugins for library sources (Rekordbox, Plex).
    - `plex/`: API client and models for Plex Media Server.
    - `query/`: Lexer/Parser and Evaluator for the selection syntax.
    - `sync/`: Multi-threaded orchestration for data movement and XML injection.
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
The selection engine uses a recursive descent parser and a universal evaluator. It supports:
- Boolean logic: `&&`, `||`, `!`
- Numeric operators: `>`, `<`, `..` (range), `:` (exact equality for numeric fields)
- Exact match: `=` (case-sensitive)
- Regex match: `::`
- Field mapping: `bpm`, `rating`, `hotcues`, `playlistcount`, etc.

### Test Boundaries and Interfaces
`djlt` uses explicit interfaces to decouple core logic from external dependencies:
- **`Library`**: Decouples the `Engine` from the concrete data structure, allowing for high-speed in-memory testing using neutral `models`.
- **`Provider`**: Capability-based interface that unifies different sources (Rekordbox, Plex) into a single queryable and writable interface.
- **`Resource`**: A universal interface for any item in a library (Track, Node), allowing for generic movement and listing logic.
- **`sys.FileSystem` & `sys.Runner`**: Abstract the OS environment (Filesystem, FFmpeg), enabling side-effect-free testing of media and sync operations.
