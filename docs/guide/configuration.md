# Configuration

`djlt` uses a persistent configuration file stored at `~/.config/djlt/config.json`. 

You should rarely need to edit this file manually. Instead, use the built-in `config` command to manage your settings safely.

## Quick Start

### 1. Set your Rekordbox XML path
This is the most important setting. It tells `djlt` where to find your library.
```bash
djlt config rekordbox.xml-path "/path/to/your/export.xml"
```

### 2. Verify your settings
You can print your entire configuration at any time:
```bash
djlt config --list
```

## Detailed Reference

For a complete list of all available keys (Plex connection, Path mapping, etc.) and advanced examples, see the full **[djlt config](../commands/config.md)** command documentation.

!!! tip "Authentication"
    While you can set `plex.token` via the config command, it is recommended to use **`djlt auth --plex`** to handle the OAuth flow automatically.
