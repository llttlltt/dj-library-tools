# djlt: the DJ’s library engine

`djlt` is a unified suite of high-performance tools for managing DJ libraries, Rekordbox XML exports, and cross-platform synchronization. 

It features a powerful boolean selection engine, parallel Plex-to-Rekordbox synchronization, and full CRUD operations on the Rekordbox playlist tree.

## Documentation

Full documentation is available at **[https://llttlltt.github.io/dj-library-tools/](https://llttlltt.github.io/dj-library-tools/)**.

- **[Getting Started](https://llttlltt.github.io/dj-library-tools/guide/getting-started/)**: Installation and first steps.
- **[Configuration](https://llttlltt.github.io/dj-library-tools/guide/configuration/)**: Setting up your paths and Plex connection.
- **[Syntax](https://llttlltt.github.io/dj-library-tools/guide/syntax/)**: Mastering the query language and operators.
- **[CLI Reference](https://llttlltt.github.io/dj-library-tools/commands/)**: Comprehensive guide to every command and flag.

## Installation

```bash
go install github.com/llttlltt/dj-library-tools/cmd/djlt@latest
```

*Note: Requires [Go](https://go.dev/doc/install). If you plan on transcoding files, you also need [FFmpeg](https://ffmpeg.org/).*

## Quick Start

```bash
# 1. Set your master Rekordbox XML path
djlt config rekordbox.xml-path ~/Documents/rekordbox.xml

# 2. Authenticate with Plex (Optional)
djlt auth --plex

# 3. Query your library
djlt ls rb/tracks "genre:House && bpm:124..128"

# 4. Sync a Plex playlist to Rekordbox
djlt sync plex/playlists name:Summer --to "rb/playlists name:'Plex Sync'"
```

## Contributing

Contributions are welcome! Please see the **[Architecture](https://llttlltt.github.io/dj-library-tools/development/architecture/)** guide for technical details on the codebase.

## Early Alpha Warning

`djlt` is currently in **early alpha**. It is functional but under active development. **Always backup your Rekordbox XML library** before performing write operations.

## License

MIT
