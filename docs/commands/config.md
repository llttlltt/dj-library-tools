# djlt config

Manage application settings for Plex, Rekordbox, and path mappings.

## Usage

`djlt config [subcommand] [flags]`

## Subcommands

### list
Show all currently saved configuration values.
```bash
djlt config list
```

### plex
Configure Plex Media Server settings.
- **auth**: Trigger the interactive PIN flow to authenticate with Plex.
- **host [value]**: Set the server address.
- **port [value]**: Set the server port (default 32400).
- **map [remote:local]**: Add a path mapping to bridge remote Plex paths to local mount points.

```bash
djlt config plex host 192.168.1.50
djlt config plex auth
djlt config plex map /music:/Volumes/Music
```

### rb
Configure Pioneer Rekordbox settings.
- **file [path]**: Set the path to your primary Rekordbox XML export.

```bash
djlt config rb file ~/Documents/rekordbox.xml
```
