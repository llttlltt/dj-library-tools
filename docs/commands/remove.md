# remove

Remove items from one or more origins

```
djlt remove [source-resource] [source-query] --from [origin-resource] [origin-query] [flags]
```
### Options

```
      --from strings   Origin resource(s) to remove from (repeatable)
  -h, --help           help for remove
```

### Inherited Options

```
      --dry-run      Preview changes without writing
  -x, --xml string   Path to the Rekordbox XML library
```

Remove items matching a source selection from one or more origin resources.
Currently supports removing tracks (rb/tracks) from playlists (rb/playlists).

Example:
  djlt remove rb/tracks artist:Four --from "rb/playlists name:Inbox"

## See also

* [djlt](index.md)	 - DJ Library Tools CLI