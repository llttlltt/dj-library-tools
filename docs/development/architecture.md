# Architecture

`djlt` is structured as a Go monorepo.

## Directory Structure

- `cmd/djlt/`: CLI entry points. Uses a **God Command** pattern (one command per resource, flags as verbs).
- `internal/`: UI-agnostic core logic.
    - `engine/`: Query engine and Rekordbox tree management.
    - `plex/`: Parallel prober and API client for Plex.
    - `query/`: Lexer/Parser for selection syntax.
    - `sync/`: Orchestration for data movement and XML injection.
    - `media/`: FFmpeg transcoding and path sanitization.
- `pkg/`: Publicly accessible packages (e.g., `rekordbox` XML types).

## Core Concepts

### Selection Engine
The selection engine uses a recursive descent parser. It supports:
- Boolean logic: `&&`, `||`, `!`
- Numeric operators: `>`, `<`, `..` (range), `:` (exact equality for numeric fields)
- Field mapping: `bpm`, `rating`, `hotcues`, `playlistcount`, etc.

### God Command Pattern
All commands follow a consistent pattern:
`djlt [resource] --[verb] [arguments]`

Example:
`djlt playlist --new "My New Playlist"`
