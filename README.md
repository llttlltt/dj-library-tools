# DJ Library Tools (djlt)

A unified suite of high-performance tools for managing DJ libraries, Rekordbox XML exports, and cross-platform synchronization.

## Installation

```bash
go install github.com/llttlltt/dj-library-tools/cmd/djlt@latest
```

## Getting Started

```bash
# Authenticate with Plex
djlt auth plex

# Map remote Plex paths to local mount points (required when syncing from Plex)
djlt config plex --map "/media/Music:/Volumes/Music"

# Remove a path map
djlt config plex --remove-map "/media/Music"

# Set your master Rekordbox XML path
djlt config rekordbox --xml "~/Documents/rekordbox.xml"

# Show all current settings
djlt config
```

---

## Discovery & Querying

`djlt` features a powerful selection engine for filtering your library. The syntax follows a `provider/resource:query` pattern.

### Command Format
```bash
# Use single quotes to prevent shell interference with !, (, and )
djlt list [source] '[query]'
```

### Sources
| Source | Description |
| :--- | :--- |
| `rb/tracks` | All tracks in your Rekordbox Collection |
| `rb/playlists` | Playlist nodes in Rekordbox (Type=1) |
| `rb/folders` | Folder nodes in Rekordbox (Type=0) |
| `plex/playlists` | All Plex playlists |
| `plex/tracks` | Tracks from a specific Plex playlist ID |

### Boolean Operators
- `&&` / `AND` — both conditions must be met
- `||` / `OR` — either condition must be met
- `!` — negation
- `( )` — grouping

### Comparison Operators
| Operator | Type | Example |
| :--- | :--- | :--- |
| `:` | Substring (text) / Exact (numeric) | `artist:Four`, `playlistcount:3` |
| `=` | Exact | `artist="Four Tet"` |
| `::` | Regex | `name::"^01"` |
| `..` | Range | `bpm:124..128` |
| `>`, `>=` | Numeric greater-than | `bpm:>=128` |
| `<`, `<=` | Numeric less-than | `rating:<3` |

> **Note:** For numeric fields (`bpm`, `rating`, `playlistcount`, etc.) the `:` operator
> performs exact numeric equality, not substring matching. Use `..` for ranges or `>`/`<`
> for comparisons.

### Query Fields

#### Metadata
`name`, `artist`, `album`, `genre`, `key`, `year`, `label`, `comment`, `remixer`, `mix`

#### Technical
`bpm`, `bitrate`, `kind` (e.g. `MP3 File`), `size`, `time` (duration in seconds)

#### Library State
`rating` (0–5 stars), `playcount`, `added` (date string), `playlistcount`, `playlist`

#### Playlist & Folder Nodes (`rb/playlists:`, `rb/folders:`)
`name`, `folder` (parent folder name, `""` = root level), `entries` (track count), `type` (`0`=folder, `1`=playlist)

#### Cues & Beatgrids
| Field | Description |
| :--- | :--- |
| `beatgrids` | Number of TEMPO markers — e.g. `beatgrids:>1` |
| `hotcues` | Count — `hotcues:8`, or color check — `hotcues:aqua` |
| `memorycues` | Count — `memorycues:2`, or color check — `memorycues:pink` |
| `hotcue:[a-h]` | Target a specific slot — `hotcue:a:green` |
| `memorycue:[n]` | Target by position, high-to-low — `memorycue:1:loop` |

**Cue sub-properties:** `label:[text\|empty]`, `loop`, `pink`, `red`, `orange`, `yellow`, `green`, `aqua`, `blue`, `purple`, `none`

### Examples
```bash
# High-rated House tracks between 124–128 BPM
djlt list rb "genre:House && rating:>=4 && bpm:124..128"

# Tracks in both a specific playlist and a BPM range
djlt list rb "playlist:Summer && bpm:120..130"

# Tracks appearing on more than 3 playlists
djlt list rb "playlistcount:>3"

# Tracks in exactly 0 playlists (orphans)
djlt list rb "playlistcount:0"

# Tracks with multiple beatgrid markers (transition tracks)
djlt list rb "beatgrids:>1"

# Tracks where Hot Cue B is Aqua and labeled "INTRO"
djlt list rb "hotcue:b:aqua:label:INTRO"

# Tracks by Four Tet that are NOT MP3s
djlt list rb 'artist:Four !kind:MP3'
```

---

## Playlist & Folder Management

The `playlist` and `folder` commands provide full CRUD on the rekordbox playlist tree. Both follow the god-command pattern: the first argument is an `rb/playlists:` or `rb/folders:` query that selects the target(s), and a flag specifies the operation.

### `djlt playlist`

```bash
# Create a new playlist at root level
djlt playlist --new "Fast Bangers"

# Create a new playlist inside a folder
djlt playlist --new "Fast Bangers" --folder "My Sets"

# Create and populate in one step
djlt playlist --new "Fast Bangers" --add "rb/tracks:bpm:128..140"

# Add tracks to one or more existing playlists
djlt playlist rb/playlists:name:Fast --add "rb/tracks:genre:Techno"

# Add to all playlists in a folder simultaneously
djlt playlist "rb/playlists:folder:My Sets" --add "rb/tracks:rating:>=4"

# Rename (requires unambiguous single match)
djlt playlist rb/playlists:name:Fast --rename "Fast Bangers"

# Move matched playlists into a folder
djlt playlist rb/playlists:name:Fast --move "Archive"

# Remove matched playlists
djlt playlist rb/playlists:name:Fast --remove

# Preview any operation without writing
djlt playlist --new "Test" --dry-run
```

| Flag | Description |
| :--- | :--- |
| `--new <name>` | Create a new playlist; combinable with `--add` |
| `--add <rb/tracks query>` | Add matched tracks; use alone to append to existing playlists |
| `--rename <name>` | Rename matched playlists (single match required) |
| `--move <folder>` | Move matched playlists into a folder |
| `--remove` | Remove matched playlists |
| `--folder <name>` | Parent folder for `--new` (default: root level) |
| `--dry-run` | Preview changes without writing |

### `djlt folder`

```bash
# Create a new folder
djlt folder --new "My Sets"

# Create a folder nested inside another
djlt folder --new "Deep Cuts" --parent "My Sets"

# Rename
djlt folder rb/folders:name:Sets --rename "My Sets"

# Move into another folder
djlt folder rb/folders:name:Sets --move "Archive"

# Remove
djlt folder rb/folders:name:Sets --remove
```

| Flag | Description |
| :--- | :--- |
| `--new <name>` | Create a new folder |
| `--rename <name>` | Rename matched folder (single match required) |
| `--move <folder>` | Move matched folder into a parent folder |
| `--remove` | Remove matched folder |
| `--parent <name>` | Parent folder for `--new` (default: root level) |
| `--dry-run` | Preview changes without writing |

> Both commands inherit `--xml` from the global flags and fall back to the configured XML path.

---

## Library Statistics

The `stat` command shows a breakdown of tracks matching a query. When no query is given, it summarises the entire library. It reads the XML path from config if `--xml` is not supplied.

```bash
# Full library summary
djlt stat

# Summary filtered to a query
djlt stat "genre:House && rating:>=4"

# With an explicit XML path
djlt stat --xml ~/Documents/rekordbox.xml "bpm:120..130"
```

Output includes total track count, average BPM, and top genres, artists, and keys.

---

## Synchronization & Export

The `sync` command moves data from Plex into Rekordbox or an M3U8 playlist.

### Plex → Rekordbox
Matches Plex tracks to your Rekordbox collection via fuzzy title/artist matching, optionally transcodes files to 320k MP3 via FFmpeg, and upserts a playlist under a **"Plex Sync"** folder in your XML. Running the same sync a second time updates the existing playlist rather than duplicating it.

```bash
# Dry-run — preview what would change
djlt sync plex:MyPlaylist rb --dry-run

# Full sync with transcoding
djlt sync plex:MyPlaylist rb --dest ./Export --format mp3

# Sync without transcoding (XML injection only)
djlt sync plex:MyPlaylist rb
```

### Plex → M3U8
```bash
djlt sync plex:MyPlaylist m3u8:~/Playlists/MyPlaylist.m3u8
```

### Flags
| Flag | Description |
| :--- | :--- |
| `--dest` | Directory to export transcoded files into |
| `--format` | Target audio format (default: `mp3`) |
| `--dry-run` | Preview all changes without writing any files or XML |

### Transcoding
Transcoding requires FFmpeg to be installed and available on `PATH`. Files that already exist at the destination with a non-zero size are skipped automatically. Artist, album, and title values are sanitized before being used as path components (e.g. `AC/DC` becomes `AC-DC`).

---

## Configuration Reference

```bash
# View all settings
djlt config

# Plex
djlt config plex --host 10.0.0.5
djlt config plex --port 32400
djlt config plex --token <token>
djlt config plex --map "/remote/path:/local/path"
djlt config plex --remove-map "/remote/path"

# Rekordbox
djlt config rekordbox --xml "~/Documents/rekordbox.xml"
```

Settings are persisted to `~/.config/djlt/config.json`.

---

## Development

```bash
# Run all tests
go test ./...

# Build local binary
go build -o djlt ./cmd/djlt
```
