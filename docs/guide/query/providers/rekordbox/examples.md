# Rekordbox Examples

## Tracks

**High-energy House**
```bash
djlt list rb/tracks "genre:House && bpm:124..128 && rating:>=4"
```

**Tracks not in any playlist**
```bash
djlt list rb/tracks "playlistcount:0"
```

**Find a specific track by ID**
```bash
djlt list rb/tracks "id:1234"
```

## Collection Tree

**Find folders containing "Sets"**
```bash
djlt list rb/folders "name:Sets"
```

**Find playlists with "2023" in the name**
```bash
djlt list rb/playlists "name:2023"
```
