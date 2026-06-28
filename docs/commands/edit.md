# edit

Update metadata, repair paths, or fix library issues

```
djlt edit [selection] [query] [flags]
```
### Options

```
  -h, --help              help for edit
      --match strings     Criteria to use for relocation matching (default [filename])
      --relocate string   Search this directory to repair missing file paths
      --repair            Perform provider-specific health/formatting repairs
      --set strings       Metadata fields to update (key:value)
```

### Inherited Options

```
      --dry-run          Preview changes without writing
      --exists           Filter for tracks where the physical file exists
  -f, --file string      Path to the primary library file (Rekordbox XML, M3U, etc.)
      --json             Output results in JSON format
      --missing          Filter for tracks where the physical file is missing
      --to-file string   Path to the destination library file for sync/move operations
  -v, --verbose          Enable verbose logging
```

A unified command for modifying resource state.

Examples:
  # Set a comment for tracks
  djlt edit rb/tracks playlists:Inbox --set comment:Great

  # Relocate missing files
  djlt edit rb/tracks --missing --relocate "/Volumes/Media/Music"

  # Run provider-specific repairs (formerly 'fix')
  djlt edit rb/tracks --repair

## See also

* [djlt](index.md)	 - DJ Library Tools CLI