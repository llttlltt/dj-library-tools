# Providers

`djlt` supports multiple library providers. Each provider exposes specific resources that can be queried using the selection syntax.

## Available Providers

- **[Rekordbox (rb)](rekordbox/index.md)**: Manage your local collection, playlists, and folders.
- **[Plex (plex)](plex/index.md)**: Sync data from your Plex Media Server.

---

### Selection Pattern
`djlt VERB [provider/resource] [query]`

The resource and query are space-separated positional arguments — the query is **not** colon-joined to the resource path. For example:

```bash
djlt ls rb/tracks "genre:House && bpm:124..128"
djlt ls rb/playlists "name:Inbox"
djlt ls rb/folders "name:Shows"
```
