# DJ Library Tools (djlt)

A unified suite of high-performance tools for managing DJ libraries, Rekordbox XML exports, and playlist hygiene.

## 🛠 Installation

Currently, the recommended way to install `djlt` is via `go install`. 

```bash
# Install the tool
go install github.com/llttlltt/dj-library-tools/cmd/djlt@latest
```

## 🚀 Usage

`djlt` provides specialized commands for different library management tasks. All library commands require a path to your Rekordbox XML via the global `-x, --xml` flag.

### 🔍 Library Exploration & Analysis

The core of `djlt` is a powerful **Selection Engine** that allows you to filter your library using a natural syntax.

#### List Tracks (`ls`)
Search and display tracks matching specific criteria:
```bash
djlt ls "artist:Four Tet bpm:120..128" --xml library.xml
```

#### Library Statistics (`stat`)
Generate a summarized report of your library or a filtered selection:
```bash
djlt stat "genre:House" --xml library.xml
```
Provides: Total tracks, Average BPM, and Top 5 Artists/Genres/Keys.

---

### 🎵 Playlist Hygiene

The `playlist fix` command is a powerful tool to migrate, repair, and enrich M3U/M3U8 playlists.

**Common Workflows:**

#### Smart Priority Search & Pruning
Search for files in a priority order (e.g., prefer MP3, fallback to FLAC). If neither is found, the track is automatically pruned to ensure zero "File Not Found" errors in Rekordbox:
```bash
djlt playlist fix my_playlist.m3u --ext mp3,flac --m3u8 -v
```

#### Batch Processing & Dry Run
Process multiple playlists at once and preview the results safely:
```bash
djlt playlist fix ./Playlists/*.m3u --ext mp3 --m3u8 --dry-run
```

---

### 📂 Metadata Operations

#### Move Metadata (Legacy Sync)
Matches tracks between two Rekordbox XML files and synchronizes `Tempo` fields.
```bash
djlt metadata move --source source.xml --destination target.xml --output merged.xml
```

---

## 🎯 Selection Syntax

Most commands support a query string as an argument.

| Operator | Type | Example |
| :--- | :--- | :--- |
| `:` | Substring | `artist:Four` |
| `=` | Exact | `artist="Four Tet"` |
| `::` | Regex | `name::"^01"` |
| `..` | Range | `bpm:124..128` |
| `!` | Negation | `!genre:Techno` |

**Supported Fields:** `name`, `artist`, `album`, `genre`, `bpm`, `key`, `label`, `rating`, `playcount`, `added`, `kind`, `size`.

---

## 🛠 Development

This project uses [mise](https://mise.jdx.dev/) for toolchain management.

```bash
# Setup environment
mise install

# Run tests
go test ./...

# Build local binary
go build -o bin/djlt ./cmd/djlt
```
