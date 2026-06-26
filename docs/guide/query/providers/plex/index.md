# Plex

The Plex provider interacts with your Plex Media Server via its API using the parallel prober for fast response times.

### Provider Alias: `plex`

## Resources

| Resource | Description | Example |
| :--- | :--- | :--- |
| `playlists` | Your Plex music playlists. | `plex/playlists name:Summer` |
| `tracks` | Tracks from your Plex music library. Use `playlists:` to scope to a specific crate or set of crates. | `plex/tracks playlists:Summer` or `plex/tracks title:'Yes'` |

## Selection Behavior

### Scoped Aggregation
When using the `playlists:` field with a substring match (`:`), the provider will aggregate tracks from **all** playlists that match the query.

For example, `plex/tracks playlists:DJ` will combine tracks from "DJ Crate 1", "DJ Crate 2", etc., into a single deduplicated list. To target a single playlist exclusively, use the exact match operator: `playlists="DJ Crate 1"`.

## Authentication

Plex requires a valid authentication token. You can set this up using:
```bash
djlt auth --plex
```

---

## Examples

**List a Plex playlist by name**
```bash
djlt ls plex/playlists name:Summer
```

**Sync a Plex playlist to Rekordbox**
```bash
djlt sync plex/playlists name:Summer --to "rb/playlists name:'Plex Sync'"
```

**Sync a Plex playlist to M3U8 (Planned)**
```bash
djlt sync plex/playlists name:Summer --to "m3u8:/path/to/playlist.m3u8"
```

