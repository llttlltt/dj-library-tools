# sync

Keep a playlist or metadata in sync with a track query

```
djlt sync [source-resource] [source-query] --to [target-resource] [target-query] [flags]
```
### Options

```
      --append             Append new tracks without removing existing ones
      --dest string        Destination directory for exported files
      --format string      Target format for exported files (default "mp3")
  -h, --help               help for sync
      --match strings      Fields to use for matching tracks (default [artist,title])
      --metadata strings   Metadata fields to synchronize (e.g. beatgrids, rating)
      --to strings         Target resource(s) to sync to (repeatable)
```

### Inherited Options

```
      --apply          Actually apply changes to the library (destructive)
      --exists           Filter for tracks where the physical file exists
  -f, --file string      Path to the primary library file (Rekordbox XML, M3U, etc.)
      --json             Output results in JSON format
      --missing          Filter for tracks where the physical file is missing
      --to-file string   Path to the destination library file for sync/move operations
  -v, --verbose          Enable verbose logging
```

Synchronizes a target (like a Rekordbox playlist or M3U file) with a source query.

The sync command is "surgical"—it only adds or removes tracks necessary to make the target
match the source. By default, it removes tracks from the target that no longer match the query.

### Metadata Reconciliation
If --metadata is specified, djlt will match tracks between the source and target using the --match keys
and synchronize specific metadata fields (e.g. beatgrids, rating).

### Examples

**Keep an "Inbox" playlist matched to specific criteria:**
```bash
djlt sync "rb/tracks added:>today" --to "rb/playlists name:Inbox"

```
**Sync beatgrids from a backup Rekordbox XML to your primary library:**
```bash
djlt sync "rb/tracks" --file backup.xml --to "rb/tracks" --metadata beatgrids

```
**Sync ratings from Plex to Rekordbox matching by filename:**
```bash
djlt sync "plex/tracks" --to "rb/tracks" --metadata rating --match filename



```
## See also

* [djlt](index.md)	 - DJ Library Tools CLI