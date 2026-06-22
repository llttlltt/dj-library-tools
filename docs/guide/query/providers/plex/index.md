# Plex (`plex`)

The Plex provider interacts with your Plex Media Server via its API using the parallel prober for fast response times.

## Resources

| Resource | Description | Example |
| :--- | :--- | :--- |
| `playlists` | Your Plex music playlists. | `plex/playlists:Summer` |

## Authentication

Plex requires a valid authentication token. You can set this up using:
```bash
djlt auth --plex
```
