# rm

Permanently delete resources or remove membership

```
djlt rm [resource] [query] [flags]
```
### Options

```
      --from strings   Origin resource(s) to remove from
  -h, --help           help for rm
```

### Inherited Options

```
      --dry-run          Preview changes without writing
  -f, --file string      Path to the primary library file (Rekordbox XML, M3U, etc.)
      --json             Output results in JSON format
      --to-file string   Path to the destination library file for sync/move operations
  -v, --verbose          Enable verbose logging
```

Permanently delete resources or remove track membership from playlists.

Use --from to specify which playlist to remove tracks from.
Without --from, the command deletes the resource itself.

Example:
  djlt rm rb/tracks "artist:Four" --from "rb/playlists name:Inbox"
  djlt rm rb/playlists name:Inbox

## See also

* [djlt](index.md)	 - DJ Library Tools CLI