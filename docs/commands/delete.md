# delete

Delete a resource from the library (destructive)

```
djlt delete [resource] [query] [flags]
```
### Options

```
  -h, --help   help for delete
```

### Inherited Options

```
      --dry-run      Preview changes without writing
  -v, --verbose      Enable verbose logging
  -x, --xml string   Path to the Rekordbox XML library
```

Permanently delete playlists or folders from the Rekordbox XML.
Warning: This is destructive to the resource, but does not delete tracks from the collection.

Example:
  djlt delete rb/playlists "name:'Old Mixes'"

## See also

* [djlt](index.md)	 - DJ Library Tools CLI