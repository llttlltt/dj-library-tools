# add

Add items from a source to one or more targets

```
djlt add [source-resource] [source-query] --to [target-resource] [target-query] [flags]
```
### Options

```
      --force        Allow adding duplicates (if supported by target)
  -h, --help         help for add
      --to strings   Target resource(s) to add to (repeatable)
```

### Inherited Options

```
      --dry-run      Preview changes without writing
  -x, --xml string   Path to the Rekordbox XML library
```

Add items from a source selection to one or more target resources.
Currently supports adding tracks (rb/tracks) to playlists (rb/playlists).

Example:
  djlt add rb/tracks artist:Four --to "rb/playlists name:Inbox"

## See also

* [djlt](index.md)	 - DJ Library Tools CLI