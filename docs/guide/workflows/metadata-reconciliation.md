# Metadata Reconciliation

This workflow shows how to use the `sync` command to reconcile metadata (like Beatgrids, Ratings, and Comments) between different libraries or providers. This replaces the old `update` command.

### 1. Identify the Source
Determine which library has the "Golden" metadata. For example, you might have fixed beatgrids in a backup XML file, or you might have rated your tracks in Plex.

### 2. Match and Preview
Use the `sync` command with the `--metadata` flag. By default, `djlt` matches tracks by **Artist** and **Title**.

```bash
# Preview syncing beatgrids from a backup XML to your primary library
djlt sync "rb/tracks" --file backup.xml --to "rb/tracks" --metadata beatgrids --dry-run
```

### 3. Change Match Criteria
If your filenames are identical but your tags are slightly different, you can match by **Filename** instead.

```bash
# Sync ratings from Plex to Rekordbox, matching by filename
djlt sync "plex/tracks" --to "rb/tracks" --metadata rating --match filename
```

### 4. Selective Reconciliation
You don't have to sync everything. You can target specific tracks using a query.

```bash
# Reconcile comments for only your "Techno" tracks
djlt sync "rb/tracks genre:Techno" --file source.xml --to "rb/tracks" --metadata comment
```

---

## Why use `sync` for metadata?
The `sync` command is agnostic. It allows you to move metadata between completely different systems (like Plex and Rekordbox) while ensuring that only unambiguous matches are updated. If a source track matches multiple targets, `djlt` will skip it for safety unless you use `--match-force`.
