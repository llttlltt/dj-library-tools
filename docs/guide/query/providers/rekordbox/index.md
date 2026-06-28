# Rekordbox

The Rekordbox provider interacts directly with your exported XML library. It allows you to query your track collection, playlists, and folder structure.

## Provider

| Alias | Description |
| :--- | :--- |
| `rb` | Rekordbox library provider. |

## Resources

| Resource | Description | Example |
| :--- | :--- | :--- |
| `tracks` | The main track collection. | `rb/tracks artist:Four` |
| `playlists` | Individual playlist ResourceGroups. | `rb/playlists name:Inbox` |
| `folders` | Folder ResourceGroups in the playlist tree. | `rb/folders name:House` |

## Fields

### Track Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `title` | String | Track title. |
| `artist` | String | Artist name. |
| `album` | String | Album name. |
| `genre` | String | Genre name. |
| `bpm` | Numeric | Beats per minute. |
| `rating` | Numeric | Star rating (Standardized 0-255). |
| `plays` | Numeric | Number of times played. |
| `playlists` | Mixed | Match by playlist name (String) or number of memberships (Numeric). |
| `year` | Numeric | Release year. |
| `key` | String | Musical key (Tonality). |
| `comment` | String | Track comments. |
| `label` | String | Record label. |
| `remixer` | String | Remixer name. |
| `mix` | String | Track mix/version. |
| `added` | String | Date added to collection. |
| `modified` | String | Date track was last modified. |
| `duration` | Numeric | Total time in seconds. |
| `bitrate` | Numeric | Audio bitrate. |
| `samplerate` | Numeric | Sample rate in Hz. |
| `size` | Numeric | File size in bytes. |
| `beatgrids` | Numeric | Number of beatgrid markers. |

### Playlist & Folder Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `name` | String | Name of the ResourceGroup. |
| `parent` | String | Name of the parent folder. |
| `items` | Numeric | Number of tracks in a playlist, or child ResourceGroups in a folder. |
| `kind` | String | `folder` or `playlist`. |

## Deep Metadata (Path Queries)

Rekordbox supports advanced path-based queries for deep analysis of cues and beatgrids.

**Syntax**: `Collection . Index / Property - Stat`

### Beatgrids

Find tracks that might need re-analysis (many markers but no BPM change):
`rb/tracks "beatgrids/bpm-drift:<0.1 && beatgrids-count:>10"`

Find "busy" variable grids:
`rb/tracks "beatgrids-density:>60"`

| Path | Description |
| :--- | :--- |
| `beatgrids-count` | Total number of markers. |
| `beatgrids-density` | Markers per minute of track duration. |
| `beatgrids/bpm-drift` | Difference between Max and Min BPM markers. |
| `beatgrids.N/bpm` | BPM value of the Nth marker. |
| `beatgrids.N/position` | Position (seconds) of the Nth marker. |

### Cues (HotCues & MemoryCues)

Search for specific pad colors or named sections:
`rb/tracks "hotcues.1/color:red && memorycues/name:Break"`

| Path | Description |
| :--- | :--- |
| `hotcues-count` | Total number of hotcues. |
| `hotcues/color` | Match if ANY hotcue has this color. |
| `hotcues.N/color` | Color of the Nth hotcue. |
| `hotcues.N/name` | Name of the Nth hotcue. |

---

## Examples

### Tracks

**High-energy House**
```bash
djlt ls rb/tracks "genre:House && bpm:124..128 && rating:>=4"
```

**Tracks not in any playlist**
```bash
djlt ls rb/tracks "playlists:0"
```

**Find a specific track by ID**
```bash
djlt ls rb/tracks "id:1234"
```

### Collection Tree

**Find folders containing "Sets"**
```bash
djlt ls rb/folders "name:Sets"
```

**Find playlists with "2023" in the name**
```bash
djlt ls rb/playlists "name:2023"
```

