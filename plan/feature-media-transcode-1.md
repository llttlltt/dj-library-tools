---
goal: Implement the Media Transcode & Export module for djlt.
version: 1.0.0
date_created: 2026-06-22
last_updated: 2026-06-22
owner: Elliott
status: 'Planned'
tags: [feature, media, transcode, ffmpeg, export]
---

# Introduction

![Status: Planned](https://img.shields.io/badge/status-Planned-blue)

To achieve full parity with the Beets `convert` plugin, `djlt` requires a transcode engine that can take lossless files (FLAC, ALAC, WAV, AIFF) and convert them to high-quality lossy formats (MP3 320kbps) while maintaining metadata, album art, and specific naming conventions.

## 1. Requirements & Constraints

- **REQ-001**: Support transcoding from FLAC, ALAC, WAV, AIFF to MP3 (320kbps, 44.1kHz).
- **REQ-002**: Maintain high-integrity metadata (ID3v2.3) and embed album art.
- **REQ-003**: Implement customizable path formatting (e.g., `Artist - Album - Title`).
- **REQ-004**: Use `ffmpeg` as the underlying engine but provide a clean Go wrapper.
- **REQ-005**: Support "Smart Copy": if a file is already the target format and meets criteria, copy it instead of transcoding.
- **REQ-006**: Parallel processing: transcode multiple files simultaneously.

## 2. Implementation Steps

### Phase 1: Transcode Engine

- GOAL-001: Build the core FFmpeg wrapper for transcoding.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-001 | Create `internal/media/config.go` for Beets-compatible configuration. | ✅ | 2026-06-22 |
| TASK-002 | Create `internal/media/ffmpeg.go` to handle FFmpeg execution. | ✅ | 2026-06-22 |
| TASK-003 | Implement `internal/media/path.go` for templated file naming. | ✅ | 2026-06-22 |

### Phase 2: Path Formatting & Export

- GOAL-002: Implement Beets-style path templating.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-004 | Create `internal/media/path.go` for templated file naming. | | |
| TASK-005 | Implement `Export` orchestrator to manage the transcode queue. | | |

### Phase 3: Integration

- GOAL-003: Connect transcoding to the Plex Sync flow.

| Task | Description | Completed | Date |
|------|-------------|-----------|------|
| TASK-006 | Add `export` flags to `sync plex` to trigger transcoding. | | |
| TASK-007 | Create `djlt media convert` command for standalone usage. | | |

## 3. Alternatives

- **ALT-001**: Use a pure Go audio library (e.g., `go-audio`). Rejected because FFmpeg is the industry standard for robust transcoding with complex metadata and art embedding.

## 4. Dependencies

- **DEP-001**: `ffmpeg` must be installed on the host system.

## 5. Files

- **FILE-001**: `internal/media/ffmpeg.go`
- **FILE-002**: `internal/media/path.go`
- **FILE-003**: `cmd/djlt/media.go`

## 6. Testing

- **TEST-001**: Verify transcode output quality and bitrates.
- **TEST-002**: Check ID3 tags and album art in the resulting MP3s.

## 7. Risks & Assumptions

- **RISK-001**: FFmpeg version variations.
- **ASSUMPTION-001**: Users want 320kbps MP3 as the primary export format for Rekordbox compatibility.

## 8. Related Specifications / Further Reading

- [architecture-djlt-system-1.md](./architecture-djlt-system-1.md)
- [Beets Convert Plugin](https://beets.readthedocs.io/en/stable/plugins/convert.html)
