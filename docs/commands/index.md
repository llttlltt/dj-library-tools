# djlt

DJ Library Tools CLI

### Options

```
      --dry-run          Preview changes without writing
      --exists           Filter for tracks where the physical file exists
  -f, --file string      Path to the primary library file (Rekordbox XML, M3U, etc.)
  -h, --help             help for djlt
      --json             Output results in JSON format
      --missing          Filter for tracks where the physical file is missing
      --to-file string   Path to the destination library file for sync/move operations
  -v, --verbose          Enable verbose logging
```

A comprehensive CLI tool for managing DJ libraries across multiple providers.
## See also

* [djlt config](config.md)	 - Manage application and provider settings
* [djlt edit](edit.md)	 - Update metadata, repair paths, or fix library issues
* [djlt ls](ls.md)	 - List items from a location (e.g. rb/tracks title:Oceans)
* [djlt mk](mk.md)	 - Create a new playlist or folder
* [djlt mv](mv.md)	 - Move items between locations
* [djlt rm](rm.md)	 - Permanently delete resources or remove membership
* [djlt sync](sync.md)	 - Keep a playlist or metadata in sync with a track query