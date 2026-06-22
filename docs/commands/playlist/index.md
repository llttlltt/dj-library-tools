# playlist

Manage rekordbox playlists

```
djlt playlist [rb/playlists query] [flags]
```

### Options

```
      --add string      Add tracks matching this rb/tracks: query
      --delete          Delete matched playlists
      --dry-run         Preview changes without writing
      --folder string   Parent folder for --new (default: root level)
  -h, --help            help for playlist
      --move string     Move matched playlists into this folder
      --new string      Create a new playlist with this name
      --remove string   Remove tracks matching this rb/tracks: query from matched playlists
      --rename string   Rename matched playlists to this name
      --sync string     Sync matched playlists to exactly match this rb/tracks: query
```

### Inherited Options

```
  -x, --xml string   Path to the Rekordbox XML library
```

## See also

* [djlt](../index.md)	 - DJ Library Tools CLI
* [djlt playlist fix](fix.md)	 - Fix playlist extensions and/or enrich with M3U8 metadata