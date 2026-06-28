# ls

List items from a location (e.g. rb/tracks title:Oceans)

```
djlt ls [resource] [query] [flags]
```

### Options

```
  -h, --help          help for ls
      --sort string   Sort results by field (e.g. bpm, artist, title)
      --stats         Show summary statistics for the selection
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

## See also

* [djlt](index.md)	 - DJ Library Tools CLI