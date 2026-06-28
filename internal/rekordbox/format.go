package rekordbox

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

// Formatter defines the interface for high-fidelity Rekordbox XML emission.
type Formatter interface {
	Format(r io.Reader, w io.Writer) error
}

// Profile defines tag-specific formatting rules.
type Profile struct {
	// AttributeOrder defines the preferred order of attributes for specific tags.
	// If a tag is not present, attributes are emitted in their source order.
	AttributeOrder map[string][]string

	// TagSpecificOrder allows different attribute ordering for the same tag name
	// based on a predicate function that inspects the StartElement.
	TagSpecificOrder []TagOrderRule
}

// TagOrderRule defines a rule for attribute ordering based on a predicate.
type TagOrderRule struct {
	TagName   string
	Predicate func(xml.StartElement) bool
	Order     []string
}

// XMLFormat stores formatting preferences for Rekordbox XML.
type XMLFormat struct {
	Indent      string
	SelfClosing bool
	LineEnding  string
	LineLength  int
	Profile     *Profile
}

// TokenStreamFormatter implements the Formatter interface using a token-stream approach.
type TokenStreamFormatter struct {
	Config *XMLFormat
}

func NewTokenStreamFormatter(config *XMLFormat) *TokenStreamFormatter {
	if config == nil {
		config = DefaultFormat()
	}
	return &TokenStreamFormatter{Config: config}
}

func (f *TokenStreamFormatter) Format(r io.Reader, w io.Writer) error {
	decoder := xml.NewDecoder(r)
	var tokens []xml.Token

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		tokens = append(tokens, xml.CopyToken(token))
	}

	currentIndent := ""
	for i := 0; i < len(tokens); i++ {
		switch t := tokens[i].(type) {
		case xml.StartElement:
			skip, selfClosing := f.selfCloseAdvance(tokens, i, t.Name)
			i += skip
			if err := f.writeStartElement(w, t, selfClosing, currentIndent); err != nil {
				return err
			}
		case xml.EndElement:
			if _, err := fmt.Fprintf(w, "</%s>", t.Name.Local); err != nil {
				return err
			}
		case xml.CharData:
			if _, err := w.Write(t); err != nil {
				return err
			}
			if lastNL := bytes.LastIndex(t, []byte("\n")); lastNL != -1 {
				currentIndent = string(t[lastNL+1:])
			}
		case xml.ProcInst:
			if _, err := fmt.Fprintf(w, "<?%s %s?>", t.Target, string(t.Inst)); err != nil {
				return err
			}
		case xml.Directive:
			if _, err := fmt.Fprintf(w, "<!%s>", string(t)); err != nil {
				return err
			}
		case xml.Comment:
			if _, err := fmt.Fprintf(w, "<!--%s-->", string(t)); err != nil {
				return err
			}
		}
	}
	return nil
}

// selfCloseAdvance reports whether the element at tokens[i] can be emitted as
// self-closing and returns the number of additional tokens to consume (1 for an
// immediate EndElement, 2 for whitespace-only CharData followed by an EndElement).
func (f *TokenStreamFormatter) selfCloseAdvance(tokens []xml.Token, i int, name xml.Name) (skip int, selfClosing bool) {
	if !f.Config.SelfClosing || i+1 >= len(tokens) {
		return 0, false
	}
	if ee, ok := tokens[i+1].(xml.EndElement); ok && ee.Name == name {
		return 1, true
	}
	if cd, ok := tokens[i+1].(xml.CharData); ok && isWhitespace(cd) && i+2 < len(tokens) {
		if ee, ok := tokens[i+2].(xml.EndElement); ok && ee.Name == name {
			return 2, true
		}
	}
	return 0, false
}

func isWhitespace(data xml.CharData) bool {
	return len(bytes.TrimSpace(data)) == 0
}

func (f *TokenStreamFormatter) writeStartElement(w io.Writer, se xml.StartElement, selfClosing bool, currentIndent string) error {
	tagName := se.Name.Local
	if _, err := fmt.Fprintf(w, "<%s", tagName); err != nil {
		return err
	}

	attrs := f.resolveAttrs(tagName, se)
	startLen := len(currentIndent) + 1 + len(tagName)
	wrap := needsWrap(attrs, startLen, f.Config.LineLength)

	currentLineLen := startLen
	for i, attr := range attrs {
		encoded := fmt.Sprintf(" %s=\"%s\"", attr.Name.Local, escapeAttr(attr.Value))
		if wrap && i > 0 && currentLineLen+decodedAttrLen(attr) > f.Config.LineLength {
			// Align continuation with the start of the first attribute.
			indentSize := startLen + 1
			if _, err := fmt.Fprintf(w, "\n%s%s", strings.Repeat(" ", indentSize), encoded[1:]); err != nil {
				return err
			}
			currentLineLen = indentSize + len(encoded) - 1
			continue
		}
		if _, err := fmt.Fprint(w, encoded); err != nil {
			return err
		}
		currentLineLen += len(encoded)
	}

	closing := ">"
	if selfClosing {
		closing = "/>"
	}
	_, err := fmt.Fprint(w, closing)
	return err
}

// resolveAttrs returns the attributes of se ordered according to the profile
// rules for tagName, falling back to source order if no rule matches.
func (f *TokenStreamFormatter) resolveAttrs(tagName string, se xml.StartElement) []xml.Attr {
	if f.Config.Profile == nil {
		return se.Attr
	}
	for _, rule := range f.Config.Profile.TagSpecificOrder {
		if rule.TagName == tagName && rule.Predicate(se) {
			return sortAttributes(se.Attr, rule.Order)
		}
	}
	if order, ok := f.Config.Profile.AttributeOrder[tagName]; ok {
		return sortAttributes(se.Attr, order)
	}
	return se.Attr
}

// needsWrap reports whether the element's attributes need to be wrapped across
// multiple lines. Wrapping is only needed when the full decoded single-line
// length exceeds lineLength+8 — matching observed Rekordbox fixture behaviour
// where lines up to 88 decoded chars are never wrapped with the default
// lineLength=80.
func needsWrap(attrs []xml.Attr, startLen, lineLength int) bool {
	if lineLength <= 0 {
		return false
	}
	total := startLen
	for _, a := range attrs {
		total += decodedAttrLen(a)
	}
	return total > lineLength+8
}

// decodedAttrLen returns the decoded (un-escaped) rendered length of a single
// attribute including its leading space: ` name="value"`.
func decodedAttrLen(a xml.Attr) int {
	return 1 + len(a.Name.Local) + 2 + len(a.Value) + 1
}

func sortAttributes(attrs []xml.Attr, order []string) []xml.Attr {
	orderMap := make(map[string]int, len(order))
	for i, name := range order {
		orderMap[name] = i
	}

	sorted := make([]xml.Attr, len(attrs))
	copy(sorted, attrs)

	sort.SliceStable(sorted, func(i, j int) bool {
		posI, okI := orderMap[sorted[i].Name.Local]
		posJ, okJ := orderMap[sorted[j].Name.Local]
		switch {
		case okI && okJ:
			return posI < posJ
		case okI:
			return true
		case okJ:
			return false
		default:
			return false // let SliceStable preserve the original relative order
		}
	})

	return sorted
}

func escapeAttr(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&apos;")
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// DefaultProfile returns the standard Rekordbox attribute ordering.
func DefaultProfile() *Profile {
	return &Profile{
		AttributeOrder: map[string][]string{
			"TRACK": {
				"TrackID", "Name", "Artist", "Composer", "Album", "Grouping", "Genre",
				"Kind", "Size", "TotalTime", "DiscNumber", "TrackNumber", "Year",
				"AverageBpm", "DateAdded", "BitRate", "SampleRate", "Comments",
				"PlayCount", "Rating", "Location", "Remixer", "Tonality", "Label", "Mix",
			},
			"POSITION_MARK": {
				"Name", "Type", "Start", "End", "Num", "Red", "Green", "Blue",
			},
			"TEMPO": {
				"Inizio", "Bpm", "Metro", "Battito",
			},
		},
		TagSpecificOrder: []TagOrderRule{
			{
				TagName: "NODE",
				Predicate: func(se xml.StartElement) bool {
					for _, a := range se.Attr {
						if a.Name.Local == "Name" && a.Value == "ROOT" {
							return true
						}
					}
					return false
				},
				Order: []string{"Type", "Name", "Count"},
			},
			{
				TagName: "NODE",
				Predicate: func(se xml.StartElement) bool {
					for _, a := range se.Attr {
						if a.Name.Local == "Type" && a.Value == "0" {
							return true
						}
					}
					return false
				},
				Order: []string{"Name", "Type", "Count"},
			},
			{
				TagName: "NODE",
				Predicate: func(se xml.StartElement) bool {
					for _, a := range se.Attr {
						if a.Name.Local == "Type" && a.Value == "1" {
							return true
						}
					}
					return false
				},
				Order: []string{"Name", "Type", "KeyType", "Entries"},
			},
		},
	}
}

// DefaultFormat returns the default Rekordbox XML formatting.
func DefaultFormat() *XMLFormat {
	return &XMLFormat{
		Indent:      "  ",
		SelfClosing: true,
		LineEnding:  "\n",
		LineLength:  80,
		Profile:     DefaultProfile(),
	}
}

// indentRe matches leading spaces before a tag, used by DetectFormat.
var indentRe = regexp.MustCompile(`^([ ]+)<`)

// DetectFormat attempts to guess the formatting from the XML data.
func DetectFormat(data []byte) *XMLFormat {
	format := DefaultFormat()

	if bytes.Contains(data, []byte("\r\n")) {
		format.LineEnding = "\r\n"
	}

	for _, line := range bytes.Split(data, []byte(format.LineEnding)) {
		if len(line) == 0 || bytes.Contains(line, []byte("<?xml")) {
			continue
		}
		if m := indentRe.FindSubmatch(line); len(m) > 1 {
			if n := len(m[1]); n > 0 && (format.Indent == "  " || n < len(format.Indent)) {
				format.Indent = string(m[1])
			}
		}
	}

	return format
}
