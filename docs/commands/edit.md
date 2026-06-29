# edit

Update metadata for resources

```
djlt edit [selection] [query] [flags]
```
### Options

```
      --exists        Filter for tracks where the physical file exists
  -h, --help          help for edit
      --missing       Filter for tracks where the physical file is missing
      --set strings   Metadata fields to update (key:value)
```

### Inherited Options

```
      --apply         Actually apply changes to the library (destructive)
  -f, --file string   Path to the primary library file (Rekordbox XML, M3U, etc.)
  -v, --verbose       Enable verbose logging
```

Modify metadata fields for tracks or other resources.
For library maintenance (deduplication, path repair), use 'djlt fix'.

Examples:
  # Set a comment for tracks
  djlt edit rb/tracks playlists:Inbox --set comment:Great

## See also

* [djlt](index.md)	 - DJ Library Tools CLI