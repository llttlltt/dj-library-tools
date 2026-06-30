# mk

Create a new playlist or folder

```
djlt mk [resource] [name] [flags]
```
### Options

```
      --at int            Insert at this 0-indexed position (-1 for end) (default -1)
  -h, --help              help for mk
      --in string         Parent folder for the new resource
  -p, --parents           Create parent folders if they don't exist
      --populate string   Source selection to populate the new resource with
```

### Inherited Options

```
      --apply         Actually apply changes to the library (destructive)
  -f, --file string   Path to the primary library file (Rekordbox XML, M3U, etc.)
  -v, --verbose       Enable verbose logging
```

Create a new Rekordbox playlist or folder.
You can optionally populate it immediately using items from a source.

Example:
  djlt mk rb/playlists "New Arrivals" --populate "rb/tracks added:>2024-01-01"
  djlt mk rb/playlists "2024/Jan/Inbox" --parents

## See also

* [djlt](index.md)	 - DJ Library Tools CLI