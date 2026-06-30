# Maintenance & Sorting

Instructions for keeping the library's sorting and inbox hierarchy synchronized.

## Sorting Folder Maintenance

**Mandatory Protocol**: Always run these commands **without** `--apply` first to verify the diff. Use `--verbose` to see specific track additions and removals.

1. **Existing Inbox (Beatgrids:>1 & BPM-Redundancy:>80)**:

   ```bash
   ./djlt sync "rb/tracks" "beatgrids/bpm-redundancy:>80 && beatgrids-count:>1" --to "rb/playlists name::Existing.*Inbox" --verbose
   ```

2. **Memory Cue Inbox (Beatgrids:=1 & !Cue=:S.O.S.)**:

   ```bash
   ./djlt sync "rb/tracks" "beatgrids-count:1 && !memorycues/name:S.O.S." --to "rb/playlists name::Memory.*Cue.*Inbox" --verbose
   ```

3. **Tagging Inbox (Beatgrids:=1 & Cue:=S.O.S.)**:

   ```bash
   ./djlt sync "rb/tracks" "beatgrids-count:1 && memorycues/name:S.O.S." --to "rb/playlists name::Tagging.*Inbox.*Beatgrids:=1" --verbose
   ```

4. **Tagging Inbox (Beatgrids:>1 & Cue:=S.O.S.)**:

   ```bash
   ./djlt sync "rb/tracks" "beatgrids-count:>1 && memorycues/name:S.O.S." --to "rb/playlists name::Tagging.*Inbox.*Beatgrids:>1" --verbose
   ```

5. **Beatgrids:=1 (Total)**:

   ```bash
   ./djlt sync "rb/tracks" "beatgrids-count:1" --to "rb/playlists name::5.*Beatgrids:=1$" --verbose
   ```

6. **Beatgrids:>1 (Total)**:

   ```bash
   ./djlt sync "rb/tracks" "beatgrids-count:>1" --to "rb/playlists name::5.*Beatgrids:>1$" --verbose
   ```

*Note: Add `--apply` only after verifying the verbose output and requesting permission or are explicitly told so.*
