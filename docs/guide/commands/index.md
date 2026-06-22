# Commands Reference

`djlt` commands are organized by the resource they manage.

## `auth`
Manages authentication with external services.
- `--plex`: Triggers PIN-based OAuth flow for Plex.

## `config`
Manages local configuration settings.
- `--list`: Prints all settings.
- `--unset <key>`: Clears a setting.

## `list`
Discovers and displays resources.
- `djlt list provider/resource:query`

## `playlist`
Manages Rekordbox playlists.
- `--new <name>`: Create a new playlist.
- `--add <query>`: Append matching tracks to a playlist.
- `--sync <query>`: Synchronize playlist content (add/remove to match query).
- `--delete`: Delete a playlist.

## `folder`
Manages Rekordbox folders.
- `--new <name>`: Create a new folder.
- `--move <name>`: Move a folder.

## `sync`
Data movement and transcoding.
- `djlt sync [source] [destination]`
- `--dry-run`: Show what would be moved without performing operations.
