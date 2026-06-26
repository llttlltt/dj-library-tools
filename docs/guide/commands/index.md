# Command Overview

`djlt` uses an action-centric (verb-centric) command structure. There are six top-level verbs — the **Surgical 6** — each with a short alias for terminal use.

## The Surgical 6

### `list` (Alias: `ls`)
The primary discovery command. Use it to see tracks, playlists, or folders matching a query. Add `--stats` to get aggregate statistics instead of a table.

```bash
# List tracks
djlt ls rb/tracks artist:Four

# Show statistics for a selection
djlt ls rb/tracks "genre:House && bpm:124..128" --stats
```

### `sync`
Mirrors a source selection into a target, ensuring an exact match. Use `--append` to add without removing existing members (replaces the legacy `add` command).

```bash
# Full sync (adds new, removes unmatched)
djlt sync rb/tracks "rating:>=4" --to "rb/playlists name:Favorites"

# Append-only (never removes existing tracks)
djlt sync rb/tracks "genre:House" --to "rb/playlists name:Inbox" --append
```

### `make` (Aliases: `mk`, `create`)
Creates a new playlist or folder. Optionally pre-populate it using `--from`.

```bash
# Create an empty playlist
djlt mk rb/playlists "New Arrivals"

# Create and populate in one step
djlt mk rb/playlists "New Arrivals" --from "rb/tracks added:>2024-01-01"
```

### `move` (Alias: `mv`)
Relocates resources between containers, or renames them in-place using `--name` (replaces the legacy `rename` command).

```bash
# Move a playlist into a folder
djlt mv rb/playlists name:Inbox --to "rb/folders name:Archive"

# Rename a playlist
djlt mv rb/playlists name:Inbox --name "Processed"

# Move tracks between playlists
djlt mv rb/tracks "bpm:>130" --from "rb/playlists name:Inbox" --to "rb/playlists name:'High Energy'"
```

### `remove` (Alias: `rm`)
Handles two distinct operations depending on whether `--from` is present:

- **Resource Deletion** (no `--from`): permanently removes a playlist or folder from the library.
- **Membership Removal** (`--from` present): unlinks tracks from a playlist without deleting them from the collection.

```bash
# Delete a playlist entirely
djlt rm rb/playlists name:Inbox

# Remove specific tracks from a playlist
djlt rm rb/tracks "rating:0" --from "rb/playlists name:Inbox"
```

### `config`
View or update persistent application settings.

```bash
djlt config rekordbox.xml-path "/path/to/export.xml"
djlt config --list
```

---

## Utility Verbs

### `fix`
Corrects common library issues, such as missing file extensions or broken metadata in M3U8 files.

### `update`
Synchronises metadata (like Beatgrids and Tempo markers) between two Rekordbox XML libraries.

---

For full flag details, see the **[CLI Reference](../../commands/index.md)**.
