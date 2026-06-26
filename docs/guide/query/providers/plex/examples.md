# Plex Examples

**List a Plex playlist by name**
```bash
djlt list plex/playlists name:Summer
```

**Sync a Plex playlist to Rekordbox**
```bash
djlt sync plex/playlists name:Summer --to "rb/playlists name:'Plex Sync'"
```

**Sync a Plex playlist to M3U8**
```bash
djlt sync plex/playlists name:Summer --to "m3u8:/path/to/playlist.m3u8"
```
