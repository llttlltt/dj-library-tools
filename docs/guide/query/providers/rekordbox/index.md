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

## Cue Filtering

You can query tracks based on their HotCues and MemoryCues using numeric count filters or custom color matches.

### Cue Counts

| Field | Description | Example |
| :--- | :--- | :--- |
| `hotcues` | Number of HotCues. | `hotcues:>3` |
| `memorycues` | Number of MemoryCues. | `memorycues:0` |

### Custom Color Matching

You can search for Hot Cues by color name using the `hotcues` field.

| Example | Description |
| :--- | :--- |
| `hotcues:red` | Find tracks with at least one Red hot cue. |
| `hotcues:hotpink` | Find tracks with at least one Hot Pink hot cue. |

### Color Palettes

Different palettes are used for Tracks and Cues to match Rekordbox's UI.

#### Track Colors

The following colors are available for track-level filtering:

- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#FF007F;margin-right:5px;"></span> `pink`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#FF0000;margin-right:5px;"></span> `red`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#FFA500;margin-right:5px;"></span> `orange`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#FFFF00;margin-right:5px;"></span> `yellow`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#00FF00;margin-right:5px;"></span> `green`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#25FDE9;margin-right:5px;"></span> `aqua`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#0000FF;margin-right:5px;"></span> `blue`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:#660099;margin-right:5px;"></span> `purple`
- <span style="display:inline-block;width:12px;height:12px;border-radius:50%;background:transparent;margin-right:5px;"></span> `none`

#### Hot Cue Colors

Use the following names to match the 16-color pad palette. Cues with no color set match `none`.

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

