# update

Update track metadata or merge markers between libraries

```
djlt update [resource] [query] --from [source-xml] [flags]
```
### Options

```
      --force           Overwrite output file if it already exists
  -f, --from string     Source Rekordbox XML to read metadata from
  -h, --help            help for update
      --merge           Merge metadata instead of overwriting
  -o, --output string   Output path for the updated Rekordbox XML
  -t, --to string       Destination Rekordbox XML to update (defaults to primary library)
```

### Inherited Options

```
      --dry-run      Preview changes without writing
      --json         Output results in JSON format
  -v, --verbose      Enable verbose logging
  -x, --xml string   Path to the Rekordbox XML library
```

Update metadata for tracks in the library using another Rekordbox XML as a source.
Currently supports updating/merging Tempo markers (Beatgrids).

Example:
  djlt update rb/tracks --from other_library.xml

## See also

* [djlt](index.md)	 - DJ Library Tools CLI