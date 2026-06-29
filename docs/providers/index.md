# Providers

`djlt` supports multiple library providers. Each provider exposes specific resources that can be queried using the selection syntax.

## Available Providers

- **[Rekordbox (rb)](rekordbox.md)**: Manage your local collection, playlists, and folders.
- **[Plex (plex)](plex.md)**: Sync data from your Plex Media Server.
- **[M3U / M3U8 (m3u/m3u8)](m3u.md)**: Read and write standard music playlist files.

---

### Selection Pattern
`djlt VERB [provider/resource] [query]`

The resource and query are space-separated positional arguments — the query is **not** colon-joined to the resource path. For example:

```bash
djlt ls rb/tracks "genre:House && bpm:124..128"
djlt ls rb/playlists "name:Inbox"
djlt ls rb/folders "name:Shows"
```
