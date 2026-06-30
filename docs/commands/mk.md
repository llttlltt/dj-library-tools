# mk

Create a new playlist or folder

```
djlt mk [resource] [name] [flags]
```
### Options

```
      --at string         Insert position: a positive integer (1-based), "start", or "end" (default: end)
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

The --at flag controls insertion position using 1-based indexing or named sentinels:
  --at start   Insert at the first position
  --at end     Append to the end (default when flag is omitted)
  --at 2       Insert at the second position (1-based)
Omitting --at is equivalent to --at end.

Examples:
  djlt mk rb/playlists "New Arrivals" --populate "rb/tracks added:>2024-01-01"
  djlt mk rb/playlists "2024/Jan/Inbox" --parents
  djlt mk rb/playlists "Inbox" --in "Sorting" --at start
  djlt mk rb/playlists "Archive" --in "Sorting" --at end
  djlt mk rb/playlists "Featured" --in "Sorting" --at 2

## See also

* [djlt](index.md)	 - DJ Library Tools CLI