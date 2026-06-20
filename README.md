# DJ Library Tools (djlt)

A unified suite of high-performance tools for managing DJ libraries, Rekordbox XML exports, and playlist hygiene.

## 🛠 Installation

Currently, the recommended way to install `djlt` is via `go install`. 

```bash
# Install the tool
go install github.com/llttlltt/dj-library-tools/cmd/djlt@latest
```

## 🚀 Usage

`djlt` provides specialized commands for different library management tasks.

### 🎵 Playlist Hygiene

The `playlist fix` command is a powerful tool to migrate, repair, and enrich M3U/M3U8 playlists. It replaces fragile Bash scripts with a native Go implementation that requires zero external dependencies (no FFmpeg/ffprobe needed).

**Common Workflows:**

#### Migrate to MP3 (Beets Workflow)
If you've converted your library from FLAC to MP3, update your playlists to point to the new files and add Rekordbox-ready metadata:
```bash
djlt playlist fix my_playlist.m3u --ext mp3 --m3u8
```

#### Smart Priority Search & Pruning
Search for files in a priority order (e.g., prefer MP3, fallback to FLAC). If neither is found, the track is automatically pruned from the playlist to ensure zero "File Not Found" errors in Rekordbox:
```bash
djlt playlist fix my_playlist.m3u --ext mp3,flac --m3u8 -v
```

#### Standardize & Enrich
Upgrade a simple list of paths to a formal `.m3u8` with Artist and Title metadata:
```bash
djlt playlist fix paths.m3u --m3u8 -o clean_library.m3u8
```

**Flags:**
- `-e, --ext strings`: Priority list of extensions to search for (e.g., `mp3,flac,wav`).
- `--m3u8`: Enrich output with `#EXTINF` metadata (Artist - Title).
- `-o, --output string`: Specific path for the output file.
- `-r, --remove-original`: Prompt to remove the source file after successful processing.
- `-v, --verbose`: Enable granular logging of every resolved path.
- `-f, --force`: Overwrite output file if it already exists.

---

### 📂 Metadata Operations (Legacy `rb-cli`)

#### Move Metadata
Matches tracks between two Rekordbox XML files and synchronizes specific fields (currently `Tempo`).

```bash
djlt metadata move --source source.xml --destination target.xml --output merged.xml
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
