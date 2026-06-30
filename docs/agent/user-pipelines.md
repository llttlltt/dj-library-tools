# Maintenance & Sorting

Instructions for keeping the library's sorting and inbox hierarchy synchronized.

## Sorting Folder Maintenance

These maintenance routines are now available as **Workflows** in the GUI (`cmd/djlt-gui`). They can be executed by selecting the workflow and pressing **Preview** to verify changes before clicking **Run**.

| Workflow | Source Query | Target Playlist Pattern |
|----------|--------------|-------------------------|
| Existing Inbox | `beatgrids/bpm-redundancy:>80 && beatgrids-count:>1` | `Existing.*Inbox` |
| Memory Cue Inbox | `beatgrids-count:1 && !memorycues/name:S.O.S.` | `Memory.*Cue.*Inbox` |
| Tagging Inbox (Beatgrids:=1) | `beatgrids-count:1 && memorycues/name:S.O.S.` | `Tagging.*Inbox.*Beatgrids:=1` |
| Tagging Inbox (Beatgrids:>1) | `beatgrids-count:>1 && memorycues/name:S.O.S.` | `Tagging.*Inbox.*Beatgrids:>1` |
| 5 - Beatgrids:=1 (Total) | `beatgrids-count:1` | `5.*Beatgrids:=1$` |
| 5 - Beatgrids:>1 (Total) | `beatgrids-count:>1` | `5.*Beatgrids:>1$` |

**Mandatory Protocol**: Always run the **Preview** first to verify the track additions and removals. Only click **Run** after confirming the diff is correct.
