# Workflow Engine

The Workflow engine executes multi-step library automation routines.

## Step Kinds

| Kind | Description | Options |
|------|-------------|---------|
| `sync` | Synchronize membership and metadata between endpoints. | `metadata`, `match`, `append_only` |
| `add` | Create a new resource and/or add tracks from source. | `in`, `at`, `kind` |
| `remove` | Remove tracks from target groups. | `recursive` |
| `edit` | Batch update track metadata. | `set` (map of field: value) |
| `fix` | Algorithmic maintenance (duplicates, orphans). | `actions` |
| `m3u_export` | Export track selection to an M3U file. | `path`, `append` |

## Ad-hoc URI Schemes

Workflows support ad-hoc URIs to reference files directly without pre-registered connection UUIDs.

- `m3u://<absolute_path>`: Reference an M3U playlist file.
- `m3u8://<absolute_path>`: Reference an M3U8 playlist file (UTF-8).

Example Endpoints:
- `m3u:///tmp/sorting.m3u`
- `m3u8:///Users/dj/Music/Playlists/Fresh.m3u8 tracks`

## Targeting

The `add` and `remove` steps use the `Targets` slice to identify which groups are affected. 

- For `add`, each target `Query` is the **name** of the new group to create/populate.
- For `remove`, each target identifies an **existing group** to remove tracks from.
