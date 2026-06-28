# Exporting to Mobile Phones

If you use a mobile player or a third-party app, you can use `djlt` to export your playlists into a standard file format (`.m3u8`) that phones can understand.

### 1. Export from Rekordbox or Plex
You can take any search or playlist and generate a mobile-ready file.

```bash
# Create a playlist file for your phone from all 5-star Rekordbox tracks
djlt sync "rb/tracks rating:5" --to "m3u8/tracks" --to-file "favorites.m3u8"
```

### 2. Keep your mobile files updated
When you add new music to your main collection, just run the command again. `djlt` will update the file with your latest tracks.

```bash
# Sync your phone's Techno list with your latest library changes
djlt sync "rb/tracks genre:Techno" --to "m3u8/tracks" --to-file "techno.m3u8"
```

### 3. Cleanup your mobile playlists
If you want to remove specific tracks from your phone's playlist file:

```bash
# Remove 1-star tracks from your mobile playlist file
djlt rm m3u8/tracks --file mobile.m3u8 "rating:1"
```

---

## Tip: Path Mapping
If your phone expects files to be in a different folder than your computer (e.g., `/sdcard/Music` vs `/Users/You/Music`), ensure you have your **Path Maps** configured in `djlt config`. This ensures the playlist file uses the correct paths for your device.
