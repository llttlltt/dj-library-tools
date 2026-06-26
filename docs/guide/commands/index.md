# Command Overview

`djlt` uses an action-centric (verb-centric) command structure. Commands are organized by the action you want to perform on your library.

## Core Verbs

### `list` (Alias: `ls`)
The primary discovery command. Use it to see tracks, playlists, or folders matching a query.
```bash
djlt list rb/tracks artist:Four
```

### `add`
Links items from a source selection to one or more targets.
```bash
djlt add rb/tracks genre:House --to "rb/playlists name:Summer"
```

### `remove` (Alias: `rm`)
Unlinks items matching a source selection from one or more origins.
```bash
djlt remove rb/tracks rating:0 --from "rb/playlists name:Inbox"
```

### `sync`
Mirrors a source selection into a target, adding or removing items to ensure an exact match.
```bash
djlt sync rb/tracks rating:5 --to "rb/playlists name:Favorites"
```

## Management Verbs

### `create`
Initializes a new resource (Playlist or Folder) and optionally populates it.
```bash
djlt create rb/playlists "New Arrivals" --from rb/tracks added:today
```

### `rename`
Changes the name of a playlist or folder.
```bash
djlt rename rb/playlists Inbox --to "Inbox (Archived)"
```

### `move` (Alias: `mv`)
Relocates resources (e.g. moving a playlist into a folder) or shifts tracks between playlists.

### `delete` (Aliases: `del`, `rm`)
Permanently removes a resource from the library.

## Utility Verbs

### `stat` (Alias: `stats`)
Provides statistical analysis (BPM, Keys, Genres) for a track selection.

### `fix`
Corrects common library issues, such as missing file extensions or broken metadata in M3U8 files.

### `update`
Synchronizes metadata (like Beatgrids/Tempo markers) between different Rekordbox XML libraries.

---

For a full list of flags and options, see the **[CLI Reference](../../commands/index.md)**.
