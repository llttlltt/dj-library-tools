# Selection Syntax

`djlt` uses a consistent selection syntax across all commands:

`[provider/resource] [query]`

## Providers & Resources

Resources are identified by their provider and type. Both parts must be specified.

- `rb/tracks`
- `rb/playlists`
- `rb/folders`
- `plex/playlists`
- `plex/tracks`

## Query Syntax

The query part supports a powerful set of operators and boolean logic.

### Operators

| Operator | Type | Example | Description |
| :--- | :--- | :--- | :--- |
| `:` | String | `artist:Daft` | Case-insensitive substring match. |
| `:` | Numeric | `bpm:124` | **Exact** numeric equality. |
| `:=` | Exact | `title:='Music'` | Case-sensitive exact match. |
| `::` | Regex | `name::'^Ye'` | Regular expression match. |
| `..` | Range | `bpm:120..130` | Inclusive range match. |
| `>`, `<` | Comparison | `rating:>3` | Greater than / Less than. |
| `>=`, `<=` | Comparison | `rating:>=4` | Greater than or equal / Less than or equal. |

## Boolean Logic

| Logic | Syntax | Example |
| :--- | :--- | :--- |
| **AND** | `&&` | `genre:House && bpm:124` |
| **OR** | `\|\|` | `genre:House \|\| genre:Techno` |
| **NOT** | `!` | `!genre:Pop` |
| **Group** | `(...)` | `(genre:House \|\| genre:Techno) && rating:>3` |

## Deep Metadata (Path Syntax)

For collections like hotcues and beatgrids, you can traverse properties using a hierarchical path:

`[Collection] . [Index] / [Property] - [Stat]`

- **`hotcues.1/color:red`**: Find tracks where the first hotcue is red.
- **`beatgrids/bpm-drift:<0.1`**: Find tracks with very consistent "dynamic" grids.
- **`beatgrids-density:>60`**: Find "busy" grids with many markers.
- **`hotcues-count:8`**: Find tracks with all 8 hotcues used.

### Supported Collections

- `hotcues`: Properties: `color`, `name`, `position`.
- `memorycues`: Properties: `color`, `name`, `position`.
- `beatgrids`: Properties: `bpm`, `position`. Stats: `-drift`, `-density`, `-count`.

## Advanced Filters

Some providers support advanced filtering patterns (e.g., membership checks or cue properties). See the individual **[Providers](query/providers/index.md)** pages for details on which fields support these.
