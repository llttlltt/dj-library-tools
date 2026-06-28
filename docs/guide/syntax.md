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

### Collection Stats

Stats allow you to perform calculations on a collection.

| Stat | Description | Example |
| :--- | :--- | :--- |
| `-count` | Total items in the collection. | `hotcues-count:8` |
| `-density` | Items per minute of duration. | `beatgrids-density:>60` |
| `-drift` | Max value minus Min value. | `beatgrids/bpm-drift:<0.1` |

### Path Query Examples

| Query | Description |
| :--- | :--- |
| `hotcues.1/color:red` | Match if the first hotcue is red. |
| `beatgrids/bpm-drift:<0.1` | Match if the BPM varies by less than 0.1. |
| `memorycues/name:Drop` | Match if any memory cue is named "Drop". |

### Supported Collections

| Collection | Properties | Stats |
| :--- | :--- | :--- |
| `hotcues` | `color`, `name`, `position` | `-count` |
| `memorycues` | `color`, `name`, `position` | `-count` |
| `beatgrids` | `bpm`, `position` | `-count`, `-drift`, `-density` |

## Advanced Filters

Some providers support advanced filtering patterns (e.g., membership checks or cue properties). See the individual **[Providers](query/providers/index.md)** pages for details on which fields support these.
