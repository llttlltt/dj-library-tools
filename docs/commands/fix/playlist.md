# fix playlist

Fix playlist extensions and/or enrich with M3U8 metadata

```
djlt fix playlist [file...] [flags]
```

### Options

```
  -e, --ext strings       Priority list of file extensions (comma-separated)
      --force             Force overwrite if output file exists
  -h, --help              help for playlist
      --m3u8              Enrich playlist with M3U8 #EXTINF tags
  -o, --output string     Specific output path (optional)
  -r, --remove-original   Remove the original playlist file after processing
```

### Inherited Options

```
      --dry-run          Preview changes without writing
  -f, --file string      Path to the primary library file (Rekordbox XML, M3U, etc.)
      --json             Output results in JSON format
      --to-file string   Path to the destination library file for sync/move operations
  -v, --verbose          Enable verbose logging
```

## See also

* [djlt fix](index.md)	 - Fix library issues or resource metadata