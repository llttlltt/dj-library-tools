# Organizing Resources

This workflow shows how to manage the hierarchy of your Rekordbox library by creating folders and moving playlists or other folders into them.

### 1. Create a container folder
Folders are useful for grouping related playlists. Use the `mk` command with the `rb/folders` resource.

```bash
djlt mk rb/folders "Sorting"
```

### 2. Move existing playlists
To move one or more playlists into a folder, use the `mv` command. The `--to` flag must point to a specific folder destination.

```bash
# Move a single playlist by name
djlt mv rb/playlists 'name:"Inbox (Simple)"' --to "rb/folders name:Sorting"

# Move multiple playlists matching a query
djlt mv rb/playlists "name:Beatgrids" --to "rb/folders name:Sorting"
```

### 3. Nesting Folders & Deep Creation
You can create a whole folder hierarchy at once using the `--parents` flag.

```bash
# Create a deep folder path in one go
djlt mk rb/folders "Archive/2024/Jan" --parents
```

### 4. Recursive Cleanup
To remove a folder and all its contents (playlists and sub-folders), use the `rm` command with the `--recursive` flag.

```bash
# Delete a folder and everything inside it
djlt rm rb/folders "name:Old Archives" --recursive
```

### 5. Renaming while moving
You can combine a move with a rename by using the `--name` flag.

```bash
# Move a playlist and rename it in one step
djlt mv rb/playlists "name:Inbox" --to "rb/folders name:Sorting" --name "To Process"
```

---

## Pro Tip: Full Resource Paths
When using `--to` or `--from` in the `mv` command, always include the full provider and resource path (e.g., `rb/folders name:...`) to ensure the tool identifies the correct target type.
