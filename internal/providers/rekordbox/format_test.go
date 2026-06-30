package rekordbox

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

// attrNames is a test helper that extracts attribute name order from a slice.
func attrNames(attrs []xml.Attr) []string {
	names := make([]string, len(attrs))
	for i, a := range attrs {
		names[i] = a.Name.Local
	}
	return names
}

// ---------------------------------------------------------------------------
// NewTokenStreamFormatter
// ---------------------------------------------------------------------------

func TestNewTokenStreamFormatter(t *testing.T) {
	t.Run("nil config falls back to DefaultFormat", func(t *testing.T) {
		f := NewTokenStreamFormatter(nil)
		if f.Config == nil {
			t.Fatal("Config should not be nil")
		}
		if f.Config.LineLength != 80 {
			t.Errorf("LineLength = %d, want 80", f.Config.LineLength)
		}
		if !f.Config.SelfClosing {
			t.Error("SelfClosing should be true")
		}
	})

	t.Run("explicit config is preserved", func(t *testing.T) {
		cfg := &XMLFormat{LineLength: 120, SelfClosing: false}
		f := NewTokenStreamFormatter(cfg)
		if f.Config.LineLength != 120 {
			t.Errorf("LineLength = %d, want 120", f.Config.LineLength)
		}
		if f.Config.SelfClosing {
			t.Error("SelfClosing should be false")
		}
	})
}

// ---------------------------------------------------------------------------
// DefaultFormat / DefaultProfile
// ---------------------------------------------------------------------------

func TestDefaultFormat(t *testing.T) {
	f := DefaultFormat()

	if f.Indent != "  " {
		t.Errorf("Indent = %q, want %q", f.Indent, "  ")
	}
	if !f.SelfClosing {
		t.Error("SelfClosing should be true")
	}
	if f.LineEnding != "\n" {
		t.Errorf("LineEnding = %q, want \"\\n\"", f.LineEnding)
	}
	if f.LineLength != 80 {
		t.Errorf("LineLength = %d, want 80", f.LineLength)
	}
	if f.Profile == nil {
		t.Error("Profile should not be nil")
	}
}

func TestDefaultProfile(t *testing.T) {
	p := DefaultProfile()

	t.Run("TRACK order starts with TrackID then Name", func(t *testing.T) {
		order, ok := p.AttributeOrder["TRACK"]
		if !ok {
			t.Fatal("TRACK not in AttributeOrder")
		}
		if order[0] != "TrackID" || order[1] != "Name" {
			t.Errorf("TRACK order[0:2] = %v, want [TrackID Name]", order[:2])
		}
	})

	t.Run("POSITION_MARK order starts with Name", func(t *testing.T) {
		order, ok := p.AttributeOrder["POSITION_MARK"]
		if !ok {
			t.Fatal("POSITION_MARK not in AttributeOrder")
		}
		if order[0] != "Name" {
			t.Errorf("POSITION_MARK order[0] = %q, want Name", order[0])
		}
	})

	t.Run("TEMPO order starts with Inizio", func(t *testing.T) {
		order, ok := p.AttributeOrder["TEMPO"]
		if !ok {
			t.Fatal("TEMPO not in AttributeOrder")
		}
		if order[0] != "Inizio" {
			t.Errorf("TEMPO order[0] = %q, want Inizio", order[0])
		}
	})

	t.Run("three NODE TagSpecificOrder rules", func(t *testing.T) {
		if len(p.TagSpecificOrder) != 3 {
			t.Fatalf("TagSpecificOrder len = %d, want 3", len(p.TagSpecificOrder))
		}
		for _, rule := range p.TagSpecificOrder {
			if rule.TagName != "NODE" {
				t.Errorf("rule.TagName = %q, want NODE", rule.TagName)
			}
		}
	})

	t.Run("ROOT node rule puts Type first", func(t *testing.T) {
		rule := p.TagSpecificOrder[0]
		rootSE := xml.StartElement{
			Name: xml.Name{Local: "NODE"},
			Attr: []xml.Attr{{Name: xml.Name{Local: "Name"}, Value: "ROOT"}},
		}
		if !rule.Predicate(rootSE) {
			t.Error("first rule should match a Name=ROOT node")
		}
		if rule.Order[0] != "Type" {
			t.Errorf("ROOT rule Order[0] = %q, want Type", rule.Order[0])
		}
	})

	t.Run("Type=0 folder rule puts Name first", func(t *testing.T) {
		rule := p.TagSpecificOrder[1]
		folderSE := xml.StartElement{
			Name: xml.Name{Local: "NODE"},
			Attr: []xml.Attr{{Name: xml.Name{Local: "Type"}, Value: "0"}},
		}
		if !rule.Predicate(folderSE) {
			t.Error("second rule should match a Type=0 node")
		}
		if rule.Order[0] != "Name" {
			t.Errorf("folder rule Order[0] = %q, want Name", rule.Order[0])
		}
	})

	t.Run("Type=1 playlist rule is Name/Type/KeyType/Entries", func(t *testing.T) {
		rule := p.TagSpecificOrder[2]
		playlistSE := xml.StartElement{
			Name: xml.Name{Local: "NODE"},
			Attr: []xml.Attr{{Name: xml.Name{Local: "Type"}, Value: "1"}},
		}
		if !rule.Predicate(playlistSE) {
			t.Error("third rule should match a Type=1 node")
		}
		want := []string{"Name", "Type", "KeyType", "Entries"}
		for i, w := range want {
			if rule.Order[i] != w {
				t.Errorf("playlist rule Order[%d] = %q, want %q", i, rule.Order[i], w)
			}
		}
	})
}

// ---------------------------------------------------------------------------
// isWhitespace
// ---------------------------------------------------------------------------

func TestIsWhitespace(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{" \t\n\r", true},
		{"\n        ", true},
		{"text", false},
		{"  x  ", false},
		{" \t text \n", false},
	}
	for _, c := range cases {
		got := isWhitespace(xml.CharData(c.input))
		if got != c.expected {
			t.Errorf("isWhitespace(%q) = %v, want %v", c.input, got, c.expected)
		}
	}
}

// ---------------------------------------------------------------------------
// escapeAttr
// ---------------------------------------------------------------------------

func TestEscapeAttr(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"no special chars", "no special chars"},
		{`say "hi"`, "say &quot;hi&quot;"},
		{"it's", "it&apos;s"},
		{"a & b", "a &amp; b"},
		{"a < b", "a &lt; b"},
		{"a > b", "a &gt; b"},
		{`"'&<>`, "&quot;&apos;&amp;&lt;&gt;"},
		{"Drum & Bass, Jungle", "Drum &amp; Bass, Jungle"},
		{"20260619 Mike's BBQ", "20260619 Mike&apos;s BBQ"},
	}
	for _, c := range cases {
		got := escapeAttr(c.input)
		if got != c.expected {
			t.Errorf("escapeAttr(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

// ---------------------------------------------------------------------------
// decodedAttrLen
// ---------------------------------------------------------------------------

func TestDecodedAttrLen(t *testing.T) {
	cases := []struct {
		name     string
		attrName string
		value    string
		expected int
	}{
		// formula: 1 (space) + len(name) + 2 (=") + len(value) + 1 (")
		{"Type=1", "Type", "1", 9},
		{"KeyType=0", "KeyType", "0", 12},
		{"Entries=19", "Entries", "19", 13},
		{"Entries=104", "Entries", "104", 14},
		{"Name=Shortlist", "Name", "20260524 Terracotta - Shortlist", 39},
		// Value is decoded: & counts as 1 char, not 5 (&amp;)
		{"Name=Drum&Bass decoded", "Name", "Drum & Bass, Jungle, Breakbeat, Bassline", 48},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := xml.Attr{Name: xml.Name{Local: c.attrName}, Value: c.value}
			got := decodedAttrLen(a)
			if got != c.expected {
				t.Errorf("decodedAttrLen(%q=%q) = %d, want %d", c.attrName, c.value, got, c.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// needsWrap
// ---------------------------------------------------------------------------

func TestNeedsWrap(t *testing.T) {
	makeAttrs := func(nameVal string, entriesVal string) []xml.Attr {
		return []xml.Attr{
			{Name: xml.Name{Local: "Name"}, Value: nameVal},
			{Name: xml.Name{Local: "Type"}, Value: "1"},
			{Name: xml.Name{Local: "KeyType"}, Value: "0"},
			{Name: xml.Name{Local: "Entries"}, Value: entriesVal},
		}
	}

	cases := []struct {
		name       string
		attrs      []xml.Attr
		startLen   int
		lineLength int
		expected   bool
	}{
		{
			name:       "lineLength=0 never wraps",
			attrs:      makeAttrs("anything", "1"),
			startLen:   5,
			lineLength: 0,
			expected:   false,
		},
		{
			// standalone: startLen=5, decoded total = 5+45+9+12+13 = 84 ≤ 88
			name:       "Terracotta Lamma standalone (total=84)",
			attrs:      makeAttrs("2026-05-16 Terracotta Lamma - Prep #1", "19"),
			startLen:   5,
			lineLength: 80,
			expected:   false,
		},
		{
			// 8-space context: startLen=13, decoded total = 13+39+9+12+14 = 87 ≤ 88
			name:       "Terracotta Shortlist 8-space (total=87)",
			attrs:      makeAttrs("20260524 Terracotta - Shortlist", "104"),
			startLen:   13,
			lineLength: 80,
			expected:   false,
		},
		{
			// 8-space context: startLen=13, decoded total = 13+45+9+12+13 = 92 > 88
			name:       "Terracotta Lamma Prep#1 8-space (total=92)",
			attrs:      makeAttrs("2026-05-16 Terracotta Lamma - Prep #1", "19"),
			startLen:   13,
			lineLength: 80,
			expected:   true,
		},
		{
			// Exactly at threshold: total = 88 → 88 > 88 is false
			name: "exactly at threshold (total=88)",
			attrs: []xml.Attr{
				{Name: xml.Name{Local: "a"}, Value: strings.Repeat("x", 80)},
			},
			startLen:   3, // 3 + (1+1+2+80+1) = 3+85 = 88
			lineLength: 80,
			expected:   false,
		},
		{
			// One over threshold: total = 89 → 89 > 88 is true
			name: "one over threshold (total=89)",
			attrs: []xml.Attr{
				{Name: xml.Name{Local: "a"}, Value: strings.Repeat("x", 81)},
			},
			startLen:   3, // 3 + (1+1+2+81+1) = 3+86 = 89
			lineLength: 80,
			expected:   true,
		},
		{
			name:       "empty attrs never wraps",
			attrs:      nil,
			startLen:   5,
			lineLength: 80,
			expected:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := needsWrap(c.attrs, c.startLen, c.lineLength)
			if got != c.expected {
				t.Errorf("needsWrap(startLen=%d, lineLength=%d) = %v, want %v",
					c.startLen, c.lineLength, got, c.expected)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// sortAttributes
// ---------------------------------------------------------------------------

func TestSortAttributes(t *testing.T) {
	t.Run("nil attrs returns empty", func(t *testing.T) {
		got := sortAttributes(nil, []string{"A"})
		if len(got) != 0 {
			t.Errorf("expected empty, got %v", got)
		}
	})

	t.Run("already in order", func(t *testing.T) {
		attrs := []xml.Attr{
			{Name: xml.Name{Local: "Name"}, Value: "x"},
			{Name: xml.Name{Local: "Type"}, Value: "1"},
		}
		got := sortAttributes(attrs, []string{"Name", "Type"})
		want := []string{"Name", "Type"}
		if names := attrNames(got); !equalSlices(names, want) {
			t.Errorf("got %v, want %v", names, want)
		}
	})

	t.Run("reversed order is sorted", func(t *testing.T) {
		attrs := []xml.Attr{
			{Name: xml.Name{Local: "Entries"}, Value: "10"},
			{Name: xml.Name{Local: "KeyType"}, Value: "0"},
			{Name: xml.Name{Local: "Type"}, Value: "1"},
			{Name: xml.Name{Local: "Name"}, Value: "x"},
		}
		want := []string{"Name", "Type", "KeyType", "Entries"}
		got := sortAttributes(attrs, want)
		if names := attrNames(got); !equalSlices(names, want) {
			t.Errorf("got %v, want %v", names, want)
		}
	})

	t.Run("attrs not in order map follow ordered attrs in original relative order", func(t *testing.T) {
		attrs := []xml.Attr{
			{Name: xml.Name{Local: "Z"}, Value: "z"},
			{Name: xml.Name{Local: "B"}, Value: "b"},
			{Name: xml.Name{Local: "A"}, Value: "a"},
		}
		got := sortAttributes(attrs, []string{"A", "B"})
		want := []string{"A", "B", "Z"}
		if names := attrNames(got); !equalSlices(names, want) {
			t.Errorf("got %v, want %v", names, want)
		}
	})

	t.Run("empty order preserves original order", func(t *testing.T) {
		attrs := []xml.Attr{
			{Name: xml.Name{Local: "B"}, Value: "b"},
			{Name: xml.Name{Local: "A"}, Value: "a"},
		}
		got := sortAttributes(attrs, nil)
		want := []string{"B", "A"}
		if names := attrNames(got); !equalSlices(names, want) {
			t.Errorf("got %v, want %v", names, want)
		}
	})

	t.Run("source slice is not modified", func(t *testing.T) {
		attrs := []xml.Attr{
			{Name: xml.Name{Local: "B"}, Value: "b"},
			{Name: xml.Name{Local: "A"}, Value: "a"},
		}
		sortAttributes(attrs, []string{"A", "B"})
		if attrs[0].Name.Local != "B" {
			t.Error("sortAttributes must not modify the source slice")
		}
	})
}

// ---------------------------------------------------------------------------
// selfCloseAdvance
// ---------------------------------------------------------------------------

func TestSelfCloseAdvance(t *testing.T) {
	name := xml.Name{Local: "NODE"}
	other := xml.Name{Local: "OTHER"}
	f := NewTokenStreamFormatter(nil)

	t.Run("SelfClosing=false returns (0, false)", func(t *testing.T) {
		f.Config.SelfClosing = false
		tokens := []xml.Token{xml.StartElement{Name: name}, xml.EndElement{Name: name}}
		skip, sc := f.selfCloseAdvance(tokens, 0, name)
		if skip != 0 || sc {
			t.Errorf("got (%d, %v), want (0, false)", skip, sc)
		}
		f.Config.SelfClosing = true
	})

	t.Run("no tokens after start returns (0, false)", func(t *testing.T) {
		tokens := []xml.Token{xml.StartElement{Name: name}}
		skip, sc := f.selfCloseAdvance(tokens, 0, name)
		if skip != 0 || sc {
			t.Errorf("got (%d, %v), want (0, false)", skip, sc)
		}
	})

	t.Run("immediate EndElement returns (1, true)", func(t *testing.T) {
		tokens := []xml.Token{xml.StartElement{Name: name}, xml.EndElement{Name: name}}
		skip, sc := f.selfCloseAdvance(tokens, 0, name)
		if skip != 1 || !sc {
			t.Errorf("got (%d, %v), want (1, true)", skip, sc)
		}
	})

	t.Run("whitespace CharData then EndElement returns (2, true)", func(t *testing.T) {
		tokens := []xml.Token{
			xml.StartElement{Name: name},
			xml.CharData("\n        "),
			xml.EndElement{Name: name},
		}
		skip, sc := f.selfCloseAdvance(tokens, 0, name)
		if skip != 2 || !sc {
			t.Errorf("got (%d, %v), want (2, true)", skip, sc)
		}
	})

	t.Run("non-whitespace CharData returns (0, false)", func(t *testing.T) {
		tokens := []xml.Token{
			xml.StartElement{Name: name},
			xml.CharData("content"),
			xml.EndElement{Name: name},
		}
		skip, sc := f.selfCloseAdvance(tokens, 0, name)
		if skip != 0 || sc {
			t.Errorf("got (%d, %v), want (0, false)", skip, sc)
		}
	})

	t.Run("mismatched EndElement returns (0, false)", func(t *testing.T) {
		tokens := []xml.Token{xml.StartElement{Name: name}, xml.EndElement{Name: other}}
		skip, sc := f.selfCloseAdvance(tokens, 0, name)
		if skip != 0 || sc {
			t.Errorf("got (%d, %v), want (0, false)", skip, sc)
		}
	})

	t.Run("whitespace then mismatched EndElement returns (0, false)", func(t *testing.T) {
		tokens := []xml.Token{
			xml.StartElement{Name: name},
			xml.CharData("\n  "),
			xml.EndElement{Name: other},
		}
		skip, sc := f.selfCloseAdvance(tokens, 0, name)
		if skip != 0 || sc {
			t.Errorf("got (%d, %v), want (0, false)", skip, sc)
		}
	})
}

// ---------------------------------------------------------------------------
// resolveAttrs
// ---------------------------------------------------------------------------

func TestResolveAttrs(t *testing.T) {
	trackSE := func(attrs ...xml.Attr) xml.StartElement {
		return xml.StartElement{Name: xml.Name{Local: "TRACK"}, Attr: attrs}
	}
	nodeSE := func(typ string, attrs ...xml.Attr) xml.StartElement {
		base := []xml.Attr{{Name: xml.Name{Local: "Type"}, Value: typ}}
		return xml.StartElement{Name: xml.Name{Local: "NODE"}, Attr: append(base, attrs...)}
	}

	t.Run("nil profile preserves source order", func(t *testing.T) {
		f := NewTokenStreamFormatter(&XMLFormat{SelfClosing: true, Profile: nil})
		se := trackSE(
			xml.Attr{Name: xml.Name{Local: "Name"}, Value: "Song"},
			xml.Attr{Name: xml.Name{Local: "TrackID"}, Value: "1"},
		)
		got := f.resolveAttrs("TRACK", se)
		if got[0].Name.Local != "Name" || got[1].Name.Local != "TrackID" {
			t.Errorf("expected source order, got %v", attrNames(got))
		}
	})

	t.Run("AttributeOrder applied for TRACK", func(t *testing.T) {
		f := NewTokenStreamFormatter(nil)
		se := trackSE(
			xml.Attr{Name: xml.Name{Local: "Name"}, Value: "Song"},
			xml.Attr{Name: xml.Name{Local: "TrackID"}, Value: "1"},
		)
		got := f.resolveAttrs("TRACK", se)
		if got[0].Name.Local != "TrackID" || got[1].Name.Local != "Name" {
			t.Errorf("expected (TrackID, Name), got %v", attrNames(got))
		}
	})

	t.Run("TagSpecificOrder used for Type=1 NODE", func(t *testing.T) {
		f := NewTokenStreamFormatter(nil)
		se := nodeSE("1",
			xml.Attr{Name: xml.Name{Local: "Entries"}, Value: "10"},
			xml.Attr{Name: xml.Name{Local: "KeyType"}, Value: "0"},
			xml.Attr{Name: xml.Name{Local: "Name"}, Value: "Test"},
		)
		got := f.resolveAttrs("NODE", se)
		want := []string{"Name", "Type", "KeyType", "Entries"}
		if names := attrNames(got); !equalSlices(names, want) {
			t.Errorf("got %v, want %v", names, want)
		}
	})

	t.Run("first matching TagSpecificOrder rule wins (ROOT before Type=0)", func(t *testing.T) {
		f := NewTokenStreamFormatter(nil)
		// ROOT node also satisfies Type=0, but the ROOT rule is listed first.
		se := xml.StartElement{
			Name: xml.Name{Local: "NODE"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "Name"}, Value: "ROOT"},
				{Name: xml.Name{Local: "Type"}, Value: "0"},
				{Name: xml.Name{Local: "Count"}, Value: "5"},
			},
		}
		got := f.resolveAttrs("NODE", se)
		// ROOT rule: Type, Name, Count
		if got[0].Name.Local != "Type" {
			t.Errorf("ROOT node: Order[0] = %q, want Type", got[0].Name.Local)
		}
	})

	t.Run("unknown tag preserves source order", func(t *testing.T) {
		f := NewTokenStreamFormatter(nil)
		se := xml.StartElement{
			Name: xml.Name{Local: "UNKNOWN"},
			Attr: []xml.Attr{
				{Name: xml.Name{Local: "B"}, Value: "b"},
				{Name: xml.Name{Local: "A"}, Value: "a"},
			},
		}
		got := f.resolveAttrs("UNKNOWN", se)
		if got[0].Name.Local != "B" || got[1].Name.Local != "A" {
			t.Errorf("expected source order, got %v", attrNames(got))
		}
	})
}

// ---------------------------------------------------------------------------
// DetectFormat
// ---------------------------------------------------------------------------

func TestDetectFormat(t *testing.T) {
	t.Run("LF line endings", func(t *testing.T) {
		f := DetectFormat([]byte("<ROOT>\n  <CHILD/>\n</ROOT>"))
		if f.LineEnding != "\n" {
			t.Errorf("LineEnding = %q, want \"\\n\"", f.LineEnding)
		}
	})

	t.Run("CRLF line endings", func(t *testing.T) {
		f := DetectFormat([]byte("<ROOT>\r\n  <CHILD/>\r\n</ROOT>"))
		if f.LineEnding != "\r\n" {
			t.Errorf("LineEnding = %q, want \"\\r\\n\"", f.LineEnding)
		}
	})

	t.Run("detects 4-space indent", func(t *testing.T) {
		f := DetectFormat([]byte("<ROOT>\n    <CHILD/>\n</ROOT>"))
		if f.Indent != "    " {
			t.Errorf("Indent = %q, want \"    \"", f.Indent)
		}
	})

	t.Run("detects 2-space indent", func(t *testing.T) {
		f := DetectFormat([]byte("<ROOT>\n  <CHILD/>\n</ROOT>"))
		if f.Indent != "  " {
			t.Errorf("Indent = %q, want \"  \"", f.Indent)
		}
	})

	t.Run("skips xml declaration line", func(t *testing.T) {
		// The <?xml ...?> line has no leading spaces; the indent is on the next line.
		f := DetectFormat([]byte("<?xml version=\"1.0\"?>\n  <ROOT/>\n"))
		if f.Indent != "  " {
			t.Errorf("Indent = %q, want \"  \"", f.Indent)
		}
	})

	t.Run("SelfClosing is always true", func(t *testing.T) {
		f := DetectFormat([]byte("<ROOT></ROOT>"))
		if !f.SelfClosing {
			t.Error("SelfClosing should be true")
		}
	})

	t.Run("Profile is set to DefaultProfile", func(t *testing.T) {
		f := DetectFormat([]byte("<ROOT/>"))
		if f.Profile == nil {
			t.Fatal("Profile should not be nil")
		}
		if _, ok := f.Profile.AttributeOrder["TRACK"]; !ok {
			t.Error("Profile should contain TRACK ordering")
		}
	})
}

// ---------------------------------------------------------------------------
// Format — token passthrough and error handling
// ---------------------------------------------------------------------------

func TestTokenStreamFormatter_Format(t *testing.T) {
	f := NewTokenStreamFormatter(nil)
	format := func(input string) (string, error) {
		var out bytes.Buffer
		err := f.Format(strings.NewReader(input), &out)
		return out.String(), err
	}

	t.Run("ProcInst passes through unchanged", func(t *testing.T) {
		input := `<?xml version="1.0" encoding="UTF-8"?>`
		got, err := format(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != input {
			t.Errorf("got %q, want %q", got, input)
		}
	})

	t.Run("Comment passes through unchanged", func(t *testing.T) {
		input := `<!-- a comment -->`
		got, err := format(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != input {
			t.Errorf("got %q, want %q", got, input)
		}
	})

	t.Run("Directive passes through unchanged", func(t *testing.T) {
		input := `<!DOCTYPE rekordbox>`
		got, err := format(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != input {
			t.Errorf("got %q, want %q", got, input)
		}
	})

	t.Run("malformed XML returns an error", func(t *testing.T) {
		_, err := format("<unclosed")
		if err == nil {
			t.Error("expected an error for malformed XML")
		}
	})

	t.Run("element with children is not self-closed", func(t *testing.T) {
		input := `<NODE Name="Shows" Type="0" Count="1"><NODE Name="Test" Type="1" KeyType="0" Entries="0"/></NODE>`
		want := `<NODE Name="Shows" Type="0" Count="1"><NODE Name="Test" Type="1" KeyType="0" Entries="0"/></NODE>`
		got, err := format(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

// ---------------------------------------------------------------------------
// TestTokenStreamFormatter_AttributeOrdering (existing)
// ---------------------------------------------------------------------------

func TestTokenStreamFormatter_AttributeOrdering(t *testing.T) {
	input := `<TRACK Name="Song" TrackID="123" Artist="Me"></TRACK>`
	expected := `<TRACK TrackID="123" Name="Song" Artist="Me"/>`

	format := DefaultFormat()
	formatter := NewTokenStreamFormatter(format)

	var out bytes.Buffer
	err := formatter.Format(strings.NewReader(input), &out)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	if out.String() != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, out.String())
	}
}

// ---------------------------------------------------------------------------
// TestTokenStreamFormatter_NodeOrdering (existing)
// ---------------------------------------------------------------------------

func TestTokenStreamFormatter_NodeOrdering(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ROOT Node",
			input:    `<NODE Name="ROOT" Count="16" Type="0"></NODE>`,
			expected: `<NODE Type="0" Name="ROOT" Count="16"/>`,
		},
		{
			name:     "Folder Node",
			input:    `<NODE Name="Shows" Type="0" Count="3"></NODE>`,
			expected: `<NODE Name="Shows" Type="0" Count="3"/>`,
		},
		{
			name:     "Playlist Node",
			input:    `<NODE Entries="114" Type="1" Name="20260619 Mike's BBQ" KeyType="0"></NODE>`,
			expected: `<NODE Name="20260619 Mike&apos;s BBQ" Type="1" KeyType="0" Entries="114"/>`,
		},
		{
			name:     "Terracotta Lamma - Prep #1 (no wrap, 84 decoded chars)",
			input:    `<NODE Name="2026-05-16 Terracotta Lamma - Prep #1" Type="1" KeyType="0" Entries="19"></NODE>`,
			expected: `<NODE Name="2026-05-16 Terracotta Lamma - Prep #1" Type="1" KeyType="0" Entries="19"/>`,
		},
		{
			name:     "Drum & Bass (no wrap standalone, 87 decoded chars)",
			input:    `<NODE Name="Drum &amp; Bass, Jungle, Breakbeat, Bassline" Type="1" KeyType="0" Entries="81"></NODE>`,
			expected: `<NODE Name="Drum &amp; Bass, Jungle, Breakbeat, Bassline" Type="1" KeyType="0" Entries="81"/>`,
		},
		{
			name: "20260524 Terracotta - Shortlist (8-space context)",
			input: `
        <NODE Name="20260524 Terracotta - Shortlist" Type="1" KeyType="0" Entries="104">
          <TRACK Key="1"/>
        </NODE>`,
			expected: `
        <NODE Name="20260524 Terracotta - Shortlist" Type="1" KeyType="0" Entries="104">
          <TRACK Key="1"/>
        </NODE>`,
		},
		{
			name: "2026-01-10 Disco to House #1 (nested in Sets folder)",
			input: `
      <NODE Name="Sets" Type="0" Count="5">
        <NODE Name="2026-01-10 Disco to House #1" Type="1" KeyType="0" Entries="20">
          <TRACK Key="1"/>
        </NODE>
      </NODE>`,
			expected: `
      <NODE Name="Sets" Type="0" Count="5">
        <NODE Name="2026-01-10 Disco to House #1" Type="1" KeyType="0" Entries="20">
          <TRACK Key="1"/>
        </NODE>
      </NODE>`,
		},
		{
			name: "2026-05-09 Jitterbug Session #1 (after sibling close)",
			input: `
      <NODE Name="Sets" Type="0" Count="2">
        <NODE Name="Prev" Type="1" KeyType="0" Entries="1">
          <TRACK Key="1"/>
        </NODE>
        <NODE Name="2026-05-09 Jitterbug Session #1" Type="1" KeyType="0" Entries="29">
          <TRACK Key="1"/>
        </NODE>
      </NODE>`,
			expected: `
      <NODE Name="Sets" Type="0" Count="2">
        <NODE Name="Prev" Type="1" KeyType="0" Entries="1">
          <TRACK Key="1"/>
        </NODE>
        <NODE Name="2026-05-09 Jitterbug Session #1" Type="1" KeyType="0" Entries="29">
          <TRACK Key="1"/>
        </NODE>
      </NODE>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format := DefaultFormat()
			formatter := NewTokenStreamFormatter(format)

			var out bytes.Buffer
			err := formatter.Format(strings.NewReader(tt.input), &out)
			if err != nil {
				t.Fatalf("Format failed: %v", err)
			}

			if out.String() != tt.expected {
				t.Errorf("expected:\n%s\ngot:\n%s", tt.expected, out.String())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
