# mk

Create a new playlist or folder

```
djlt mk [resource] [name] [flags]
```
### Options

```
      --at int        Insert at this 0-indexed position (-1 for end) (default -1)
      --from string   Initial items to populate the resource with
  -h, --help          help for mk
      --in string     Parent folder for the new resource
```

### Inherited Options

```
      --dry-run          Preview changes without writing
  -f, --file string      Path to the primary library file (Rekordbox XML, M3U, etc.)
      --json             Output results in JSON format
      --to-file string   Path to the destination library file for sync/move operations
  -v, --verbose          Enable verbose logging
```

Create a new Rekordbox playlist or folder.
You can optionally populate it immediately using items from a source.

Example:
  djlt mk rb/playlists "New Arrivals" --from "rb/tracks added:>2024-01-01"

## See also

* [djlt](index.md)	 - DJ Library Tools CLI