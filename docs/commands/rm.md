# rm

Permanently delete resources or remove membership

```
djlt rm [resource] [query] [flags]
```
### Options

```
      --from strings   Origin resource(s) to remove from
  -h, --help           help for rm
  -r, --recursive      Delete folder and all its contents
```

### Inherited Options

```
      --apply         Actually apply changes to the library (destructive)
  -f, --file string   Path to the primary library file (Rekordbox XML, M3U, etc.)
  -v, --verbose       Enable verbose logging
```

Permanently delete resources or remove track membership from playlists.

Use --from to specify which playlist to remove tracks from.
Without --from, the command deletes the resource itself.

Example:
  djlt rm rb/tracks "artist:Four" --from "rb/playlists name:Inbox"
  djlt rm rb/playlists name:Inbox
  djlt rm rb/folders name:OldSets --recursive

## See also

* [djlt](index.md)	 - DJ Library Tools CLI