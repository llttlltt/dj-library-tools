# mv

Move items between locations

```
djlt mv [resource] [query] --to [destination] [--from origin] [flags]
```
### Options

```
      --from string   Origin playlist (required for tracks)
  -h, --help          help for mv
      --name string   New name for the resource (renames)
      --to string     Destination playlist or folder
```

### Inherited Options

```
      --dry-run          Preview changes without writing
  -f, --file string      Path to the primary library file (Rekordbox XML, M3U, etc.)
      --json             Output results in JSON format
      --to-file string   Path to the destination library file for sync/move operations
  -v, --verbose          Enable verbose logging
```

Move items between locations.
For tracks, both --from and --to are required.
For playlists and folders, only --to (the parent folder) is required.

Use the --name flag to rename a resource.

Example:
  djlt mv rb/tracks "bpm:>130" --from "name:Inbox" --to "name:'High Energy'"
  djlt mv rb/playlists name:Inbox --name "Processed"

## See also

* [djlt](index.md)	 - DJ Library Tools CLI