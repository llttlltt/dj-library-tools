# Automating Playlists

This workflow shows how to use the `sync` command to manage a "to-do" playlist that stays in sync with your track criteria.

### 1. Verify your query
Before doing anything else, use the `list` command to see which tracks currently need work. This ensures your query is selecting exactly what you expect.

```bash
djlt list rb/tracks "beatgrids:1 && hotcues:<1 && memorycues:<1"
```

### 2. Create the playlist
Next, create an empty playlist in Rekordbox to hold these tracks.

```bash
djlt mk rb/playlists "Inbox (Simple)"
```

### 3. Add the tracks
Now, populate your new playlist with the tracks found in step 1. Use `--append` so existing members are never removed.

```bash
djlt sync rb/tracks "beatgrids:1 && hotcues:<1 && memorycues:<1" \
  --to "rb/playlists name:'Inbox (Simple)'" --append
```

!!! tip "All-in-one"
    You can combine the creation and population steps into a single command:
    ```bash
    djlt mk rb/playlists "Inbox (Simple)" --from "rb/tracks beatgrids:1 && hotcues:<1 && memorycues:<1"
    ```

### 4. Work in Rekordbox
Open Rekordbox and start working on the tracks. As you add **Hot Cues**, **Memory Cues**, or fix the **Beatgrid** (adding more markers), those tracks no longer match your "simple" criteria.

### 5. Synchronize the playlist
Whenever you want to refresh your list, run the `sync` command. `djlt` will compare the current state of your library against the playlist and bring it back in line with your query.

!!! tip "Dry Run first"
    It's a good habit to use the `--dry-run` flag first to see exactly how many tracks will be added or removed before committing the changes to your XML:
    ```bash
    djlt sync rb/tracks "..." --to "rb/playlists name:Inbox" --dry-run
    ```

```bash
djlt sync rb/tracks "beatgrids:1 && hotcues:<1 && memorycues:<1" \
  --to "rb/playlists name:'Inbox (Simple)'"
```

**How it works:**

1.  **Removes**: Any tracks you have already fixed (because they no longer match the query).
2.  **Preserves**: Any tracks that still match the query, keeping them in their current order.
3.  **Adds**: Any new tracks that now match the query (appending them to the end).

---

## Why use `sync`?
While you could delete and recreate the playlist each time, the `sync` command is **surgical**. It only adds or removes the tracks that are relevant to the query.

If you have manually dragged other tracks into that playlist (like reference tracks or favourites), `sync` will leave them alone — whereas recreating the playlist would delete them.
