# djlt config plex

Configure Plex Media Server settings.

## Subcommands

### auth
Interactive Plex authentication (PIN flow).
```bash
djlt config plex auth
```

### host
Set or get the Plex server hostname or IP.
```bash
djlt config plex host [value]
```

### port
Set or get the Plex server port (default 32400).
```bash
djlt config plex port [value]
```

### map
Add a path mapping to bridge remote Plex paths to local mount points.
```bash
djlt config plex map [remote:local]
```

## Options

```
  -h, --help   help for plex
```

## Inherited Options

```
      --apply         Actually apply changes to the library (destructive)
  -f, --file string   Path to the primary library file (Rekordbox XML, M3U, etc.)
  -v, --verbose       Enable verbose logging
```
