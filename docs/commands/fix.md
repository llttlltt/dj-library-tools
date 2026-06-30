# fix

Perform health, formatting, and structural repairs on the library

```
djlt fix [selection] [query] [flags]
```
### Options

```
      --duplicates strings   Remove duplicates (targets: members, tracks)
  -h, --help                 help for fix
      --metadata strings     Fix/normalize metadata (targets: artist, album, etc.)
      --orphans strings      Remove orphaned resources (targets: all)
      --paths strings        Repair file paths (targets: relocate, normalize)
```

### Inherited Options

```
      --apply         Actually apply changes to the library (destructive)
  -f, --file string   Path to the primary library file (Rekordbox XML, M3U, etc.)
  -v, --verbose       Enable verbose logging
```

A multi-purpose repair command for library maintenance.

Examples:
  # Remove duplicate tracks from specific playlists
  djlt fix rb/playlists "Inbox,Recently Added" --duplicates members

  # Normalize metadata for matching tracks
  djlt fix rb/tracks "genre:Techno" --metadata artist,album

  # Repair broken file paths for missing tracks
  djlt fix rb/tracks "missing:true" --paths normalize

## See also

* [djlt](index.md)	 - DJ Library Tools CLI