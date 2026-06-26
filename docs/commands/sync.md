# sync

Sync items between a source and one or more targets

```
djlt sync [source-resource] [source-query] --to [target-resource] [target-query] [flags]
```

### Options

```
      --dest string     Destination directory for exported files
      --dry-run         Preview changes without writing files or XML
      --format string   Target format for exported files (default "mp3")
  -h, --help            help for sync
      --to strings      Target resource(s) to sync to (repeatable)
```

### Inherited Options

```
  -v, --verbose      Enable verbose logging
  -x, --xml string   Path to the Rekordbox XML library
```

## See also

* [djlt](index.md)	 - DJ Library Tools CLI