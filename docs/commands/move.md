# move

Move items between locations

```
djlt move [resource] [query] --to [destination] [--from origin] [flags]
```
### Options

```
      --from string   Origin playlist (required for tracks)
  -h, --help          help for move
      --name string   New name for the resource (renames)
      --to string     Destination playlist or folder
```

### Inherited Options

```
      --dry-run      Preview changes without writing
      --json         Output results in JSON format
  -v, --verbose      Enable verbose logging
  -x, --xml string   Path to the Rekordbox XML library
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