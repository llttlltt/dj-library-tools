# sync

Keep a playlist in sync with a track query

```
djlt sync [source-resource] [source-query] --to [target-resource] [target-query] [flags]
```
### Options

```
      --append          Append new tracks without removing existing ones
      --dest string     Destination directory for exported files
      --format string   Target format for exported files (default "mp3")
  -h, --help            help for sync
      --to strings      Target resource(s) to sync to (repeatable)
```

### Inherited Options

```
      --dry-run      Preview changes without writing
      --json         Output results in JSON format
  -v, --verbose      Enable verbose logging
  -x, --xml string   Path to the Rekordbox XML library
```

Synchronizes a target (like a Rekordbox playlist or M3U file) with a source query.

The sync command is "surgical"—it only adds or removes tracks necessary to make the target
match the source. By default, it removes tracks from the target that no longer match the query.

### Examples

**Keep an "Inbox" playlist matched to specific criteria:**
```bash
djlt sync "rb/tracks added:>today" --to "rb/playlists name:Inbox"

```
**Add new tracks to a playlist without removing existing ones:**
```bash
djlt sync "rb/tracks rating:5" --to "rb/playlists name:Favorites" --append

```
**Sync a query to an external M3U playlist file:**
```bash
djlt sync "rb/tracks genre:House" --to "m3u/path/to/playlist.m3u"



```
## See also

* [djlt](index.md)	 - DJ Library Tools CLI