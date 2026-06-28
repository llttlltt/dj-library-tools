# update

Update track metadata or merge markers between libraries

```
djlt update [resource] [query] --from [source-file] [flags]
```
### Options

```
      --force           Overwrite output file if it already exists
  -F, --from string     Source library file to read metadata from
  -h, --help            help for update
      --merge           Merge metadata instead of overwriting
  -o, --output string   Output path for the updated Rekordbox XML
  -t, --to string       Destination Rekordbox XML to update (defaults to primary library)
```

### Inherited Options

```
      --dry-run          Preview changes without writing
  -f, --file string      Path to the primary library file (Rekordbox XML, M3U, etc.)
      --json             Output results in JSON format
      --to-file string   Path to the destination library file for sync/move operations
  -v, --verbose          Enable verbose logging
```

Update metadata for tracks in the library using another Rekordbox XML as a source.
Currently supports updating/merging Tempo markers (Beatgrids).

Example:
  djlt update rb/tracks --from other_library.xml

## See also

* [djlt](index.md)	 - DJ Library Tools CLI