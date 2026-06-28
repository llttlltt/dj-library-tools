# Library Repair & Enrichment

This workflow shows how to identify and fix physical health issues in your library, such as missing files or missing metadata tags. This replaces the old `fix` command.

### 1. Identify Missing Files
Broken links are a common issue when moving music between drives. Use the global `--missing` flag to find them.

```bash
# List all tracks in Rekordbox that are missing their physical file
djlt ls rb/tracks --missing
```

### 2. Automated Relocation
If you've moved your music to a new drive, use the `modify` command to repair the paths in bulk.

```bash
# Search a new directory and update paths for all missing tracks
djlt modify rb/tracks --missing --relocate "/Volumes/NewDrive/Music"
```

### 3. Cleanup dead references
If you want to permanently remove "ghost" tracks that no longer exist on disk.

```bash
# Remove all missing tracks from your Rekordbox collection
djlt rm rb/tracks --missing
```

### 4. Format Enrichment (Self-Healing)
For file-based providers like M3U, `djlt` uses a self-healing approach. Every time you perform an operation (like a `sync` or `mk`), the provider automatically ensures the output file has full metadata tags (`#EXTINF`).

```bash
# Simply 'fixing' an M3U8 file by re-saving it through the provider
djlt fix m3u8/playlists --file ./my_list.m3u8
```

---

## Proactive Health Checks
It's good practice to run `djlt ls --missing` after every major library move to catch broken links before you start your next performance.
