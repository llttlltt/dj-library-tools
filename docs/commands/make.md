# make

Create a new playlist or folder

```
djlt make [resource] [name] [flags]
```
### Options

```
      --at int        Insert at this 0-indexed position (-1 for end) (default -1)
      --from string   Initial items to populate the resource with
  -h, --help          help for make
      --in string     Parent folder for the new resource
```

### Inherited Options

```
      --dry-run      Preview changes without writing
      --json         Output results in JSON format
  -v, --verbose      Enable verbose logging
  -x, --xml string   Path to the Rekordbox XML library
```

Create a new Rekordbox playlist or folder.
You can optionally populate it immediately using items from a source.

Example:
  djlt mk rb/playlists "New Arrivals" --from "rb/tracks added:>2024-01-01"

## See also

* [djlt](index.md)	 - DJ Library Tools CLI