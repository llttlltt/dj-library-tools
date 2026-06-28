# Metadata Hardening

A "hardened" library is one where every track is fully analyzed, cue-pointed, and correctly tagged. This workflow helps you identify gaps in your metadata and reconcile them in bulk.

### 1. Identify Missing Beatgrids
Analyze which tracks are missing beatgrid information (Tempo markers). These are tracks you haven't "locked in" for syncing yet.

```bash
# Find tracks with zero beatgrids
djlt ls rb/tracks "beatgrids:0"
```

### 2. Standardize Color Coding
Use color coding to flag tracks that need specific types of work. For example, use **Red** for tracks that have broken metadata and **Green** for those that are "Performance Ready".

```bash
# Find all "Red" tracks that you've since analyzed (have hotcues)
# and need to be marked as "Green"
djlt ls rb/tracks "color:red && hotcues:>3"
```

### 3. Reconcile from External Sources
If you have metadata (like ratings or grids) in a different library, use the `sync` command to reconcile it without moving files.

```bash
# Sync ratings from Plex into your Rekordbox collection
djlt sync plex/tracks --to rb/tracks --metadata rating
```

### 4. Bulk Processing
If you find a group of tracks that are missing a specific tag, you can move them into a temporary "To Tag" playlist to process them in Rekordbox.

```bash
# Move all un-tagged tracks into a processing playlist
djlt mv rb/tracks "comment:none" --from "rb/playlists name:Inbox" --to "rb/playlists name:'To Tag'"
```

---

## The analyze-tag loop
The most effective way to harden your library is the query-reconcile loop:
1. Run `ls` with a specific metadata gap (e.g. `hotcues:0`).
2. Fix the metadata in your provider or reconcile from a backup using `sync --metadata`.
3. Repeat until the query returns zero results.
