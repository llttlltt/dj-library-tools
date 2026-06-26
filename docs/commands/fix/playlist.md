# fix playlist

Fix playlist extensions and/or enrich with M3U8 metadata

```
djlt fix playlist [file...] [flags]
```

### Options

```
      --dry-run           Show what would be done without modifying files
  -e, --ext strings       Priority list of file extensions (comma-separated)
  -f, --force             Force overwrite if output file exists
  -h, --help              help for playlist
      --m3u8              Enrich playlist with M3U8 #EXTINF tags
  -o, --output string     Specific output path (optional)
  -r, --remove-original   Remove the original playlist file after processing
  -v, --verbose           Enable verbose logging
```

### Inherited Options

```
  -x, --xml string   Path to the Rekordbox XML library
```

## See also

* [djlt fix](index.md)	 - Fix library issues or resource metadata