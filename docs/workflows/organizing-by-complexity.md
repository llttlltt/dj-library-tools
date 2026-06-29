# Organizing by Complexity

This workflow shows how to categorize your library based on track complexity—specifically by how many Beatgrid markers are required.

### 1. Identify "Simple" vs "Complex" tracks
Many tracks have a steady tempo and only require a single Beatgrid marker. Others may have tempo shifts or live drums, requiring multiple markers. You can use the `beatgrids` field to identify these.

```bash
# View tracks with a single beatgrid
djlt ls rb/tracks "beatgrids:1"

# View tracks with multiple beatgrid markers
djlt ls rb/tracks "beatgrids:>1"
```

### 2. Create complexity-based playlists
Use the `mk` command with the `--from` flag to create and populate playlists based on these criteria.

```bash
# Create the "Simple Beatgrids" playlist
djlt mk rb/playlists "Simple Beatgrids" --from "rb/tracks beatgrids:1"

# Create the "Complex Beatgrids" playlist
djlt mk rb/playlists "Complex Beatgrids" --from "rb/tracks beatgrids:>1"
```

### 3. Maintain the organization
As you add new music to your library, you can periodically refresh these playlists using the `sync` command to ensure they always reflect your current beatgrid state.

```bash
djlt sync rb/tracks "beatgrids:>1" --to "rb/playlists name:'Complex Beatgrids'"
```

!!! tip "Organizing further"
    Once created, you can move these playlists into a folder for better organization. See the **[Organizing Resources](organizing-resources.md)** workflow for more details on managing your library hierarchy.
