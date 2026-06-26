# create

Create a new playlist or folder

```
djlt create [resource] [name] [flags]
```
### Options

```
      --at int        Insert at this 0-indexed position (-1 for end) (default -1)
      --dry-run       Preview changes without writing
      --from string   Initial items to populate the resource with
  -h, --help          help for create
      --in string     Parent folder for the new resource
```

### Inherited Options

```
  -x, --xml string   Path to the Rekordbox XML library
```

Create a new Rekordbox playlist or folder.
You can optionally populate it immediately using items from a source.

Example:
  djlt create rb/playlists "New Arrivals" --from "rb/tracks added:>2024-01-01"

## See also

* [djlt](index.md)	 - DJ Library Tools CLI