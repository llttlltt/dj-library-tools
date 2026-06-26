# M3U / M3U8

The M3U provider allows you to read from and write to standard music playlist files (`.m3u` and `.m3u8`). It treats individual files as playlist resources that can be queried, modified, or synced with other libraries.

## Providers

| Alias | Description |
| :--- | :--- |
| `m3u` | Standard M3U library provider. |
| `m3u8` | UTF-8 encoded M3U library provider (preferred for paths with special characters). |

## Resources

| Resource | Description | Example |
| :--- | :--- | :--- |
| `tracks` | The tracks contained within the file. | `m3u8:/path/to/list.m3u8/tracks` |
| `playlists` | The M3U file itself as a node (Default). | `m3u8:/path/to/list.m3u8` |

!!! tip "Path Syntax"
    Unlike other providers, the "Resource" part for M3U is the **file path**. 
    You can omit `/tracks` to list the playlist itself, or append it to filter the tracks inside.

## Fields

### Track Fields

The M3U provider parses metadata from `#EXTINF` tags.

| Field | Type | Description |
| :--- | :--- | :--- |
| `title` | String | Track title. Extracted from `#EXTINF` or filename. |
| `artist` | String | Artist name. Extracted from `#EXTINF`. |
| `location` | String | Full path to the audio file. |

### Playlist Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `name` | String | Filename of the playlist. |
| `items` | Numeric | Number of tracks in the file. |

## Technical Details

- **Relative Paths**: When an M3U file is loaded, any relative paths are automatically resolved against the directory containing the `.m3u` file.
- **UTF-8 Support**: Using the `m3u8` prefix ensures that UTF-8 encoding is handled correctly for non-ASCII characters in file paths or metadata.
- **Save Operation**: Modifications (adding/removing tracks) are not written to disk until the provider's `Save` method is called. This is handled automatically by commands like `mk`, `rm`, and `sync`.

## Examples

### Reading

**List tracks in a playlist**
```bash
djlt ls m3u8:~/Music/Favorites.m3u8
```

**Filter tracks within a file**
```bash
djlt ls m3u8:./Techno.m3u8 "artist:Derrick"
```

### Management

**Create a new playlist from a Rekordbox query**
```bash
djlt mk m3u8:./Techno.m3u8 "Techno" --from "rb/tracks genre:Techno rating:5"
```

**Sync from Plex to M3U8**
```bash
djlt sync plex/playlists "name:Road Trip" --to m3u8:~/Desktop/RoadTrip.m3u8
```

**Remove specific tracks from an M3U8 file**
```bash
djlt rm m3u8:./list.m3u8/tracks "artist:Unknown" --from m3u8:./list.m3u8
```
