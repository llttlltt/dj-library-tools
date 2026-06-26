# rename

Rename a playlist or folder

```
djlt rename [resource] [query] --to [new-name] [flags]
```
### Options

```
      --dry-run     Preview changes without writing
  -h, --help        help for rename
      --to string   The new name for the resource
```

### Inherited Options

```
  -x, --xml string   Path to the Rekordbox XML library
```

Rename a Rekordbox playlist or folder.
The target must resolve to a single resource.

Example:
  djlt rename rb/playlists Inbox --to "Inbox (Processed)"

## See also

* [djlt](index.md)	 - DJ Library Tools CLI