package query

import (
	"testing"
)

func TestQueryValidation(t *testing.T) {
	allowedFields := []string{"title", "artist", "bpm"}

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{"Valid static field", "title:Oceans", false},
		{"Invalid static field", "nonsense:foo", true},
		{"Valid path field", "beatgrids-count:>0", false},
		{"Valid indexed path", "hotcues.1/color:red", false},
		{"Invalid collection in path", "invalid/prop:val", true},
		{"Complex valid query", "artist:Four && beatgrids/bpm-drift:<0.1", false},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := parser.Parse(tt.query)
			err := q.ValidateWithFields(allowedFields)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWithFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
