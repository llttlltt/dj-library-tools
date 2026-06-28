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

### Playlist & Folder Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `name` | String | Name of the ResourceGroup. |
| `parent` | String | Name of the parent folder. |
| `items` | Numeric | Number of tracks in a playlist, or child ResourceGroups in a folder. |
| `kind` | String | `folder` or `playlist`. |

## Collections (Path Queries)

Rekordbox supports advanced path-based queries for deep analysis of cues and beatgrids. 

**Syntax**: `Collection . Index / Property - Stat`

| Collection | Description | Properties | Stats |
| :--- | :--- | :--- | :--- |
| `beatgrids` | Beatgrid markers and tempo information. | `bpm`, `position` | `-count`, `-drift`, `-density` |
| `hotcues` | Performance pads (A-H). | `color`, `name`, `position` | `-count` |
| `memorycues` | Standard markers. | `color`, `name`, `position` | `-count` |
| `playlists` | Membership in playlists. | `name`, `folder` | `-count` |

### Property Reference

| Property | Type | Description |
| :--- | :--- | :--- |
| `bpm` | Numeric | Beats per minute at the marker. |
| `position` | Numeric | Time in seconds from the start of the track. |
| `color` | String | Color name (e.g., `red`, `skyblue`). |
| `name` | String | User-assigned label or comment. |
| `folder` | String | Name of the folder containing the playlist (Playlists only). |

### Indexing

Collections support **1-based indexing** to target a specific item. 

- `hotcues.1/color`: The color of the first hotcue.
- `playlists.2/name`: The name of the second playlist the track belongs to.

If no index is provided (e.g., `playlists/name:House`), the query returns true if **any** item in the collection matches.

### Stat Reference

| Stat | Type | Description |
| :--- | :--- | :--- |
| `-count` | Numeric | Total number of items in the collection. |
| `-min` | Numeric | Minimum value found in the collection. |
| `-max` | Numeric | Maximum value found in the collection. |
| `-avg` | Numeric | Mathematical average of values. |
| `-drift` | Numeric | Difference between the Maximum and Minimum values. |
| `-jitter` | Numeric | Average change between consecutive items. |
| `-redundancy` | Numeric | Percentage (0-1) of items identical to the previous one. |
| `-stability` | Numeric | Composite score (0-100) where 100 is perfectly steady. |
| `-density` | Numeric | Items per minute (based on track duration). |

## Color Palettes

Different palettes are used for Tracks and Cues to match Rekordbox's UI.

### Track Colors

The following colors are available for track-level filtering (e.g. `color:red`):

- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#FF007F;margin-right:5px;"></span> `pink`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#FF0000;margin-right:5px;"></span> `red`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#FFA500;margin-right:5px;"></span> `orange`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#FFFF00;margin-right:5px;"></span> `yellow`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#00FF00;margin-right:5px;"></span> `green`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#25FDE9;margin-right:5px;"></span> `aqua`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#0000FF;margin-right:5px;"></span> `blue`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#660099;margin-right:5px;"></span> `purple`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:transparent;margin-right:5px;"></span> `none`

### Hot Cue Colors

Use the following names to match the 16-color pad palette (e.g. `hotcues.1/color:skyblue`). Cues with no color set match `none`.

<table style="border-collapse: separate; border-spacing: 5px;">
  <tr>
    <td style="background:#DE44CF; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>hotpink</code></td>
    <td style="background:#B432FF; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>purple</code></td>
    <td style="background:#AA72FF; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>violet</code></td>
    <td style="background:#6473FF; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>indigo</code></td>
  </tr>
  <tr>
    <td style="background:#305AFF; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>blue</code></td>
    <td style="background:#50B4FF; color:black; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>skyblue</code></td>
    <td style="background:#00E0FF; color:black; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>aqua</code></td>
    <td style="background:#1FA392; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>darkgreen</code></td>
  </tr>
  <tr>
    <td style="background:#10B176; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>brightgreen</code></td>
    <td style="background:#28E214; color:black; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>green</code></td>
    <td style="background:#A5E116; color:black; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>yellowgreen</code></td>
    <td style="background:#B4BE04; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>yellow</code></td>
  </tr>
  <tr>
    <td style="background:#C3AF04; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>orange</code></td>
    <td style="background:#E0641B; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>darkorange</code></td>
    <td style="background:#E62828; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>red</code></td>
    <td style="background:#FF127B; color:white; padding:10px; text-align:center; border-radius:4px; width:25%;"><code>pink</code></td>
  </tr>
</table>

---

## Examples

### Deep Analysis (Paths)

**Identify unstable "dynamic" grids**
```bash
djlt ls rb/tracks "beatgrids/bpm-drift:>0.5 && beatgrids-count:>10"
```

**Find tracks with red HotCue A**
```bash
djlt ls rb/tracks "hotcues.1/color:red"
```

**Find "busy" grids with high marker density**
```bash
djlt ls rb/tracks "beatgrids-density:>60"
```

### Basic Metadata

**High-energy House**
```bash
djlt ls rb/tracks "genre:House && bpm:124..128 && rating:>=4"
```

**Tracks not in any playlist**
```bash
djlt ls rb/tracks "playlists-count:0"
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
