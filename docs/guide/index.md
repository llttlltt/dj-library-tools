# Getting Started

Welcome to **djlt**! This guide will help you get started with managing your DJ library, syncing with Plex, and performing powerful queries across your collection.

## Quick Installation

`djlt` is a Go-based tool and can be installed with a single command:

```bash
go install github.com/llttlltt/dj-library-tools/cmd/djlt@latest
```

!!! tip "Requirements"
    You must have [Go](https://go.dev/doc/install) installed on your system. If you plan on transcoding files, you will also need [FFmpeg](https://ffmpeg.org/).

## Basic Configuration

Before using `djlt`, you need to tell it where your Rekordbox XML library is located.

1. **Set your XML path:**
   ```bash
   djlt config rb file "/path/to/your/export.xml"
   ```
   This path is used by default for most commands.

2. **Authenticate with Plex (Optional):**
   If you want to sync playlists from Plex, run the authentication flow:
   ```bash
   djlt config plex auth
   ```
   Follow the link and enter the PIN provided in your terminal.

3. **Verify your setup:**
   Check that everything is configured correctly:
   ```bash
   djlt config list
   ```

## Seeing Your Music

Once configured, you can explore your library using the `list` command. `djlt` uses a consistent **[Selection Syntax](../query/syntax.md)** across all providers.

### Basic Searching
To list tracks from your Rekordbox library matching a specific criteria:

```bash
djlt list rb/tracks "artist:'Daft Punk'"
```

### Searching Specific Fields
You can combine fields and use boolean logic for complex searches:

```bash
djlt list rb/tracks "genre:House && bpm:124..128"
```

## Syncing Your Library

One of the most powerful features of `djlt` is its ability to sync data between providers.

### Sync Plex to Rekordbox
To take a playlist from Plex and inject it into your Rekordbox tree:

```bash
djlt sync plex/playlists name:Summer --to "rb/playlists name:'Plex Sync'"
```

### Export and Transcode
You can also export files to a local directory while syncing:

```bash
djlt sync plex/playlists name:Summer --to "rb/playlists name:'Plex Sync'" \
  --dest ~/Music/Export --format mp3
```

## Library Statistics

Get quick insights into any part of your library using the `--stats` flag:

```bash
djlt ls rb/tracks "rating:>=4" --stats
```

---

## Ready for more?

Now that you've mastered the basics, continue your journey:

- **[Syntax](../query/syntax.md)**: Master the query language and operators.
- **[Providers](../providers/index.md)**: Explore all available resources in Rekordbox and Plex.
- **[CLI Reference](../commands/index.md)**: Comprehensive guide to all commands and flags.
