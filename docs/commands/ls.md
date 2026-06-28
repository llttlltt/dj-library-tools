# ls

List items from a location (e.g. rb/tracks title:Oceans)

```
djlt ls [resource] [query] [flags]
```

### Options

```
      --columns strings   Comma-separated list of columns to display
      --exists            Filter for tracks where the physical file exists
  -h, --help              help for ls
      --json              Output results in JSON format
      --missing           Filter for tracks where the physical file is missing
      --sort string       Sort results by any available field (e.g. artist, title, bpm, etc.)
      --stats             Show summary statistics for the selection
```

### Inherited Options

```
      --apply         Actually apply changes to the library (destructive)
  -f, --file string   Path to the primary library file (Rekordbox XML, M3U, etc.)
  -v, --verbose       Enable verbose logging
```

## See also

* [djlt](index.md)	 - DJ Library Tools CLI