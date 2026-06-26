# Syncing to Mobile (M3U8)

!!! info "Roadmap Feature"
    M3U8 export is currently on the development roadmap and is scheduled for an upcoming release. The examples below show the planned syntax.

If you use a mobile player or a third-party app that doesn't support Plex or Rekordbox directly, you will be able to use `djlt` to export your playlists to a standard M3U8 format.

### 1. Export a Plex Playlist (Planned)
You will be able to take any selection from Plex and generate an M3U8 file. This is useful for creating offline backups or for use in mobile music players.

```bash
# Sync a Plex playlist to a local M3U8 file
djlt sync plex/playlists name:Summer --to "m3u8:/path/to/playlist.m3u8"
```

!!! note "Path Mapping"
    Ensure your `plex.map` is correctly configured in `djlt config` so that the file paths in the M3U8 point to your local files, not the remote Plex server paths.

### 2. Export from Rekordbox (Planned)
You will also be able to generate M3U8 files from your Rekordbox playlists or queries.

```bash
# Export all 5-star tracks to an M3U8 file
djlt sync rb/tracks "rating:5" --to "m3u8:~/Desktop/Favorites.m3u8"
```

---

## Roadmap: Rekordbox to Plex
Currently, `djlt` supports **Plex → Rekordbox**. 

**Rekordbox → Plex** syncing (uploading playlists back to your Plex server) and **M3U8 Exports** are highly requested features and are currently in development. Because the Plex API for playlist creation is restrictive, this requires a specialized write-provider which is scheduled for a future release.
