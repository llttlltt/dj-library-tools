# rb-cli

A Go-based command-line utility for managing Rekordbox XML library metadata.

`rb-cli` provides tools to inspect and manipulate Rekordbox XML files, enabling precise control over track metadata across different library exports.

## 🛠 Installation

Currently, the recommended way to install `rb-cli` is via `go install`. This project uses `mise` for environment management; ensure you are using the specified Go version.

```bash
# Ensure Go 1.25 is available via mise
mise install

# Install the tool
go install github.com/llttlltt/dj-library-tools@latest
```

_Note: Future support for GitHub Releases and Homebrew is planned as the feature set matures._

## 🚀 Usage

`rb-cli` uses a nested command structure. The primary functionality currently implemented is within the `metadata` command group.

### Metadata Operations

#### Move Metadata

The `move` command matches tracks from a source XML to a destination XML based on strict metadata criteria and synchronizes specific fields (currently `Tempo`).

**Syntax:**

```bash
rb-cli metadata move [flags] --source <path> --destination <path> --output <path>
```

**Matching Criteria:**
Tracks are matched using strict equality on the following fields:

- Name
- Artist
- Composer
- Album
- Comments
- Disc Number
- Track Number
- Year

**Behavior:**
When a match is found, the `Tempo` from the source track is copied to the corresponding track in the destination file.

## 🗺 Roadmap & Architecture

`rb-cli` is built on a **"Selection + Action"** philosophy. A powerful core Query Engine defines the _selection_ (the "who"), and a set of command primitives performs the _action_ (the "what").

### 🔍 The Query Engine

The core of the project is a sophisticated query language. Future commands will support syntax like:
`rb-cli ls artist:"Four Tet" album::"Sixteen Oceans|There Is Love in You|New Energy" bpm:120..140`

#### Query Operators & Syntax

- **`:` (Colon)**: Performs a case-insensitive substring match (the default).
- **`=` (Equal Sign)**: Performs a strict, exact string match.
- **`::` (Double Colon)**: Signals a Regular Expression (regex) match.
- **`..` (Double Dot)**: Defines a range for numeric values or dates (e.g., `120..140`).
- **`~` (Tilde)**: Performs a "fuzzy" match (requires the `fuzzy` plugin).
- **`^` (Caret)**: Used within a regex (`::`) to match the start of a field.
- **`$` (Dollar Sign)**: Used within a regex (`::`) to match the end of a field.
- **`|` (Pipe)**: Used within a regex (`::`) as an "OR" operator to match multiple terms.
- **`!` (Exclamation Point)**: Used at the start of a query to negate it (finds everything that **doesn't** match).

### Core Primitives (Planned)

The following commands will serve as the fundamental building blocks of the tool:

- **`ls [QUERY]`**: List tracks that match the given query.
- **`stat [QUERY]`**: Analyze track statistics (summary, distribution, or musical mapping).
- **`modify [QUERY] <changes>`**: The primary engine for bulk updating track metadata.
- **`remove [QUERY]`**: Remove tracks or metadata entries from the library.
- **`sync [QUERY]`**: Reconcile differences between two library files.

### Workflow Shortcuts (Planned)

Once the primitives are stable, specialized "shortcuts" will be implemented to optimize common DJ workflows:

- **`tag [QUERY]`**: A high-level wrapper for `modify`, optimized for rapid addition/removal of semantic tags.
- **`audit [QUERY]`**: A specialized tool to find common library hygiene issues (missing BPM, empty artist fields, etc.).
- **`dedupe [QUERY]`**: A specialized utility to identify and resolve duplicate tracks.

## 🛠 Development

This project uses [mise](https://mise.jdx.dev/) for toolchain management.

```bash
# Install dependencies
mise install

# Run tests
go test ./...
```

## ⚖️ License

[Specify your license here, e.g., MIT]
