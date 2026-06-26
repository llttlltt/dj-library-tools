# Selection Syntax

`djlt` uses a consistent selection syntax across all commands:

`[provider/resource] [query]`

## Providers & Resources

Resources are identified by their provider prefix. If no resource is specified, a default is used based on the provider.

- `rb/tracks` (Default for `rb`)
- `rb/playlists`
- `rb/folders`
- `plex/playlists` (Default for `plex`)
- `plex/tracks`

## Query Syntax

The query part supports a powerful set of operators and boolean logic.

### Operators

| Operator | Type | Example | Description |
| :--- | :--- | :--- | :--- |
| `:` | String | `artist:Daft` | Case-insensitive substring match. |
| `:` | Numeric | `bpm:124` | **Exact** numeric equality. |
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

## Advanced Filters

Some providers support advanced filtering patterns (e.g., membership checks or cue properties). See the individual **[Providers](query/providers/index.md)** pages for details on which fields support these.
