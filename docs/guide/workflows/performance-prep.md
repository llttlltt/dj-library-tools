# Performance Preparation

Preparing for a set isn't just about picking tracks—it's about understanding the energy and range of your selection. This workflow uses `ls --stats` and `mk` to analyze and assemble performance-ready crates.

### 1. Analyze Your Crate
Before you play, check the "shape" of your playlist. Do you have enough variety in keys? Is the BPM range too narrow or too wide?

```bash
# Get a statistical summary of a specific playlist
djlt ls rb/tracks "playlists:Summer" --stats
```

The output will show you:
- **Average BPM**
- **Top 5 Keys** (to ensure Harmonic coverage)
- **Top 5 Genres**

### 2. Discover "Power Pairs"
Find tracks that match the energy of your favorite song. For example, if you love playing a 124 BPM track in 8A, find other "power pairs" in that same range.

```bash
# Find 4+ star tracks between 123-125 BPM in 8A or 8B
djlt ls rb/tracks "bpm:123..125 && (key:8A || key:8B) && rating:>=4"
```

### 3. Build the Set-List
Once you've refined a query that represents the "vibe" you want, you can generate a performance crate instantly.

```bash
# Create a new 'Peak Time' playlist from your query
djlt mk rb/playlists "Peak Time" --from "rb/tracks bpm:124..128 && rating:5"
```

### 4. Verify Cues
Ensure every track in your set has at least some cue points set for navigation.

```bash
# Find tracks in your 'Peak Time' crate with no hotcues
djlt ls rb/tracks "playlists:'Peak Time' && hotcues:0"
```

---

## Pro Tip: Scoped Stats
You can run stats on a query *within* a playlist. If you have a huge 'House' crate and want to know the stats for just the 5-star tracks:

```bash
djlt ls rb/tracks "playlists:House && rating:5" --stats
```
