# Syncing to Mobile (M3U8)

If you use a mobile player or a third-party app that doesn't support Plex or Rekordbox directly, you can use `djlt` to export your playlists to a standard M3U8 format.

### 1. Export a Plex Selection
You can take any selection from Plex and generate an M3U8 file. This is useful for creating offline backups or for use in mobile music players.

```bash
# Sync a Plex playlist to a local M3U8 file
djlt sync plex/playlists name:Summer --to "m3u8:/path/to/playlist.m3u8"
```

!!! note "Path Mapping"
    Ensure your `plex.map` is correctly configured in `djlt config` so that the file paths in the M3U8 point to your local files, not the remote Plex server paths.

### 2. Export from Rekordbox
You can also generate M3U8 files from your Rekordbox playlists or queries.

```bash
# Export all 5-star tracks to an M3U8 file
djlt sync rb/tracks "rating:5" --to "m3u8:~/Desktop/Favorites.m3u8"
```

### 3. Creating and Managing M3U8 Files
You can create M3U8 files directly using the `mk` command:

```bash
# Create a new M3U8 file from a Rekordbox folder
djlt mk m3u8/house.m3u8 "House Favorites" --from "rb/folders name:House"
```

You can even remove items from an existing M3U8 file:

```bash
# Remove 1-star tracks from a mobile playlist
djlt rm m3u8/mobile.m3u8/tracks "rating:1" --from m3u8/mobile.m3u8
```

---

## Roadmap: Rekordbox to Plex
Currently, `djlt` supports **Plex → Rekordbox** and **Any → M3U8**.

**Rekordbox → Plex** syncing (uploading playlists back to your Plex server) is a highly requested feature and is currently in development. Because the Plex API for playlist creation is restrictive, this requires a specialized write-provider which is scheduled for a future release.
