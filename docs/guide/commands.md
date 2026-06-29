# Command Overview

`djlt` uses an action-centric (verb-centric) command structure.

## The Core Commands

### `ls`
The primary discovery command. Use it to see tracks, playlists, or folders matching a query. Add `--stats` to get aggregate statistics instead of a table.

```bash
# List tracks
djlt ls rb/tracks artist:Four

# Show statistics for a selection
djlt ls rb/tracks "genre:House && bpm:124..128" --stats
```

### `sync`
Mirrors a source selection into a target, ensuring an exact match. Use `--append` to add without removing existing members.

```bash
# Full sync (adds new, removes unmatched)
djlt sync rb/tracks "rating:>=4" --to "rb/playlists name:Favorites"

# Append-only (never removes existing tracks)
djlt sync rb/tracks "genre:House" --to "rb/playlists name:Inbox" --append
```

### `edit`
A unified command for modifying resource state. It replaces the legacy `modify` and `fix` commands.

```bash
# Set metadata in bulk
djlt edit rb/tracks "rating:0" --set "rating:3"

# Repair missing file paths
djlt edit rb/tracks --missing --relocate "/Volumes/Music"

# Run provider-specific repairs
djlt edit rb/playlists --repair
```

### `mk`
Creates a new playlist or folder. Optionally pre-populate it using `--from`.

```bash
# Create an empty playlist
djlt mk rb/playlists "New Arrivals"

# Create a folder hierarchy
djlt mk rb/folders "2024/Jan/Sorting" --parents
```

### `mv`
Relocates resources between containers, or renames them in-place using `--name`.

```bash
# Move a playlist into a folder
djlt mv rb/playlists name:Inbox --to "rb/folders name:Archive"

# Move tracks between playlists
djlt mv rb/tracks "bpm:>130" --from "rb/playlists name:Inbox" --to "rb/playlists name:'High Energy'"
```

### `rm`
Removes resources or track membership.

```bash
# Delete a playlist entirely
djlt rm rb/playlists name:Inbox

# Delete a folder and everything inside
djlt rm rb/folders name:OldSets --recursive

# Remove specific tracks from a playlist
djlt rm rb/tracks "rating:0" --from "rb/playlists name:Inbox"
```

### `config`
View or update persistent application settings and provider authentication.

```bash
djlt config rb file "/path/to/export.xml"
djlt config plex auth
djlt config list
```

---

For full flag details, see the **[CLI Reference](../../commands/index.md)**.
