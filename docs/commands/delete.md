# delete

Delete a resource from the library (destructive)

```
djlt delete [resource] [query] [flags]
```
### Options

```
      --dry-run   Preview changes without writing
  -h, --help      help for delete
```

### Inherited Options

```
  -x, --xml string   Path to the Rekordbox XML library
```

Permanently delete playlists or folders from the Rekordbox XML.
Warning: This is destructive to the resource, but does not delete tracks from the collection.

Example:
  djlt delete rb/playlists "Old Mixes"

## See also

* [djlt](index.md)	 - DJ Library Tools CLI