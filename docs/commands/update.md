# update

Reconcile metadata between two libraries

```
djlt update [target-selection] --from [source-selection] [flags]
```
### Options

```
  -F, --from string     Source selection to read metadata from
  -h, --help            help for update
      --match strings   Fields to use for matching tracks (default [artist,title])
```

### Inherited Options

```
      --dry-run          Preview changes without writing
  -f, --file string      Path to the primary library file (Rekordbox XML, M3U, etc.)
      --json             Output results in JSON format
      --to-file string   Path to the destination library file for sync/move operations
  -v, --verbose          Enable verbose logging
```

Matches tracks between a source and target selection and synchronizes metadata.

Example:
  djlt update rb/tracks --from plex/tracks --metadata beatgrids

## See also

* [djlt](index.md)	 - DJ Library Tools CLI