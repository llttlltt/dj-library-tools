# Creating and Updating Playlists

This guide shows you how to take any search criteria and turn it into a playlist that stays up to date, whether inside your main DJ software or as a separate file for your phone.

### 1. Build a playlist from a search
The most basic workflow is taking a set of tracks you've found and saving them as a new playlist.

```bash
# Find all 5-star tracks and save them to a new Rekordbox playlist
djlt mk rb/playlists "Top Tracks" --from "rb/tracks rating:5"
```

### 2. Keep a playlist in sync
As you add new music to your library, you want your playlists to update automatically. The `sync` command makes your target playlist match your search exactly.

```bash
# Refresh your "Top Tracks" playlist so it includes any new 5-star music
djlt sync "rb/tracks rating:5" --to "rb/playlists name:'Top Tracks'"
```

**How it works:**
*   **Adds**: Any new tracks that now match your search.
*   **Removes**: Any tracks that no longer match (e.g., if you changed a rating to 3 stars).
*   **Preserves**: Everything else stays exactly as it was.

### 3. Create playlists for other apps (Mobile/Phone)
If you want to move music to a phone or another app, you can "sync" your library to a standard playlist file (`.m3u8`).

```bash
# Create or update a playlist file for your mobile phone from your Rekordbox tracks
djlt sync "rb/tracks genre:House" --to "m3u8/tracks" --to-file "mobile.m3u8"
```

### 4. Moving music between apps
You can even use one app as the source for another. For example, if you have a playlist in Plex that you want to use in Rekordbox:

```bash
# Sync a Plex playlist directly into a Rekordbox playlist
djlt sync "plex/playlists name:Summer" --to "rb/playlists name:Summer"
```

---

## Tip: Append vs. Sync
By default, `sync` makes the target an exact match of the source. If you want to **add** new tracks but **never remove** anything you've manually added to the playlist, use the `--append` flag.

```bash
djlt sync "rb/tracks added:>today" --to "rb/playlists name:Inbox" --append
```
