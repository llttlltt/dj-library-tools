# High-Fidelity Rekordbox XML Formatting

The `internal/rekordbox` package provides a `TokenStreamFormatter` designed to emit XML that is bit-for-bit compatible with the idiosyncratic formatting used by Rekordbox.

## Key Formatting Rules

### 1. Attribute Ordering
Rekordbox expects attributes in a specific order for different tag types. For example, a `TRACK` ResourceGroup must have `TrackID` before `Name`. These orders are defined in `DefaultProfile()`.

- **TRACK**: `TrackID`, `Name`, `Artist`, `Composer`, `Album`, ...
- **NODE (Playlist)**: `Name`, `Type`, `KeyType`, `Entries`
- **NODE (Folder)**: `Name`, `Type`, `Count`
- **NODE (Root)**: `Type`, `Name`, `Count`

### 2. Intelligent Attribute Wrapping
Unlike standard XML formatters that wrap at a fixed character limit, Rekordbox uses a specific "pre-wrap tolerance" rule:

1. **Threshold**: The full single-line string length is calculated using **decoded** characters.
2. **Decision**: Wrapping is only triggered if this total exceeds **88 characters** (assuming a base `LineLength` of 80).
3. **Greedy Wrap**: Once wrapping is triggered, attributes are placed one by one. If an attribute would push the current line past **80 characters** (decoded), it is wrapped to a new line.
4. **Alignment**: Wrapped attributes are indented to align with the start of the first attribute on the opening line (typically `len(indent) + len("<TAG ") + 1`).

### 3. Entity Encoding
Attributes are escaped using standard XML entities, but the formatter ensures consistency with Rekordbox's expectations for characters like apostrophes (`&apos;`) and ampersands (`&amp;`).

### 4. Self-Closing Tags
The formatter performs a look-ahead on the token stream. If a `StartElement` is immediately followed by its matching `EndElement` (or only by whitespace), it is emitted as a self-closing tag: `<NODE ... />`.

## Usage

```go
format := rekordbox.DefaultFormat()
formatter := rekordbox.NewTokenStreamFormatter(format)

var output bytes.Buffer
err := formatter.Format(xmlReader, &output)
```

## Testing
Formatting rules are pinned by comprehensive tests in `internal/rekordbox/format_test.go`. Any changes to the wrapping or ordering logic must pass the high-fidelity fixture tests:

- `20260524 Terracotta - Shortlist (8-space context)`
- `Drum & Bass (Correct Indent Wrap)`
- `Terracotta Lamma (Wrap at 80 chars absolute)`
