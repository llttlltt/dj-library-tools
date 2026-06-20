# DJ Library Tools (djlt)

A unified suite of high-performance tools for managing DJ libraries, Rekordbox XML exports, and playlist hygiene.

## 🛠 Installation

Currently, the recommended way to install `djlt` is via `go install`. 

```bash
# Install the tool
go install github.com/llttlltt/dj-library-tools/cmd/djlt@latest
```

## 🚀 Getting Started

`djlt` uses a persistent configuration stored in `~/.config/djlt/config.json`. Initialize your setup:

```bash
# Authenticate with Plex
djlt auth plex

# Configure your server (optional if using auto-discovery)
djlt config plex --host 10.0.10.151 --port 32400

# Map remote Plex paths to local mount points
djlt config plex --map "/media/Music/Master:/Volumes/Media/Music/Master"

# Set your master Rekordbox XML path
djlt config rekordbox --xml "/path/to/rekordbox.xml"
```

---

## 🔍 Discovery & Querying

`djlt` uses a **Location-based URI** syntax: `provider/resource:query`.

### Listing Items (`list`)
```bash
# List Plex playlists
djlt list plex

# List Plex tracks matching a query
djlt list plex/tracks:9102

# List Rekordbox tracks matching a query
djlt list rb:bpm:120..128
```

### Selection Syntax
| Operator | Type | Example |
| :--- | :--- | :--- |
| `:` | Substring | `artist:Four` |
| `=` | Exact | `artist="Four Tet"` |
| `::` | Regex | `name::"^01"` |
| `..` | Range | `bpm:124..128` |
| `!` | Negation | `!genre:Techno` |

---

## 🔄 Synchronization & Export

The `sync` command orchestrates data and media movement between providers.

### Sync Plex to Rekordbox (With Transcoding)
Automatically matches Plex tracks to Rekordbox metadata, transcodes files via FFmpeg, and updates the Rekordbox XML.

```bash
djlt sync plex:MyPlaylistName rb \
  --dest ./ExportFolder \
  --format mp3
```

- **Matching**: Uses fuzzy logic to pair tracks.
- **Transcoding**: Inherits settings from Beets (320k MP3, ID3v2.3).
- **XML Injection**: Creates a "Plex Sync" folder in Rekordbox with your playlist.

---

## 🎵 Playlist Hygiene

The `playlist fix` command repairs and enriches local M3U/M3U8 playlists.

```bash
# Fix extensions and upgrade to M3U8 with metadata
djlt playlist fix my_playlist.m3u --ext mp3,flac --m3u8
```

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
