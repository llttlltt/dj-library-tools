# Library Repair & Enrichment

This workflow shows how to identify and fix physical health issues in your library, such as missing files or missing metadata tags.

### 1. Identify Missing Files
Broken links are a common issue when moving music between drives. Use the global `--missing` flag to find them.

```bash
# List all tracks in Rekordbox that are missing their physical file
djlt ls rb/tracks --missing
```

### 2. Automated Relocation
If you've moved your music to a new drive, use the `fix` command to repair the paths in bulk.

```bash
# Attempt to find missing files by checking for common DJ extensions
djlt fix rb/tracks --missing --paths normalize --apply
```

### 3. Cleanup dead references
If you want to permanently remove "ghost" tracks that no longer exist on disk.

```bash
# Remove all missing tracks from your Rekordbox collection
djlt rm rb/tracks --missing
```

### 4. Automated Maintenance

The `fix` command is a unified way to perform algorithmic repairs on your library.

```bash
# Perform a full library health sweep (duplicates, normalization, path repair)
djlt fix rb/tracks --duplicates tracks --metadata all --paths normalize --apply
```

#### Deduplication
You can remove duplicate tracks from the master collection or duplicate memberships within a playlist.

```bash
# Remove duplicate track IDs from a specific playlist
djlt fix "rb/playlists name:MyPlaylist" --duplicates members --apply

# Remove duplicate physical tracks from the entire collection
djlt fix rb/tracks --duplicates tracks --apply
```

#### Metadata Normalization
Clean up messy tags like trailing whitespace or "None" comments.

```bash
# Normalize all metadata fields
djlt fix rb/tracks --metadata all --apply
```

---

## Proactive Health Checks
It's good practice to run `djlt ls --missing` after every major library move to catch broken links before you start your next performance.
