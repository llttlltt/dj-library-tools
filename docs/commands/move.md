# move

Move items between locations

```
djlt move [resource] [query] --to [destination] [--from origin] [flags]
```
### Options

```
      --dry-run       Preview changes without writing
      --from string   Origin playlist (required for tracks)
  -h, --help          help for move
      --to string     Destination playlist or folder
```

### Inherited Options

```
  -x, --xml string   Path to the Rekordbox XML library
```

Move items between Rekordbox locations.
For tracks, both --from and --to are required.
For playlists and folders, only --to (the parent folder) is required.

Example:
  djlt move rb/tracks "bpm:>130" --from "name:Inbox" --to "name:'High Energy'"
  djlt move rb/playlists "name:'Deep House'" --to "name:Genres"

## See also

* [djlt](index.md)	 - DJ Library Tools CLI