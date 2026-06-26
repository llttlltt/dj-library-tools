# remove

Remove a resource or membership from the library

```
djlt remove [resource] [query] [flags]
```
### Options

```
      --from strings   Origin resource(s) to remove from (repeatable)
  -h, --help           help for remove
```

### Inherited Options

```
      --dry-run      Preview changes without writing
      --json         Output results in JSON format
  -v, --verbose      Enable verbose logging
  -x, --xml string   Path to the Rekordbox XML library
```

Permanently delete resources or remove track membership from playlists.

Use --from to specify which playlist to remove tracks from.
Without --from, the command deletes the resource itself.

Example:
  djlt rm rb/tracks artist:Four --from "rb/playlists name:Inbox"
  djlt rm rb/playlists name:Inbox

## See also

* [djlt](index.md)	 - DJ Library Tools CLI