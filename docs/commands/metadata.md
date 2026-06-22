# metadata

Manage track metadata between Rekordbox XML libraries

```
djlt metadata [flags]
```
### Options

```
  -d, --destination string   Destination Rekordbox XML (tracks receive the Tempo markers)
  -f, --force                Overwrite output file if it already exists
  -h, --help                 help for metadata
  -o, --output string        Output path for the merged Rekordbox XML
  -s, --source string        Source Rekordbox XML (Tempo markers are read from here)
```

### Inherited Options

```
  -x, --xml string   Path to the Rekordbox XML library
```

Reads two Rekordbox XML libraries: a source library from which Tempo
markers are copied, and a destination library whose tracks receive them.
Tracks are matched by strict metadata equality (Name, Artist, Album, etc.).
A merged library is written to the output path.

## See also

* [djlt](index.md)	 - DJ Library Tools CLI