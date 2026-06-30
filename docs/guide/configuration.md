# Configuration

`djlt` uses a persistent configuration file stored at `~/.config/djlt/config.json`. 

You should rarely need to edit this file manually. Instead, use the built-in `config` command to manage your settings safely.

## Core Settings

### Rekordbox
The primary setting for Rekordbox is the path to your XML export file.
```bash
djlt config rb file "/path/to/your/export.xml"
```

### Plex
Plex configuration requires a host, port, and authentication token.
```bash
djlt config plex host 192.168.1.50
djlt config plex port 32400
```

!!! tip "Authentication"
    While you can set the token manually, it is recommended to use **`djlt config plex auth`** to handle the interactive PIN flow automatically.

### Path Mapping
Path maps allow you to bridge remote file paths (from a NAS or Plex server) to your local mount points.
```bash
djlt config plex map /music/remote:/Volumes/Music
```

## Viewing Settings
You can print your entire configuration at any time:
```bash
djlt config list
```

For more advanced examples, see the full **[djlt config](../commands/config/index.md)** command documentation.
