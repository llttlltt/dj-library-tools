# DJ Library Tools (djlt)

A unified suite of high-performance tools for managing DJ libraries, Rekordbox XML exports, and cross-platform synchronization.

## 🛠 Installation

```bash
go install github.com/llttlltt/dj-library-tools/cmd/djlt@latest
```

## 🚀 Getting Started

Initialize your configuration:

```bash
# Authenticate with Plex
djlt auth plex

# Map remote Plex paths to local mount points (if syncing from Plex)
djlt config plex --map "/media/Music:/Volumes/Music"

# Set your master Rekordbox XML path
djlt config rekordbox --xml "~/Documents/rekordbox.xml"
```

---

## 🔍 Discovery & Querying

`djlt` features a powerful selection engine for filtering your library. The syntax follows a `source query` or `provider/resource:query` pattern.

### Command Format
```bash
# Recommendation: Use single quotes '' for queries to prevent shell interference with ! or ( )
djlt list [source] '[query]'
```

### Sources
| Source | Alias | Description |
| :--- | :--- | :--- |
| `rb/tracks` | `rb` | All tracks in your Rekordbox Collection |
| `rb/playlists` | | All playlists in Rekordbox |
| `plex/playlists` | `plex` | All Plex playlists |
| `plex/tracks` | | Fetch tracks from a specific Plex playlist ID |

### Boolean Operators
You can combine multiple criteria using logic:
- `&&` or `AND` : Both conditions must be met.
- `||` or `OR` : Either condition must be met.
- `!` : Negation (must NOT match).
- `( )` : Parentheses for grouping logic.

### Comparison Operators
| Operator | Type | Example |
| :--- | :--- | :--- |
| `:` | Substring | `artist:Four` |
| `=` | Exact | `artist="Four Tet"` |
| `::` | Regex | `name::"^01"` |
| `..` | Range | `bpm:124..128` |
| `>`, `>=` | Numeric | `bpm:>=128` |
| `<`, `<=` | Numeric | `rating:<3` |

### Query Properties
You can filter by any standard Rekordbox field:

#### Metadata
- `name`, `artist`, `album`, `genre`, `key`, `year`, `label`, `comment`, `remixer`, `mix`

#### Technical
- `bpm`, `bitrate`, `kind` (e.g., MP3 File), `size`, `time` (duration in seconds)

#### Library State
- `rating` (0-5 stars), `playcount`, `added` (date), `playlistcount` (number of playlists)
- `playlist` (match by name, e.g. `playlist:Summer`)

#### Cues & Beatgrids (Advanced)
| Property | Description |
| :--- | :--- |
| `beatgrids` | Number of TEMPO markers (e.g. `beatgrids:>1`) |
| `hotcues` | Number of hot cues (e.g. `hotcues:8`) or check color `hotcues:aqua` |
| `memorycues` | Number of memory cues or check color `memorycues:pink` |
| `hotcue:[a-h]` | Target a specific slot (e.g. `hotcue:a:green`) |
| `memorycue:[idx]`| Target by position high-to-low (e.g. `memorycue:1:loop`) |

**Cue Sub-properties:**
- `label:[text|empty]` : Match by text or check if unlabeled (`label:empty`).
- `loop` : Match if the cue is an active loop.
- `[color]` : `pink`, `red`, `orange`, `yellow`, `green`, `aqua`, `blue`, `purple`, `none`.

### Examples
```bash
# High rated House tracks between 124-128 BPM
djlt list rb "genre:House && rating:>=4 && bpm:124..128"

# Tracks with multiple beatgrid markers (transition tracks)
djlt list rb "beatgrids:>1"

# Tracks where Hot Cue B is Aqua and labeled "INTRO"
djlt list rb "hotcue:b:aqua:label:INTRO"

# Tracks by Four Tet that are NOT MP3s
djlt list rb 'artist:Four !kind:MP3'

# Tracks appearing on more than 3 playlists
djlt list rb "playlistcount:>3"

# Tracks appearing in both specific playlists
djlt list rb "playlist:Summer && playlist:Beach"
```

---

## 🔄 Synchronization & Export

The `sync` command orchestrates data and media movement.

### Sync Plex to Rekordbox
Matches Plex tracks to Rekordbox, transcodes to 320k MP3 via FFmpeg, and injects a "Plex Sync" folder into your XML.

```bash
djlt sync plex:MyPlaylist rb --dest ./Export --format mp3
```

---

## 🛠 Development

```bash
# Run tests
go test ./...

# Build local binary
go build -o djlt ./cmd/djlt
```
