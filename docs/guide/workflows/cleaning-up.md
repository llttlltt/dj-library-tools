# Cleaning Up Your Library

Over time, any DJ library accumulates "dead weight"—duplicate tracks, low-quality files, or music you simply never play. This workflow shows how to use `ls` and `rm` to identify and prune your collection.

### 1. Find Low-Quality Files
If you've upgraded your library over the years, you might still have old, low-bitrate files lurking in your collection.

```bash
# List all tracks with a bitrate lower than 320kbps
djlt ls rb/tracks "bitrate:<320"
```

### 2. Identify "Ghost" Tracks
Ghost tracks are songs that aren't in any of your playlists. They take up space but are never seen during a performance.

```bash
# Find tracks that belong to zero playlists
djlt ls rb/tracks "playlists:0"
```

!!! tip "Aggregation"
    You can combine these to find low-quality tracks that are also not in any playlists:
    ```bash
    djlt ls rb/tracks "bitrate:<320 && playlists:0"
    ```

### 3. Surface Unplayed Music
Sometimes the best way to clean up is to find the music you haven't touched in years.

```bash
# Find tracks with zero plays added more than 2 years ago
djlt ls rb/tracks "plays:0 && added:<-2y"
```

### 4. Pruning Membership
If you find a group of tracks that don't belong in a specific "Inbox" or "Process" crate anymore, you can unlink them in bulk without deleting them from your collection.

```bash
# Remove all 1-star tracks from the 'Recent Imports' playlist
djlt rm rb/tracks "rating:1" --from "rb/playlists name:'Recent Imports'"
```

---

## Safety First
By default, all `djlt` commands run in **Preview mode**. This allows you to see exactly what would happen without modifying your library.

Once you are satisfied with the preview, add the **`--apply`** flag to commit the changes:

```bash
djlt rm rb/tracks "rating:1" --from "..." --apply
```
