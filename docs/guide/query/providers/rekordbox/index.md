# Rekordbox (`rb`)

The Rekordbox provider interacts directly with your exported XML library. It allows you to query your track collection, playlists, and folder structure.

## Resources

| Resource | Description | Example |
| :--- | :--- | :--- |
| `tracks` | The main track collection. | `rb/tracks:bpm:124` |
| `playlists` | Individual playlist nodes. | `rb/playlists:name:Summer` |
| `folders` | Folder nodes in the playlist tree. | `rb/folders:name:House` |

## Fields

### Track Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `name`, `title` | String | Track title. |
| `artist` | String | Artist name. |
| `album` | String | Album name. |
| `genre` | String | Genre name. |
| `bpm`, `tempo` | Numeric | Beats per minute. |
| `rating` | Numeric | Star rating (0-5 stars). |
| `playcount` | Numeric | Number of times played. |
| `playlist` | String | Matches tracks in a specific playlist. |
| `playlistcount` | Numeric | Number of playlists a track belongs to. |
| `year` | Numeric | Release year. |
| `key` | String | Musical key (Tonality). |
| `comment` | String | Track comments. |
| `label` | String | Record label. |
| `remixer` | String | Remixer name. |
| `mix`, `version` | String | Track mix/version. |
| `added` | String | Date added to collection. |
| `time`, `length` | Numeric | Total time in seconds. |
| `bitrate`, `kbps` | Numeric | Audio bitrate. |
| `size` | Numeric | File size in bytes. |

### Playlist & Folder Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `name` | String | Name of the node. |
| `folder`, `parent` | String | Name of the parent folder. |
| `entries` | Numeric | Number of tracks or sub-items. |
| `type` | Numeric | `0` for folder, `1` for playlist. |

## Cue Filtering

You can query tracks based on their HotCues and MemoryCues using specific count and property filters.

### Cue Counts

| Field | Description | Example |
| :--- | :--- | :--- |
| `hotcues` | Number of HotCues. | `hotcues:>3` |
| `memorycues` | Number of MemoryCues. | `memorycues:0` |

### Specific Cue Properties

You can target a specific cue by its ID and check its properties.

**Syntax**: `field:ID[:Property[:Value]]`

| Resource | IDs | Example |
| :--- | :--- | :--- |
| `hotcue` | `a` through `h` | `hotcue:a:red` |
| `memorycue` | `1`, `2`, ... | `memorycue:1:label:Intro` |

### Available Properties

| Property | Value | Description |
| :--- | :--- | :--- |
| **Color** | `red`, `blue`, etc. | Match by cue color. |
| `label` | `Text` | Substring match on the cue label. |
| `loop` | (none) | Match if it is an active loop. |

**Available Colors**: `red`, `orange`, `yellow`, `green`, `aqua`, `blue`, `purple`, `pink`, `none`.
