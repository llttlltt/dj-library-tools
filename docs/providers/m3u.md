# M3U / M3U8

The M3U provider allows you to read from and write to standard music playlist files (`.m3u` and `.m3u8`). It treats individual files as libraries containing track collections.

## Providers

| Alias | Description |
| :--- | :--- |
| `m3u` | Standard M3U library provider. |
| `m3u8` | UTF-8 encoded M3U library provider (preferred for paths with special characters). |

## Resources

| Resource | Description | Example |
| :--- | :--- | :--- |
| `tracks` | The tracks contained within the file. | `m3u8/tracks --file list.m3u8` |
| `playlists` | The M3U file itself as a ResourceGroup. | `m3u8/playlists --file list.m3u8` |

## Fields

### Track Fields

The M3U provider parses metadata from native `#EXTINF` tags.

| Field | Type | Description |
| :--- | :--- | :--- |
| `display` | String | The raw display name string from the `#EXTINF` tag. |
| `location` | String | The physical path to the audio file. |
| `duration` | Numeric | Track duration in seconds. |

### Playlist Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `name` | String | Filename of the playlist. |
| `items` | Numeric | Number of tracks in the file. |

## Technical Details

- **File-Based**: M3U providers are strictly file-based and require the `--file` (`-f`) flag for initialization.
- **Relative Paths**: When an M3U file is loaded, any relative paths are automatically resolved against the directory containing the `.m3u` file.
- **UTF-8 Support**: Using the `m3u8` alias ensures that UTF-8 encoding is handled correctly for non-ASCII characters.
- **Save Operation**: Modifications (adding/removing tracks) are not written to disk until the provider's `Save` method is called. This is handled automatically by commands like `sync`.

## Examples

### Reading

**List all tracks in a playlist**
```bash
djlt ls m3u8/tracks --file ~/Music/Favorites.m3u8
```

**Filter tracks by part of the filename**
```bash
djlt ls m3u8/tracks --file ./Techno.m3u8 "location:Techno"
```

### Management

**Sync from Rekordbox to M3U8**
```bash
djlt sync "rb/tracks genre:Techno" --to "m3u8/tracks" --to-file "./Techno.m3u8"
```

**Remove specific tracks from an M3U8 file by display name**
```bash
djlt rm m3u8/tracks --file ./list.m3u8 "display:Unknown"
```
