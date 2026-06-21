---
goal: Implement the Sync Engine to reconcile Plex playlists with Rekordbox XML.
version: 1.0.0
date_created: 2026-06-22
last_updated: 2026-06-22
owner: Elliott
status: 'Planned'
tags: [feature, sync, rekordbox, plex]
---

# Introduction

![Status: Planned](https://img.shields.io/badge/status-Planned-blue)

The Sync Engine is the bridge between Plex metadata and the Rekordbox XML library. It matches tracks from Plex playlists to existing tracks in the Rekordbox library and generates/updates Rekordbox playlists.

## 1. Requirements & Constraints

- **REQ-001**: Match tracks using fuzzy or exact matching on Title, Artist, and Album.
- **REQ-002**: Support reconciliation of file paths (Plex paths vs. Local/Rekordbox paths).
- **REQ-003**: Create or update Rekordbox XML Playlists based on Plex playlist data.
- **REQ-004**: Handle "missing" tracks (tracks in Plex but not in Rekordbox) gracefully.
- **PAT-001**: Use the `rekordbox.Track` models from `pkg/rekordbox`.

## 2. Implementation Steps

### Phase 1: Matching Logic

- GOAL-001: Implement the track matching core.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Create `internal/sync/matcher.go` with track matching logic. | | |
| TASK-002 | Implement scoring/fuzzy matching for Artist/Title. | | |

### Phase 2: Reconciliation & Export

- GOAL-002: Reconcile paths and generate Rekordbox XML output.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-003 | Create `internal/sync/engine.go` to orchestrate the sync process. | | |
| TASK-004 | Implement path mapping logic (e.g., prefix replacement for Plex-to-Local). | | |
| TASK-005 | Implement `UpdateXML` to inject new playlists into a `rekordbox.XML` object. | | |

### Phase 3: CLI Integration

- GOAL-003: Expose sync functionality via the CLI.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-006 | Create `cmd/djlt/sync.go` with `sync plex` command. | | |

## 3. Alternatives

- **ALT-001**: Direct database injection. Rejected; XML is safer and matches Rekordbox's primary import method.

## 4. Dependencies

- **DEP-001**: `pkg/rekordbox` - XML models.
- **DEP-002**: `internal/plex` - Plex data retrieval.

## 5. Files

- **FILE-001**: `internal/sync/matcher.go`
- **FILE-002**: `internal/sync/engine.go`
- **FILE-003**: `cmd/djlt/sync.go`

## 6. Testing

- **TEST-001**: Unit tests for matching logic with various metadata edge cases.
- **TEST-002**: Integration test verifying XML playlist generation.

## 7. Risks & Assumptions

- **RISK-001**: Metadata mismatch between Plex and Rekordbox (e.g., "Feat." vs "ft.").
- **ASSUMPTION-001**: Users have already imported their music into both Plex and Rekordbox.

## 8. Related Specifications / Further Reading

- [architecture-djlt-system-1.md](./architecture-djlt-system-1.md)
