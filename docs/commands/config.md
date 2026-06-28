# config

View or update application configuration

```
djlt config [key] [value] [flags]
```
### Options

```
  -h, --help    help for config
      --list    Show all configuration values
      --unset   Remove a configuration value
```

### Inherited Options

```
      --dry-run          Preview changes without writing
  -f, --file string      Path to the primary library file (Rekordbox XML, M3U, etc.)
      --json             Output results in JSON format
      --to-file string   Path to the destination library file for sync/move operations
  -v, --verbose          Enable verbose logging
```

Manage djlt configuration using dot-namespaced keys. Settings are stored in ~/.config/djlt/config.json.

## Keys

- **plex.host**: Plex server hostname or IP.
- **plex.port**: Plex server port (default: 32400).
- **plex.token**: Plex authentication token (usually set via 'djlt auth --plex').
- **plex.map**: Remote-to-local path map entry. Used to bridge Plex remote paths to your local mount points.
- **rekordbox.file-path**: Absolute path to your Rekordbox XML export file.

## Examples

**List all settings**
```bash
djlt config --list

```
**Configure Rekordbox library**
```bash
djlt config rekordbox.file-path ~/Documents/rekordbox.xml

```
**Set up Plex connection**
```bash
djlt config plex.host 192.168.1.50
djlt config plex.port 32400

```
**Add a Plex path mapping**
```bash
djlt config plex.map /music/remote:/Volumes/Music

```
**Remove a specific path mapping**
```bash
djlt config --unset plex.map /music/remote

```
**Unset a scalar value**
```bash
djlt config --unset plex.host


```
## See also

* [djlt](index.md)	 - DJ Library Tools CLI